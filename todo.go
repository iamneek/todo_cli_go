package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

func write_todo_json() {
	data, _ := json.Marshal(all_todos)
	os.WriteFile("todo.json", data, 0644)
}

func parse_sts(status string) Status {
	var sts Status
	switch strings.ToLower(status) {
	case "completed":
		sts = StatusCompleted
	case "wip":
		sts = StatusInProgress
	case "pending":
		sts = StatusPending
	default:
		sts = StatusPending
	}
	return sts
}

func check_task_id_length(task_id string) {
	if len(task_id) < 4 {
		fmt.Println("The ID prefix should be greater than 4 characters.")
		os.Exit(1)
		return
	}
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
		new_status := os.Args[3]
		check_task_id_length(task_id)
		mark_task(new_status, task_id)

	case "del":
	case "delete":
		task_id := os.Args[2]
		check_task_id_length(task_id)
		remove_task(task_id)
	default:
		fmt.Println("Unknown command")
	}
}

func add_task(task string) {
	load_all_todo()
	id := uuid.New()
	var t = todo{Id: id, Task: task, TaskStatus: StatusPending}
	all_todos = append(all_todos, t)
	write_todo_json()
}

func list_task() {
	load_all_todo()
	for i := range all_todos {
		fmt.Println("Task ID: ", all_todos[i].Id, "\nTask: ", all_todos[i].Task, "\nTask Status: ", all_todos[i].TaskStatus)
	}
}

func mark_task(status string, taskID string) {
	load_all_todo()
	var sts Status
	sts = parse_sts(status)
	for i := range all_todos {
		if strings.HasPrefix(all_todos[i].Id.String(), taskID) {
			all_todos[i].TaskStatus = sts
			fmt.Println(all_todos[i])
		}
	}
	write_todo_json()
}

func remove_task(taskID string) {
	load_all_todo()
	for i := 0; i < len(all_todos); i++ {
		if strings.HasPrefix(all_todos[i].Id.String(), taskID) {
			all_todos = append(all_todos[:i], all_todos[i+1:]...)
			fmt.Println("Task Deleted")
		}
	}
	write_todo_json()
}
