package email

import (
	"context"
	"fitness-trainer/internal/clients/kafka"
	"fitness-trainer/internal/domain"

	"github.com/opentracing/opentracing-go"
)

type Service struct {
	producer kafka.Producer
	topic    string
}

func NewService(producer kafka.Producer, topic string) *Service {
	return &Service{
		producer: producer,
		topic:    topic,
	}
}

func (s *Service) SendWelcomeEmail(ctx context.Context, email, name string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.email.SendWelcomeEmail")
	defer span.Finish()

	payload := domain.WelcomePayload{
		Email: email,
		Name:  name,
	}
	message := domain.EmailMessage{
		Type:    domain.WelcomeEmail,
		Payload: payload,
	}
	return s.producer.Publish(ctx, s.topic, email, message)
}
