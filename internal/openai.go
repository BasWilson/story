package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type GPTClient struct {
	apiKey string
	model  string
    baseUrl string
}

func NewGPTClient(apiKey string, model string, baseUrl string) *GPTClient {
	return &GPTClient{
		apiKey: apiKey,
		model:  model,
		baseUrl: baseUrl,
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

	requestBody := map[string]interface{}{
		"model": g.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/chat/completions", g.baseUrl), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	generated := result.Choices[0].Message.Content
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