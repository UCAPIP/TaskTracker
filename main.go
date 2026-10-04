package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Structure of the Task
type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func (t *Task) Touch() {
	t.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
}

func (t *Task) SetStatus(status string) {
	t.Status = status
	t.Touch()
}

func (t *Task) UpdateDescription(description string) {
	t.Description = description
	t.Touch()
}

// Task Statuses
const (
	StatusDone       = "done"
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
)

// JSON Path
var jsonPath = "Tasks.json"

// Loading Tasks from a File
func loadTasks(path string) ([]Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The file does not exist")
			return []Task{}, nil
		}
		return nil, fmt.Errorf("File read error: %w", err)
	}

	if len(data) == 0 {
		fmt.Println("There are no tasks. To create a new task, enter the command: task-cli add [task text]")
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w", err)
	}
	return tasks, nil
}

// Adding a New Task
func add(data *[]Task, description string) {
	maxID := func(tasks []Task) int {
		maxId := 0
		for _, t := range tasks {
			if t.Id > maxId {
				maxId = t.Id
			}
		}
		return maxId
	}(*data)

	newTask := Task{
		Id:          maxID + 1,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	*data = append(*data, newTask)

	newData, err := json.MarshalIndent(*data, "", "  ")
	if err != nil {
		fmt.Println("serialization error:", err)
		return
	}

	err = os.WriteFile(jsonPath, newData, 0644)
	if err != nil {
		fmt.Println("JSON writing error:", err)
		return
	}

	fmt.Println("The task has been added, and the file has been updated")

}

// Updating the task description
func update(data *[]Task, id int, description string) {

	for i, task := range *data {
		if task.Id == id {
			fmt.Println("Old task text:", (*data)[i].Description)
			fmt.Println("New task text: ", description)
			(*data)[i].UpdateDescription(description)

			newData, err := json.MarshalIndent(*data, "", "  ")
			if err != nil {
				fmt.Println("serialization error:", err)
				return
			}

			err = os.WriteFile(jsonPath, newData, 0644)
			if err != nil {
				fmt.Println("JSON writing error:", err)
				return
			}

			return

		}
	}

	fmt.Println("The task with id", id, "was not found")

}

// Changing a task's status (done/in progress)
func markStatus(data *[]Task, id int, status string) {
	for i, task := range *data {
		if task.Id == id {
			(*data)[i].SetStatus(status)

			newData, err := json.MarshalIndent(*data, "", "  ")
			if err != nil {
				fmt.Println("serialization error:", err)
				return
			}

			err = os.WriteFile(jsonPath, newData, 0644)
			if err != nil {
				fmt.Println("JSON writing error:", err)
				return
			}

			return

		}
	}

	fmt.Println("The task with id", id, "was not found")

}

// Deleting a Task
func deleteTask(data *[]Task, id int) {

	for i, task := range *data {
		if task.Id == id {
			*data = append((*data)[:i], (*data)[i+1:]...)

			newData, err := json.MarshalIndent(*data, "", "  ")
			if err != nil {
				fmt.Println("serialization error:", err)
				return
			}

			err = os.WriteFile(jsonPath, newData, 0644)
			if err != nil {
				fmt.Println("JSON writing error:", err)
				return
			}

			fmt.Println("The task", id, "was deleted")

			return

		}
	}

	fmt.Println("The task with id", id, "was not found")

}

// Display a list of tasks; filtering by status (all, done, todo, in-progress) is available. By default, “all” is displayed.
func listTasks(tasks []Task, filter string) {
	for _, t := range tasks {
		switch filter {
		case "all":
			printTask(t)
		case "done":
			if t.Status == StatusDone {
				printTask(t)
			}
		case "todo":
			if t.Status == StatusTodo {
				printTask(t)
			}
		case "in-progress":
			if t.Status == StatusInProgress {
				printTask(t)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown filter: %s\n", filter)
			os.Exit(1)
		}
	}
}

// A function for formatting the output of print Tasks
func printTask(t Task) {
	fmt.Printf("ID: %d | %s | %s | created: %s | updated: %s\n",
		t.Id, t.Description, t.Status, t.CreatedAt, t.UpdatedAt)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: task-cli <command> [args]")
		fmt.Println("Commands: add, update, delete, mark-done, mark-in-progress, list")
		os.Exit(1)
	}

	tasks, err := loadTasks(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add <description>")
			os.Exit(1)
		}
		description := strings.Join(os.Args[2:], " ")
		add(&tasks, description)

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: task-cli update <id> <description>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		description := strings.Join(os.Args[3:], " ")
		update(&tasks, id, description)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli delete <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		deleteTask(&tasks, id)

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli mark-done <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		markStatus(&tasks, id, StatusDone)

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli mark-in-progress <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		markStatus(&tasks, id, StatusInProgress)

	case "list":
		filter := "all"
		if len(os.Args) >= 3 {
			filter = os.Args[2]
		}
		listTasks(tasks, filter)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

}
