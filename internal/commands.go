package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// For user input reading
var reader = bufio.NewReader(os.Stdin)

func SetContext(data *AppData) error {
    fmt.Print("Enter programming language: ")
    lang, _ := reader.ReadString('\n')
    data.ProjectContext.Language = strings.TrimSpace(lang)

    fmt.Print("Enter framework: ")
    fw, _ := reader.ReadString('\n')
    data.ProjectContext.Framework = strings.TrimSpace(fw)

    fmt.Print("Enter your role: ")
    role, _ := reader.ReadString('\n')
    data.ProjectContext.Role = strings.TrimSpace(role)

	fmt.Print("Enter project description: ")
	description, _ := reader.ReadString('\n')
	data.ProjectContext.ProjectDescription = strings.TrimSpace(description)

    // Save to JSON
    return SaveData(data)
}

func NewStory(data *AppData, gpt *GPTClient) error {
    fmt.Print("Enter user story: ")
    story, _ := reader.ReadString('\n')
    story = strings.TrimSpace(story)

    tasks, err := gpt.GenerateTasks(context.Background(), data.ProjectContext, story)
    if err != nil {
        return err
    }

    userStory := UserStory{
        ID:    data.NextStoryID,
        Story: story,
        Tasks: make([]Task, len(tasks)),
    }

    for i, t := range tasks {
        userStory.Tasks[i] = Task{
            Description: t,
            Completed:   false,
        }
    }

    data.UserStories = append(data.UserStories, userStory)
    data.NextStoryID++

    return SaveData(data)
}

func NextTask(data *AppData) {
    story := getActiveStory(data)
    if story == nil {
        fmt.Println("No active user story. Add a new story first!")
        return
    }

    for _, task := range story.Tasks {
        if !task.Completed {
            fmt.Println("Next task:", task.Description)
            return
        }
    }
    fmt.Println("All tasks in this story are complete!")
}

func CompleteTask(data *AppData, taskIndexStr string) error {
    story := getActiveStory(data)
    if story == nil {
        fmt.Println("No active user story. Add a new story first!")
        return nil
    }

	index := 0
	if taskIndexStr == "" {
		// get the first incomplete task
		for i, task := range story.Tasks {
			if !task.Completed {
				index = i
				break
			}
		}
	} else {
		index, err := strconv.Atoi(taskIndexStr)
		if err != nil || index < 1 || index > len(story.Tasks) {
			fmt.Println("Invalid task index. Must be between 1 and", len(story.Tasks))
			return nil
		}
	}


    // Mark task as complete
    story.Tasks[index].Completed = true

	// show which task was completed
	fmt.Printf("\033[1;32mTask #%d completed: %s\033[0m\n", index+1, story.Tasks[index].Description)

    // Check if story is complete
    if allTasksCompleted(story.Tasks) {
        story.Complete = true
        fmt.Printf("\033[1;32mUser story #%d is now complete!\033[0m\n", story.ID)
    }

    return SaveData(data)
}

func ShowStatus(data *AppData) {
    story := getActiveStory(data)
    if story == nil {
        fmt.Println("No active user story. Add a new story first!")
        return
    }

    fmt.Printf("Current story (\033[1;36mID: %d\033[0m): %s\n\n", story.ID, story.Story)
    // Show next incomplete task
    for i, task := range story.Tasks {
        if !task.Completed {
            fmt.Printf("Next task (\033[1;36m#%d\033[0m): %s\n", i+1, task.Description)
            return
        }
    }
    fmt.Println("\033[1;32mAll tasks are complete for this story.\033[0m")
}

func Help() {
    fmt.Println("\033[1;36mAvailable commands:\033[0m")
    fmt.Println("  set-context                - Set project context (language, framework, role)")
    fmt.Println("  new-story                  - Add a new user story and generate tasks")
    fmt.Println("  next-task                  - Show the next incomplete task for the active story")
    fmt.Println("  complete-task [taskIndex?]  - Mark the current or a specific task as complete")
    fmt.Println("  status                - Show the current user story and next task")
    fmt.Println("  help                       - Show this help menu")
}

// Helper: returns the first incomplete user story, or nil if none
func getActiveStory(data *AppData) *UserStory {
    for i := range data.UserStories {
        if !data.UserStories[i].Complete {
            return &data.UserStories[i]
        }
    }
    return nil
}

func allTasksCompleted(tasks []Task) bool {
    for _, t := range tasks {
        if !t.Completed {
            return false
        }
    }
    return true
}