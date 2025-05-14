package domain

type EmailEventType string

const (
	WelcomeEmail EmailEventType = "welcome"
)

type WelcomePayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type EmailMessage struct {
	Type    EmailEventType `json:"type"`
	Payload any            `json:"payload"`
}
