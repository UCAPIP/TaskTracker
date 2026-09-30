package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Структура задачи
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

// Статусы задач
const (
	StatusDone       = "done"
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
)

// Путь JSON
var jsonPath = "Tasks.json"

// Загрузка задач из файла
func loadTasks(path string) ([]Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Файл не существует")
			return []Task{}, nil
		}
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if len(data) == 0 {
		fmt.Println("Задач нет. Чтобы создать новую задачу введите команду: task-cli add [текст задачи]")
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	return tasks, nil
}

// Добавление новой задачи
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
		fmt.Println("ошибка сериализации:", err)
		return
	}

	err = os.WriteFile(jsonPath, newData, 0644)
	if err != nil {
		fmt.Println("ошибка записи в JSON:", err)
		return
	}

	fmt.Println("Таска добавлена, файл обновлён")

}

// Обновление текста(description) задачи
func update(data *[]Task, id int, description string) {

	for i, task := range *data {
		if task.Id == id {
			fmt.Println("Старый текст задачи:", (*data)[i].Description)
			fmt.Println("Новый текст задачи: ", description)
			(*data)[i].UpdateDescription(description)

			newData, err := json.MarshalIndent(*data, "", "  ")
			if err != nil {
				fmt.Println("ошибка сериализации:", err)
				return
			}

			err = os.WriteFile(jsonPath, newData, 0644)
			if err != nil {
				fmt.Println("ошибка записи в JSON:", err)
				return
			}

			return

		}
	}

	fmt.Println("Задача с id", id, "не найдена")

}

// Смена статуса задачи (done\in-progress)
func markStatus(data *[]Task, id int, status string) {
	for i, task := range *data {
		if task.Id == id {
			(*data)[i].SetStatus(status)

			newData, err := json.MarshalIndent(*data, "", "  ")
			if err != nil {
				fmt.Println("ошибка сериализации:", err)
				return
			}

			err = os.WriteFile(jsonPath, newData, 0644)
			if err != nil {
				fmt.Println("ошибка записи в JSON:", err)
				return
			}

			return

		}
	}

	fmt.Println("Задача с id", id, "не найдена")

}

// Удаление задачи
func deleteTask(data *[]Task, id int) {

	for i, task := range *data {
		if task.Id == id {
			*data = append((*data)[:i], (*data)[i+1:]...)

			newData, err := json.MarshalIndent(*data, "", "  ")
			if err != nil {
				fmt.Println("ошибка сериализации:", err)
				return
			}

			err = os.WriteFile(jsonPath, newData, 0644)
			if err != nil {
				fmt.Println("ошибка записи в JSON:", err)
				return
			}

			fmt.Println("Задача", id, "удалена")

			return

		}
	}

	fmt.Println("Задача с id", id, "не найдена")

}

// Вывод списка задач, доступена филтрация по статусу(all\done\todo\in-progress). По дефолту выводим all
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
			fmt.Fprintf(os.Stderr, "Неизвестный фильтр: %s\n", filter)
			os.Exit(1)
		}
	}
}

// Функция для форматирования вывода задач на печать
func printTask(t Task) {
	fmt.Printf("ID: %d | %s | %s | создано: %s | обновлено: %s\n",
		t.Id, t.Description, t.Status, t.CreatedAt, t.UpdatedAt)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Использование: task-cli <command> [args]")
		fmt.Println("Команды: add, update, delete, mark-done, mark-in-progress, list")
		os.Exit(1)
	}

	tasks, err := loadTasks(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки задач: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Использование: task-cli add <description>")
			os.Exit(1)
		}
		description := strings.Join(os.Args[2:], " ")
		add(&tasks, description)

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Использование: task-cli update <id> <description>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Некорректный ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		description := strings.Join(os.Args[3:], " ")
		update(&tasks, id, description)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Использование: task-cli delete <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Некорректный ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		deleteTask(&tasks, id)

	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Использование: task-cli mark-done <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Некорректный ID: %s\n", os.Args[2])
			os.Exit(1)
		}
		markStatus(&tasks, id, StatusDone)

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Использование: task-cli mark-in-progress <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Некорректный ID: %s\n", os.Args[2])
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
		fmt.Fprintf(os.Stderr, "Неизвестная команда: %s\n", os.Args[1])
		os.Exit(1)
	}

}
