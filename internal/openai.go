package internal

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type GPTClient struct {
    client *openai.Client
}

func NewGPTClient(apiKey string) *GPTClient {
    c := openai.NewClient(apiKey)
    return &GPTClient{
        client: c,
    }
}

// GenerateTasks uses the user story & context to produce a list of tasks
func (g *GPTClient) GenerateTasks(ctx context.Context, projectContext ProjectContext, userStory string) ([]string, error) {
    prompt := fmt.Sprintf(`
Given the following project context and user story, please produce a detailed list of tasks:

Project Context:
- Language: %s
- Framework: %s
- Role: %s
- Project Description: %s

User Story: %s

Output ONLY the tasks as bullet points, no other text.
`,
        projectContext.Language,
        projectContext.Framework,
        projectContext.Role,
        projectContext.ProjectDescription,
        userStory,
    )

    resp, err := g.client.CreateChatCompletion(
        ctx,
        openai.ChatCompletionRequest{
            Model: openai.GPT4, // or whatever GPT-4 model you have access to
            Messages: []openai.ChatCompletionMessage{
                {
                    Role:    openai.ChatMessageRoleSystem,
                    Content: "You are a helpful assistant.",
                },
                {
                    Role:    openai.ChatMessageRoleUser,
                    Content: prompt,
                },
            },
            Temperature: 0.7,
        },
    )
    if err != nil {
        return nil, err
    }

    generated := resp.Choices[0].Message.Content
    // (We’ll want to parse out each bullet point into a slice of task descriptions)
    tasks := parseTasksFromGPTOutput(generated)

    return tasks, nil
}

// parseTasksFromGPTOutput is a naive bullet point parser
func parseTasksFromGPTOutput(output string) []string {
    lines := strings.Split(output, "\n")
    tasks := []string{}

    for _, line := range lines {
        trimmed := strings.TrimSpace(line)

        // Skip empty lines
        if trimmed == "" {
            continue
        }

        // Remove leading bullet characters (like "-", "*", or "•")
        // so the final string is just the task description
        trimmed = strings.TrimPrefix(trimmed, "- ")
        trimmed = strings.TrimPrefix(trimmed, "* ")
        trimmed = strings.TrimPrefix(trimmed, "• ")

        // Only add to tasks if there's something left
        if trimmed != "" {
            tasks = append(tasks, trimmed)
        }
    }

    return tasks
}