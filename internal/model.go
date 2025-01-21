package internal

type ProjectContext struct {
    Language   string `json:"language,omitempty"`
    Framework  string `json:"framework,omitempty"`
    Role       string `json:"role,omitempty"`
	ProjectDescription string `json:"project_description,omitempty"`
}

type Task struct {
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}

type UserStory struct {
    ID       int    `json:"id"`
    Story    string `json:"story"`
    Tasks    []Task `json:"tasks"`
    Complete bool   `json:"complete"`
}

type AppData struct {
    ProjectContext ProjectContext `json:"project_context"`
    UserStories    []UserStory    `json:"user_stories"`
    NextStoryID    int            `json:"next_story_id"`
}