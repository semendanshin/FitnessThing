package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"fitness-trainer/internal/app"
	s3_client "fitness-trainer/internal/clients/s3"
	"fitness-trainer/internal/db"
	"fitness-trainer/internal/jwt"
	"fitness-trainer/internal/logger"
	"fitness-trainer/internal/repository"
	"fitness-trainer/internal/service"
	"fitness-trainer/internal/tracer"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	logger.Init()
	godotenv.Load()
	log.SetOutput(io.Discard)
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

	endpoint := os.Getenv("AWS_ENDPOINT")
	bucket := os.Getenv("AWS_S3_BUCKET")

	awsConfig := getAWSConfig(ctx)
	s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	s3ClientWrapper := s3_client.New(s3Client, bucket)

	ContextManager := db.NewContextManager(pool)

	Repo := repository.NewPGXRepository(ContextManager)

	jwtSecret := os.Getenv("JWT_SECRET")

	JWTProvider := jwt.NewProvider(
		jwt.WithCredentials(
			jwt.NewSecretCredentials(jwtSecret),
		),
		jwt.WithAccessTTL(
			30*time.Minute,
		),
	)

	Service := service.New(
		ContextManager,
		JWTProvider,     // JWT Provider
		s3ClientWrapper, // S3 Client
		Repo,            // Auth
		Repo,            // User
		Repo,            // Exercise
		Repo,            // Workout
		Repo,            // ExerciseInstance
		Repo,            // MuscleGroup
		Repo,            // Workout
		Repo,            // ExerciseLog
		Repo,            // SetLog
		Repo,            // Set
		Repo,            // ExpectedSet
	)

	App := app.New(
		Service,
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

func getAWSConfig(ctx context.Context) aws.Config {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	customRegion := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(customRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)

	if err != nil {
		log.Fatal("Unable to load AWS config:", err)
	}

	return cfg
}
