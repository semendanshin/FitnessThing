package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, topic string, key string, message interface{}) error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Async:        false,
		},
	}
}

func headersFromContext(ctx context.Context) ([]kafka.Header, error) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil, nil
	}

	carrier := opentracing.TextMapCarrier{}
	err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		return nil, err
	}

	headers := make([]kafka.Header, 0, len(carrier))
	for k, v := range carrier {
		headers = append(headers, kafka.Header{Key: k, Value: []byte(v)})
	}
	return headers, nil
}

func (p *producer) Publish(ctx context.Context, topic string, key string, message any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "kafka.Producer.Publish")
	defer span.Finish()

	span.SetTag("topic", topic)
	span.SetTag("key", key)
	span.SetTag("message", message)

	value, err := json.Marshal(message)
	if err != nil {
		return err
	}

	headers, err := headersFromContext(ctx)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic:   topic,
		Key:     []byte(key),
		Value:   value,
		Time:    time.Now(),
		Headers: headers,
	}

	return p.writer.WriteMessages(ctx, msg)
}
