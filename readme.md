# story - AI-Powered Task Management for Developers

story is a CLI tool that helps developers manage their tasks by generating actionable steps from user storys using AI (OpenAI GPT-4). It maintains project context and tracks progress through user storys and tasks.

## Features

-   🧠 AI-powered task generation from user storys
-   📝 Maintains project context (language, framework, role, description)
-   ✅ Tracks task completion status
-   💾 Persistent storage of project data per repository
-   🎯 Focus on the next actionable task

## Installation

1. Clone the repository:

```bash
git clone https://github.com/baswilson/story.git
```

2. Install dependencies:

3. Create a `.env` file with your OpenAI API key:

```env
OPENAI_API_KEY=your-api-key-here
```

4. Build the project:

```bash
go build -o story cmd/tool/main.go
```

5. Run the tool:

```bash
./story
```

6. Recommend (but optional) add to profile

Add the `OPEN_API_KEY` to your ~/.zshrc profile for use everywhere. Add an alias to the path of the executable in your ~/.zshrc profile for use everywhere like so:

```bash
nano ~/.zshrc
```

```bash
alias story="/Path/To/Your/Executable/story"
```

After that you can just run it in any folder with persistance by using `story`

## Usage

### Commands

| Command               | Description                                                           |
| --------------------- | --------------------------------------------------------------------- |
| `set-context`         | Set project context (language, framework, role, description)          |
| `new-story`           | Add a new user story and generate tasks                               |
| `next-task`           | Show the next incomplete task for the active story                    |
| `complete-task [idx]` | Mark the current or a specific task as complete (optional task index) |
| `show-status`         | Show the current user story and next task                             |
| `help`                | Show the help menu                                                    |

### Example Workflow

1. Set your project context:

```bash
./story set-context
```

2. Add a new user story:

```bash
./story new-story
```

3. Check your next task:

```bash
./story next-task
```

4. Get current task / story

```bash
./story
```

5. Complete current task

```bash
./story complete-task
```

6. Complete a specific task

```bash
./story complete-task 2
```

7. Check progress

```bash
./story show-status
```

## Data Storage

All project data is stored in `memory.json` in the project root. This includes:

-   Project context
-   User stories
-   Task status
-   Next story ID

## Requirements

-   Go 1.22.2 or higher
-   OpenAI API key (GPT-4 access recommended)

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT License
