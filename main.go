// go-todo-cli — простий консольний список справ на Go.
// Дані зберігаються у файлі tasks.json поруч із програмою.
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

// Task — одна справа у списку.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// loadTasks читає список справ із файлу (якщо файлу ще немає — повертає порожній список).
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

// saveTasks записує список справ у файл у гарному (з відступами) JSON.
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

// nextID повертає наступний вільний номер справи.
func nextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// findIndex шукає справу за номером і повертає її індекс у слайсі (або -1).
func findIndex(tasks []Task, id int) int {
	for i, t := range tasks {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func printHelp() {
	fmt.Println(`go-todo-cli — список справ у терміналі

Команди:
  add <текст>     додати справу
  list            показати всі справи
  done <номер>    позначити справу виконаною
  delete <номер>  видалити справу
  clear           видалити всі виконані справи
  help            показати цю довідку`)
}

func printList(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("Список порожній. Додай справу: todo add \"Купити хліб\"")
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
	fmt.Printf("\nВиконано: %d з %d\n", done, len(tasks))
}

// parseID перетворює аргумент командного рядка на номер справи.
func parseID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("вкажи номер справи")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("номер має бути числом: %q", args[0])
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
		return fmt.Errorf("не вдалося прочитати %s: %w", dataFile, err)
	}

	command, rest := args[0], args[1:]
	switch command {
	case "add":
		title := strings.TrimSpace(strings.Join(rest, " "))
		if title == "" {
			return fmt.Errorf("напиши текст справи: todo add \"Купити хліб\"")
		}
		task := Task{ID: nextID(tasks), Title: title, CreatedAt: time.Now()}
		tasks = append(tasks, task)
		fmt.Printf("Додано справу #%d: %s\n", task.ID, task.Title)

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
			return fmt.Errorf("справу #%d не знайдено", id)
		}
		tasks[i].Done = true
		fmt.Printf("Справу #%d виконано: %s\n", id, tasks[i].Title)

	case "delete":
		id, err := parseID(rest)
		if err != nil {
			return err
		}
		i := findIndex(tasks, id)
		if i == -1 {
			return fmt.Errorf("справу #%d не знайдено", id)
		}
		fmt.Printf("Видалено справу #%d: %s\n", id, tasks[i].Title)
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
		fmt.Printf("Видалено виконаних справ: %d\n", removed)

	case "help", "-h", "--help":
		printHelp()
		return nil

	default:
		return fmt.Errorf("невідома команда %q (дивись: todo help)", command)
	}

	return saveTasks(tasks)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Помилка:", err)
		os.Exit(1)
	}
}
