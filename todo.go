package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("File not found, creating one...")
		all_todos = []todo{}
		write_todo_json()
	}
	json.Unmarshal(content, &all_todos)
}

func write_todo_json() {
	data, _ := json.Marshal(all_todos)
	os.WriteFile("todo.json", data, 0644)
}

func parse_sts(status string) (Status, int) {
	var sts Status
	var def int = 0
	switch strings.ToLower(status) {
	case "completed":
		sts = StatusCompleted
	case "wip":
		sts = StatusInProgress
	case "pending":
		sts = StatusPending
	default:
		sts = StatusPending
		def = 1
	}
	return sts, def
}

func parse_sts_code(status Status) string {
	switch status {
	case StatusPending:
		return " Pending "
	case StatusInProgress:
		return "   WIP   "
	case StatusCompleted:
		return "Completed"
	default:
		return " Pending "
	}
}

func check_task_id_length(task_id string) {
	if len(task_id) < 4 {
		fmt.Println("The ID prefix should be greater than 4 characters.")
		return
	}
}

func main() {
	load_all_todo()
	if len(os.Args) < 2 {
		print_help()
		return
	}
	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			print_help()
			return
		}
		task := os.Args[2]
		add_task(task)
	case "list":
		list_mode := "all"
		if len(os.Args) == 3 {
			list_mode = os.Args[2]
		}
		list_task(list_mode)
	case "mark":
		if len(os.Args) < 4 {
			print_help()
			return
		}
		task_id := os.Args[2]
		new_status := strings.ToLower(os.Args[3])
		if new_status != "pending" && new_status != "wip" && new_status != "completed" {
			print_help()
			return
		}
		check_task_id_length(task_id)
		mark_task(new_status, task_id)

	case "del", "delete":
		if len(os.Args) < 3 {
			print_help()
			return
		}
		task_id := os.Args[2]
		check_task_id_length(task_id)
		remove_task(task_id)
	case "help":
		print_help()
	default:
		fmt.Println("Unknown command...\nAccess help via: todo help")
	}
}

func add_task(task string) {
	load_all_todo()
	id := uuid.New()
	var t = todo{Id: id, Task: task, TaskStatus: StatusPending}
	all_todos = append(all_todos, t)
	write_todo_json()
	fmt.Println("Task: ", task, "\nAdded successfully!")
}

func list_task(filter string) {
	load_all_todo()
	fmt.Println()
	filter_status, def := parse_sts(filter)

	for _, t := range all_todos {
		if def != 1 && t.TaskStatus != filter_status {
			continue
		}
		color := ""
		switch t.TaskStatus {
		case StatusPending:
			color = Yellow
		case StatusInProgress:
			color = Blue
		case StatusCompleted:
			color = Green
		}
		fmt.Println(Purple+t.Id.String()[:5]+Reset, "[", color+parse_sts_code(t.TaskStatus)+Reset, "]", "-", t.Task)
	}
	fmt.Println()
}

func mark_task(status string, taskID string) {
	load_all_todo()
	var sts Status
	sts, _ = parse_sts(status)
	for i := range all_todos {
		if strings.HasPrefix(all_todos[i].Id.String(), taskID) {
			all_todos[i].TaskStatus = sts
			fmt.Println("Status for Task: ", all_todos[i].Id, " updated to ", strings.ToUpper(parse_sts_code(all_todos[i].TaskStatus)))
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

func print_help() {
	fmt.Println(`
TODO CLI - Usage Guide

Commands:

  add <task>
      Add a new todo task
      Example: todo add "Learn Go"

  list [all|pending|wip|completed]
      List all tasks or filter by status
      Example: todo list
      Example: todo list pending

  mark <id_prefix> <pending|wip|completed>
      Update status of a task
      Example: todo mark a1b2 wip

  del <id_prefix>
      Delete a task
      Example: todo del a1b2

  help
      Show this help menu`)
}
