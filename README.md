# go-todo-cli

A simple console todo list written in **Go**. A learning project covering command-line arguments, structs, slices, JSON and file I/O.

## Features

- add, list, complete and delete tasks
- clear all completed tasks
- tasks are saved to `tasks.json`
- Go standard library only, no third-party packages

## Usage

```bash
go run . add "Learn Go"
go run . add "Finish Django homework"
go run . list
go run . done 1
go run . delete 2
go run . clear
```

Or build a binary:

```bash
go build -o todo.exe
./todo.exe list
```

## Example

```
$ todo list
  1. [x] Learn Go  (14.11 20:30)
  2. [ ] Finish Django homework  (14.11 20:31)

Done: 1 of 2
```

## What it covers

| Topic | Where in the code |
|---|---|
| structs and JSON tags | `type Task struct` |
| file I/O | `os.ReadFile`, `os.WriteFile` |
| JSON | `json.Unmarshal`, `json.MarshalIndent` |
| slices | `append`, removing an element |
| error handling | `error`, `fmt.Errorf` |
| command-line arguments | `os.Args`, `switch` |
