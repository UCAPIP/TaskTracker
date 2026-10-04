# Task Tracker CLI

A simple console application for managing a task list. It allows adding, updating, and deleting tasks, changing their status, and viewing the list with filtering. All data is stored in a local `Tasks.json` file.

The project is implemented in **Go** without using any external libraries — only the standard library.

---

## ✨ Features

- ➕ Add tasks
- ✏️ Update task description
- 🗑 Delete a task
- 🔄 Change status: `todo` → `in-progress` → `done`
- 📋 View all tasks
- 🔍 Filter tasks by status
- 💾 Automatic saving to `Tasks.json`
- 🕐 Automatic tracking of creation and update times

---

## 📋 Requirements

- [Go](https://go.dev/dl/) version **1.21** or newer

Check your version:
```bash
go version
```

---

## 🚀 Installation and Build

1. Clone the repository:
```bash
git clone <repository-url>
cd task-tracker-cli
```

2. Build the binary file:
```bash
go build -o task-cli main.go
```

3. Run commands using the built file:
```bash
./task-cli <command> [args]
```

> **For Windows**, use `task-cli.exe` instead of `./task-cli`.

Alternatively, you can run it without building:
```bash
go run main.go <command> [args]
```

---

## 📖 Usage

### Add a task

```bash
task-cli add "Buy groceries"
```

A description with spaces can be passed as a single argument in quotes or as multiple arguments — they will be concatenated:

```bash
task-cli add "Buy groceries and cook dinner"
task-cli add Buy groceries and cook dinner
```

### Update a task description

```bash
task-cli update 1 "Buy groceries, cook dinner, and wash the dishes"
```

### Delete a task

```bash
task-cli delete 1
```

### Mark a task as done

```bash
task-cli mark-done 1
```

### Mark a task as in progress

```bash
task-cli mark-in-progress 1
```

### Show all tasks

```bash
task-cli list
```

### Show tasks by status

```bash
task-cli list done         # only completed tasks
task-cli list todo         # only planned tasks
task-cli list in-progress  # only tasks in progress
```

---

## 📊 Example Session

```bash
$ task-cli add "Learn about pointers in Go"
Task added, file updated

$ task-cli add "Write README"
Task added, file updated

$ task-cli list
ID: 1 | Learn about pointers in Go | todo | created: 2026-09-18 10:00:00 | updated: 2026-09-18 10:00:00
ID: 2 | Write README | todo | created: 2026-09-18 10:01:00 | updated: 2026-09-18 10:01:00

$ task-cli mark-done 1
Status of task 1 changed to done

$ task-cli list done
ID: 1 | Learn about pointers in Go | done | created: 2026-09-18 10:00:00 | updated: 2026-09-18 10:02:00
```

---

## 🗂 Data Structure

Each task is stored as an object with the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `id` | `int` | Unique task identifier |
| `description` | `string` | Task text |
| `status` | `string` | Status: `todo`, `in-progress`, `done` |
| `createdAt` | `string` | Creation date and time (`2006-01-02 15:04:05`) |
| `updatedAt` | `string` | Date and time of the last update |

### Example `Tasks.json`

```json
[
  {
    "id": 1,
    "description": "Learn about pointers in Go",
    "status": "done",
    "createdAt": "2026-09-18 10:00:00",
    "updatedAt": "2026-09-18 10:02:00"
  },
  {
    "id": 2,
    "description": "Write README",
    "status": "todo",
    "createdAt": "2026-09-18 10:01:00",
    "updatedAt": "2026-09-18 10:01:00"
  }
]
```

---

## 🎯 Available Statuses

| Status | Meaning | Description |
|--------|---------|-------------|
| `todo` | Default | Task is created but not yet started |
| `in-progress` | In progress | Task is currently being worked on |
| `done` | Completed | Task is finished |

---

## ⚠️ Error Handling

The program correctly handles the following situations:

- **File does not exist** — created automatically on the first run.
- **Invalid ID** — an error message is displayed.
- **Task not found** — displays the message `Task with id X not found`.
- **Unknown command** — displays a list of available commands.
- **Unknown filter** — displays an error message.

All errors are written to `stderr`, and the program exits with code `1`.

---

## 🛠 Technologies

- Language: **Go**
- Data storage: **JSON file** (`Tasks.json`)
- Only standard library packages are used:
  - `encoding/json` — serialization/deserialization
  - `os` — file and command-line argument handling
  - `strconv` — string to number conversion
  - `strings` — string manipulation
  - `time` — date and time handling

---

## 📄 License

This project is created for educational purposes.
