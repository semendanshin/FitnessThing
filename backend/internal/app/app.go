package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"fitness-trainer/internal/app/fitness-trainer/api/auth"
	"fitness-trainer/internal/app/fitness-trainer/api/exercise"
	"fitness-trainer/internal/app/fitness-trainer/api/routine"
	"fitness-trainer/internal/app/fitness-trainer/api/user"
	"fitness-trainer/internal/app/fitness-trainer/api/workout"
	"fitness-trainer/internal/app/interceptors"
	"fitness-trainer/internal/logger"
	desc "fitness-trainer/pkg/workouts"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/swaggest/swgui/v5emb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	host              string = "0.0.0.0"
	accessTokenHeader string = "X-Access-Token"
)

func CustomMatcher(key string) (string, bool) {
	switch key {
	case accessTokenHeader:
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

type Options struct {
	grpcPort         int
	gatewayPort      int
	enableGateway    bool
	enableReflection bool
}

var defaultOptions = &Options{
	grpcPort:         50051,
	gatewayPort:      8080,
	enableGateway:    true,
	enableReflection: true,
}

type OptionsFunc func(*Options)

func WithGrpcPort(port int) OptionsFunc {
	return func(o *Options) {
		o.grpcPort = port
	}
}

func WithGatewayPort(port int) OptionsFunc {
	return func(o *Options) {
		o.gatewayPort = port
	}
}

func WithEnableReflection(enableReflection bool) OptionsFunc {
	return func(o *Options) {
		o.enableReflection = enableReflection
	}
}

func WithEnableGateway(enableGateway bool) OptionsFunc {
	return func(o *Options) {
		o.enableGateway = enableGateway
	}
}

type App struct {
	authService     auth.Service
	userService     user.Service
	workoutService  workout.Service
	exerciseService exercise.Service
	routineService  routine.Service

	options *Options
}

func New(
	authService auth.Service,
	userService user.Service,
	workoutService workout.Service,
	exerciseService exercise.Service,
	routineService routine.Service,
	options ...OptionsFunc,
) *App {
	opts := defaultOptions
	for _, o := range options {
		o(opts)
	}
	return &App{
		authService:     authService,
		userService:     userService,
		workoutService:  workoutService,
		exerciseService: exerciseService,
		routineService:  routineService,
		options:         opts,
	}
}

func (a *App) Run(ctx context.Context) error {
	grpcEndpoint := fmt.Sprintf(":%d", a.options.grpcPort)
	httpEndpoint := fmt.Sprintf(":%d", a.options.gatewayPort)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.TracingInterceptor,
			interceptors.RecovertInterceptor,
			interceptors.NewAuth(
				a.authService,
				map[string]struct{}{
					"/fitness_trainer.api.workout.AuthService/Login":      {},
					"/fitness_trainer.api.workout.AuthService/Refresh":    {},
					"/fitness_trainer.api.workout.UserService/CreateUser": {},
				},
			),
			interceptors.ErrCodesInterceptor,
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(),
		),
	)

	workoutService := workout.New(a.workoutService)
	exerciseService := exercise.New(a.exerciseService)
	routineService := routine.New(a.routineService)
	authServiceServer := auth.New(a.authService)
	userServiceServer := user.New(a.userService)

	// Register the service
	desc.RegisterWorkoutServiceServer(srv, workoutService)
	desc.RegisterExerciseServiceServer(srv, exerciseService)
	desc.RegisterRoutineServiceServer(srv, routineService)
	desc.RegisterUserServiceServer(srv, userServiceServer)
	desc.RegisterAuthServiceServer(srv, authServiceServer)

	// Reflect the service
	if a.options.enableReflection {
		reflection.Register(srv)
	}

	// Create gatewayx
	gatewayMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(CustomMatcher),
	)

	if err := registerGateway(ctx, gatewayMux, grpcEndpoint); err != nil {
		return err
	}

	gatewayMuxWithCORS := cors.New(
		cors.Options{
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			AllowedHeaders: []string{"*"},
		},
	).Handler(gatewayMux)

	// Create swagger ui
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/swagger", func(w http.ResponseWriter, request *http.Request) {
		http.ServeFile(w, request, "pkg/workouts/workouts.swagger.json")
	})
	httpMux.Handle("/docs/", v5emb.NewHandler(
		"PVZ Service",
		"/swagger",
		"/docs/",
	))

	//httpMux.Handle("/metrics", promhttp.Handler())

	httpMux.Handle("/", gatewayMuxWithCORS)

	httpSrv := &http.Server{
		Addr:    httpEndpoint,
		Handler: httpMux,
	}

	// Start the gateway and swagger ui
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Errorf("error starting http server: %v", err)
			}
		}
	}()

	// Handle shutdown signals
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

		<-stop

		logger.Info("shutting down server...")

		err := httpSrv.Shutdown(ctx)
		if err != nil {
			logger.Errorf("error shutting down http server: %v", err)
		}

		srv.Stop()
	}()

	// Create listener
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		return err
	}

	logger.Infof("grpc server listening on port %d", a.options.grpcPort)

	// Start the server
	if err := srv.Serve(lis); err != nil {
		return err
	}

	logger.Infof("grpc server stopped")

	return nil
}

func registerGateway(ctx context.Context, mux *runtime.ServeMux, grpcEndpoint string) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterWorkoutServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	err = desc.RegisterExerciseServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	err = desc.RegisterRoutineServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	err = desc.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	err = desc.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return err
	}

	return nil
}
