package handlers

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ollama/ollama/api"
	"github.com/sudo-JP/tad/internal/prompt"
)

func handleOllamaMessage(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	userPrompt := strings.Join(args, " ")
	systemPrompt, err := prompt.System()
	if err != nil {
		return
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return
	}

	messages := []api.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	req := &api.ChatRequest{
		Model:    "deepseek-r1:1.5b",
		Messages: messages,
	}

	var responseBuilder strings.Builder
	err = client.Chat(ctx, req, func(resp api.ChatResponse) error {
		responseBuilder.WriteString(resp.Message.Content)
		return nil
	})
	if err != nil {
		return
	}

	finalOutput := responseBuilder.String()
	if finalOutput != "" {
		reply(s, m, finalOutput)
	}
}
