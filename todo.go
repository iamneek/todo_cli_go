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

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
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

func parse_sts_code(status Status) string {
	switch status {
	case StatusPending:
		return "Pending  "
	case StatusInProgress:
		return "WIP      "
	case StatusCompleted:
		return "Completed"
	default:
		return "Pending  "
	}
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
	fmt.Println()
	for i := range all_todos {
		switch all_todos[i].TaskStatus {
		case StatusPending:
			fmt.Println(Purple+all_todos[i].Id.String()[:5]+Reset, "[", Yellow+parse_sts_code(all_todos[i].TaskStatus)+Reset, "]", "-", all_todos[i].Task)
		case StatusInProgress:
			fmt.Println(Purple+all_todos[i].Id.String()[:5]+Reset, "[", Blue+parse_sts_code(all_todos[i].TaskStatus)+Reset, "]", "-", all_todos[i].Task)
		case StatusCompleted:
			fmt.Println(Purple+all_todos[i].Id.String()[:5]+Reset, "[", Green+parse_sts_code(all_todos[i].TaskStatus)+Reset, "]", "-", all_todos[i].Task)

		}
	}
	fmt.Println()
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
