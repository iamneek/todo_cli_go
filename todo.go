package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type Status int

const (
	StatusPending Status = iota
	StatusInProgress
	StatusCompleted
)

type todo struct {
	Id         uuid.UUID `json:"id"`
	Task       string    `json:"todo"`
	TaskStatus Status    `json:"status"`
}

var all_todos []todo

func load_all_todo() {
	content, err := os.ReadFile("todo.json")
	if err != nil {
		log.Fatal("Failed to read file.")
		return
	}
	json.Unmarshal(content, &all_todos)
}

func main() {
	load_all_todo()
	command := os.Args[1]
	switch command {
	case "add":
		task := os.Args[2]
		add_task(task)
	case "list":
		list_task()
	case "mark":
		task_id := os.Args[2]
		fmt.Println(task_id)
	default:
		fmt.Println("Unknown command")
	}
}

func add_task(task string) {
	load_all_todo()
	id := uuid.New()
	var t = todo{Id: id, Task: task, TaskStatus: StatusPending}
	all_todos = append(all_todos, t)
	data, _ := json.Marshal(all_todos)
	os.WriteFile("todo.json", data, 0644)
}

func list_task() {
	load_all_todo()
	for i := range all_todos {
		fmt.Println("Task ID: ", all_todos[i].Id, "\nTask: ", all_todos[i].Task, "\nTask Status: ", all_todos[i].TaskStatus)
	}
}
