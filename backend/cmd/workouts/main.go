package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"fitness-trainer/internal/app"
	genai_client "fitness-trainer/internal/clients/gemini"
	"fitness-trainer/internal/clients/ratelimiter"
	s3_client "fitness-trainer/internal/clients/s3"
	"fitness-trainer/internal/db"
	"fitness-trainer/internal/jwt"
	"fitness-trainer/internal/logger"
	"fitness-trainer/internal/repository"
	"fitness-trainer/internal/service"
	workout_generator_service "fitness-trainer/internal/service/workout_generator"
	"fitness-trainer/internal/tracer"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/generative-ai-go/genai"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
	apiOpts "google.golang.org/api/option"
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

	genaiClient, err := newGeminiClient(ctx)
	if err != nil {
		return err
	}

	clientWrapper := genai_client.New(genaiClient)

	// openaiClient := newOpenAIClient()

	// clientWrapper := openai_client.New(openaiClient, os.Getenv("OPENAI_ASS_ID"))

	WorkoutGenerator := workout_generator_service.New(clientWrapper)

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerDay(5),
		MaxBurst: 5,
	}

	inmemmoryStore, err := memstore.NewCtx(65536)
	if err != nil {
		return fmt.Errorf("failed to create in memory store: %w", err)
	}

	rateLimiter, err := throttled.NewGCRARateLimiterCtx(inmemmoryStore, quota)
	if err != nil {
		return fmt.Errorf("failed to create rate limiter: %w", err)
	}

	rateLimiterWrapper := ratelimiter.New(rateLimiter)

	Service := service.New(
		ContextManager,
		JWTProvider,
		s3ClientWrapper,
		WorkoutGenerator,
		rateLimiterWrapper,
		Repo, // Auth
		Repo, // User
		Repo, // Exercise
		Repo, // Workout
		Repo, // ExerciseInstance
		Repo, // MuscleGroup
		Repo, // Workout
		Repo, // ExerciseLog
		Repo, // SetLog
		Repo, // Set
		Repo, // ExpectedSet
		Repo, // Generation Settings
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

type ProxyRoundTripper struct {
	proxy  *url.URL
	apiKey string
}

func (t *ProxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	if t.proxy != nil {
		transport.Proxy = http.ProxyURL(t.proxy)
	}

	newReq := req.Clone(req.Context())
	q := newReq.URL.Query()
	q.Add("key", t.apiKey)
	newReq.URL.RawQuery = q.Encode()

	return transport.RoundTrip(newReq)
}

func loadProxyData() *url.URL {
	proxyURL := os.Getenv("PROXY_URL")
	proxyUser := os.Getenv("PROXY_USER")
	proxyPassword := os.Getenv("PROXY_PASSWORD")

	if proxyURL == "" {
		return nil
	}

	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if proxyUser != "" && proxyPassword != "" {
		parsedURL.User = url.UserPassword(proxyUser, proxyPassword)
	}

	return parsedURL
}

func newHTTPClient(proxyURL *url.URL, apiKey string) *http.Client {
	return &http.Client{
		Transport: &ProxyRoundTripper{
			apiKey: apiKey,
			proxy:  proxyURL,
		},
		Timeout: 30 * time.Second,
	}
}

func newGeminiClient(ctx context.Context) (*genai.Client, error) {
	proxy := loadProxyData()

	return genai.NewClient(
		ctx,
		apiOpts.WithHTTPClient(newHTTPClient(proxy, os.Getenv("GENAI_API_KEY"))),
		apiOpts.WithAPIKey(os.Getenv("GENAI_API_KEY")),
	)
}

func newOpenAIClient() *openai.Client {
	proxyURL := loadProxyData()

	return openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		option.WithHeader("OpenAI-Beta", "assistants=v2"),
		option.WithHTTPClient(&http.Client{
			Transport: &ProxyRoundTripper{
				proxy: proxyURL,
			},
			Timeout: 30 * time.Second,
		}),
	)
}
