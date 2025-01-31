package main

import (
	"fmt"
	"os"

	"github.com/baswilson/storie/internal"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	godotenv.Load()

	// 1. Load data from JSON
	data, err := internal.LoadData()
	if err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}

	// 2. Create GPT client
	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the AI_API_KEY environment variable.")
		os.Exit(1)
	}

	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "gpt-4o-2024-11-20"
	}

	baseUrl := os.Getenv("AI_API_BASE")
	if baseUrl == "" {
		baseUrl = "https://api.openai.com/v1"
	}

	gptClient := internal.NewGPTClient(apiKey, model, baseUrl)

	// 3. Check CLI args
	if len(os.Args) < 2 {
		// Show status by default, or show help
		internal.ShowStatus(data)
		fmt.Println()
		internal.Help()
		return
	}

	command := os.Args[1]

	switch command {
	case "set-context":
		err = internal.SetContext(data)
	case "new-story":
		err = internal.NewStory(data, gptClient)
		internal.ShowStatus(data)
	case "next-task":
		internal.NextTask(data)
	case "complete-task":
		taskIndex := ""
		if len(os.Args) >= 3 {
			taskIndex = os.Args[2]
		}
		err = internal.CompleteTask(data, taskIndex)
	case "status":
		internal.ShowStatus(data)
	case "help":
		internal.Help()
	case "config-info":
		internal.ConfigInfo()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		internal.Help()
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
}