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

type options struct {
	ServiceName       string
	CollectorEndpoint string
}

type OptFunc func(*options)

func WithServiceName(serviceName string) OptFunc {
	return func(o *options) {
		o.ServiceName = serviceName
	}
}

func WithCollectorEndpoint(collectorEndpoint string) OptFunc {
	return func(o *options) {
		o.CollectorEndpoint = collectorEndpoint
	}
}

var defaultOptions = &options{
	ServiceName:       "service",
	CollectorEndpoint: "http://localhost:14268/api/traces",
}

func MustSetup(ctx context.Context, opts ...OptFunc) {
	o := defaultOptions
	for _, opt := range opts {
		opt(o)
	}

	log.Printf("Initializing tracer for service: %s", o.ServiceName)
	cfg := traceconfig.Configuration{
		ServiceName: o.ServiceName,
		Sampler: &traceconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &traceconfig.ReporterConfig{
			// LogSpans: true,
			CollectorEndpoint: o.CollectorEndpoint,
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
