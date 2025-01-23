package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"time"

	"fitness-trainer/internal/app"
	"fitness-trainer/internal/jwt"
	"fitness-trainer/internal/logger"
	"fitness-trainer/internal/repository"
	"fitness-trainer/internal/service"
	"fitness-trainer/internal/tracer"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	logger.Init()
	godotenv.Load()
}

func loadPostgresURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)
}

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tracer.MustSetup(
		ctx, 
		tracer.WithServiceName("fitness-trainer"),
		tracer.WithCollectorEndpoint(os.Getenv("JAEGER_COLLECTOR_ENDPOINT")),
	)

	postgresURL := loadPostgresURL()

	pool, err := pgxpool.New(ctx, postgresURL)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal(err.Error())
	}

	Repo := repository.NewPGXRepository(pool)

	JWTProvider := jwt.NewProvider(
		jwt.WithCredentials(
			jwt.NewSecretCredentials(os.Getenv("JWT_SECRET")),
		),
		jwt.WithAccessTTL(
			30*time.Minute,
		),
	)

	Service := service.New(
		JWTProvider, // JWT Provider
		Repo,        // Auth
		Repo,        // User
		Repo,        // Exercise
		Repo,        // Workout
		Repo,        // ExerciseInstance
		Repo,        // MuscleGroup
		Repo,        // Workout
		Repo,        // ExerciseLog
		Repo,        // SetLog
	)

	App := app.New(
		Service,
		Service,
		Service,
		Service,
		Service,
		app.WithHTTPPathPrefix("/api"),
	)

	if err := App.Run(ctx); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}
