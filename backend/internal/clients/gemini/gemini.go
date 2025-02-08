package genai_client

import (
	"context"
	"fitness-trainer/internal/domain"
	"fitness-trainer/internal/logger"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/opentracing/opentracing-go"
)

var responseSchema = &genai.Schema{
	Type:     genai.TypeObject,
	Enum:     []string{},
	Required: []string{"reasoning"},
	Properties: map[string]*genai.Schema{
		"exercises": &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type:     genai.TypeObject,
				Enum:     []string{},
				Required: []string{"id", "name"},
				Properties: map[string]*genai.Schema{
					"id": &genai.Schema{
						Type: genai.TypeString,
					},
					"name": &genai.Schema{
						Type: genai.TypeString,
					},
				},
			},
		},
		"reasoning": &genai.Schema{
			Type: genai.TypeString,
		},
	},
}

type Client struct {
	client *genai.Client
}

func New(client *genai.Client) *Client {
	return &Client{client: client}
}

func (c *Client) CreateCompletion(ctx context.Context, userID domain.ID, systemPrompt, prompt string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "genai.CreateCompletion")
	defer span.Finish()

	logger.Debugf("creating completion for user %s", userID)
	logger.Debugf("system prompt: %s", systemPrompt)
	logger.Debugf("user prompt: %s", prompt)

	model := c.client.GenerativeModel("gemini-2.0-flash")

	model.SetTemperature(1.8)
	model.SetTopK(40)
	model.SetTopP(0.9)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "application/json"
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}
	model.ResponseSchema = responseSchema

	session := model.StartChat()

	resp, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	var response string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			response += string(text)
		}
	}

	logger.Debugf("response: %s", response)

	return response, nil
}
