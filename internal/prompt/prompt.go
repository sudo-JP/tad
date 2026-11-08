package prompt

import (
	_ "embed"
	"fmt"
)

//go:embed system.txt
var systemPrompt string

func System() (string, error) {
	if systemPrompt == "" {
		return "", fmt.Errorf("system prompt is empty")
	}
	return systemPrompt, nil
}
