package service

import (
	"context"

	"github.com/itsLeonB/billsplittr/internal/config"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/rotisserie/eris"
)

type openAILLMService struct {
	client openai.Client
	model  string
}

func NewLLMService(cfg config.LLM) LLMService {
	client := openai.NewClient(option.WithAPIKey(cfg.ApiKey), option.WithBaseURL(cfg.BaseUrl))
	return &openAILLMService{client, cfg.Model}
}

func (llm *openAILLMService) Prompt(ctx context.Context, systemMsg, userMsg string) (string, error) {
	if userMsg == "" {
		return "", eris.New("empty user message")
	}

	msgs := make([]openai.ChatCompletionMessageParamUnion, 0, 2)

	if systemMsg != "" {
		msgs = append(msgs, openai.SystemMessage(systemMsg))
	}

	msgs = append(msgs, openai.UserMessage(userMsg))

	params := openai.ChatCompletionNewParams{
		Model:    llm.model,
		Messages: msgs,
	}

	response, err := llm.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", eris.Wrap(err, "error getting LLM response")
	}

	if len(response.Choices) < 1 {
		return "", eris.New("no response from LLM")
	}

	return response.Choices[0].Message.Content, nil
}
