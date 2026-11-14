// go-todo-cli is a simple console todo list written in Go.
// Tasks are stored in the tasks.json file next to the program.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const dataFile = "tasks.json"

// Task is a single item in the todo list.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// loadTasks reads the task list from the file (returns an empty list if the file does not exist yet).
func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(dataFile)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

// saveTasks writes the task list to the file as pretty-printed JSON.
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

// nextID returns the next free task number.
func nextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// findIndex looks up a task by its number and returns its index in the slice (or -1).
func findIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func printHelp() {
	fmt.Println(`go-todo-cli - a todo list in your terminal

Commands:
  add <text>      add a task
  list            show all tasks
  done <number>   mark a task as done
  delete <number> delete a task
  clear           remove all completed tasks
  help            show this help`)
}

func printList(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("The list is empty. Add a task: todo add \"Buy bread\"")
		return
	}
	done := 0
	for _, t := range tasks {
		mark := "[ ]"
		if t.Done {
			mark = "[x]"
			done++
		}
		fmt.Printf("%3d. %s %s  (%s)\n", t.ID, mark, t.Title, t.CreatedAt.Format("02.01 15:04"))
	}
	fmt.Printf("\nDone: %d of %d\n", done, len(tasks))
}

// parseID converts a command-line argument into a task number.
func parseID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("please specify a task number")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("task number must be a number: %q", args[0])
	}
	return id, nil
}

func run(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}

	tasks, err := loadTasks()
	if err != nil {
		return fmt.Errorf("could not read %s: %w", dataFile, err)
	}

	command, rest := args[0], args[1:]
	switch command {
	case "add":
		title := strings.TrimSpace(strings.Join(rest, " "))
		if title == "" {
			return fmt.Errorf("please write the task text: todo add \"Buy bread\"")
		}
		task := Task{ID: nextID(tasks), Title: title, CreatedAt: time.Now()}
		tasks = append(tasks, task)
		fmt.Printf("Added task #%d: %s\n", task.ID, task.Title)

	case "list":
		printList(tasks)
		return nil

	case "done":
		id, err := parseID(rest)
		if err != nil {
			return err
		}
		i := findIndex(tasks, id)
		if i == -1 {
			return fmt.Errorf("task #%d not found", id)
		}
		tasks[i].Done = true
		fmt.Printf("Task #%d is done: %s\n", id, tasks[i].Title)

	case "delete":
		id, err := parseID(rest)
		if err != nil {
			return err
		}
		i := findIndex(tasks, id)
		if i == -1 {
			return fmt.Errorf("task #%d not found", id)
		}
		fmt.Printf("Deleted task #%d: %s\n", id, tasks[i].Title)
		tasks = append(tasks[:i], tasks[i+1:]...)

	case "clear":
		left := tasks[:0]
		removed := 0
		for _, t := range tasks {
			if t.Done {
				removed++
				continue
			}
			left = append(left, t)
		}
		tasks = left
		fmt.Printf("Removed completed tasks: %d\n", removed)

	case "help", "-h", "--help":
		printHelp()
		return nil

	default:
		return fmt.Errorf("unknown command %q (see: todo help)", command)
	}

	return saveTasks(tasks)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
