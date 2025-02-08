package openai_client

import (
	"context"
	"fitness-trainer/internal/domain"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/opentracing/opentracing-go"
)

type Client struct {
	client      *openai.Client
	assistantID string
}

func New(client *openai.Client, assistantID string) *Client {
	return &Client{
		client:      client,
		assistantID: assistantID,
	}
}

func (c *Client) CreateCompletion(ctx context.Context, userID domain.ID, systemPrompt, prompt string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client.CreateCompletion")
	defer span.Finish()

	// Get assistant
	assistant, err := c.client.Beta.Assistants.Get(ctx, c.assistantID)
	if err != nil {
		return "", fmt.Errorf("failed to get assistant: %w", err)
	}

	// Create thread
	thread, err := c.client.Beta.Threads.New(ctx, openai.BetaThreadNewParams{
		Messages: openai.F([]openai.BetaThreadNewParamsMessage{
			{
				Role: openai.F(openai.BetaThreadNewParamsMessagesRoleAssistant),
				Content: openai.F([]openai.MessageContentPartParamUnion{
					openai.TextContentBlockParam{
						Type: openai.F(openai.TextContentBlockParamTypeText),
						Text: openai.String(systemPrompt),
					},
				}),
			},
			{
				Role: openai.F(openai.BetaThreadNewParamsMessagesRoleUser),
				Content: openai.F([]openai.MessageContentPartParamUnion{
					openai.TextContentBlockParam{
						Type: openai.F(openai.TextContentBlockParamTypeText),
						Text: openai.String(prompt),
					},
				}),
			},
		}),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}
	defer c.client.Beta.Threads.Delete(ctx, thread.ID)

	// Create completion
	run, err := c.client.Beta.Threads.Runs.NewAndPoll(ctx, thread.ID, openai.BetaThreadRunNewParams{
		AssistantID: openai.F(assistant.ID),
	}, 0)
	if err != nil {
		return "", fmt.Errorf("failed to create run: %w", err)
	}

	messages, err := c.client.Beta.Threads.Messages.List(ctx, thread.ID, openai.BetaThreadMessageListParams{
		Limit: openai.Int(1),
		Order: openai.F(openai.BetaThreadMessageListParamsOrderDesc),
		RunID: openai.String(run.ID),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list messages: %w", err)
	}

	if len(messages.Data) == 0 {
		return "", fmt.Errorf("no messages returned")
	}

	return messages.Data[0].Content[0].Text.Value, nil
}
