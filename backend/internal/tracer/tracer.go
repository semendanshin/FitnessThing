package tracer

import (
	"context"
	"log"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	traceconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

func MustSetup(ctx context.Context, serviceName string) {
	log.Printf("Initializing tracer for service: %s", serviceName)
	cfg := traceconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &traceconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &traceconfig.ReporterConfig{
			// LogSpans: true,
			CollectorEndpoint: "http://localhost:14268/api/traces",
		},
	}

	tracer, closer, err := cfg.NewTracer(traceconfig.Logger(jaeger.StdLogger), traceconfig.Metrics(prometheus.New()))
	if err != nil {
		log.Fatalf("ERROR: cannot init Jaeger %s", err)
	}
	log.Printf("Successfully initialized Jaeger tracer")

	go func() {
		onceCloser := sync.OnceFunc(func() {
			log.Println("closing tracer")
			if err = closer.Close(); err != nil {
				log.Fatalf("ERROR: cannot close Jaeger %s", err)
			}
		})

		for {
			<-ctx.Done()
			onceCloser()
		}
	}()

	opentracing.SetGlobalTracer(tracer)
	log.Printf("Set global tracer successfully")
}
