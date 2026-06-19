package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

// TODO: Handle creation of todo.json in case it doesnot exist

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
		list_mode := "all"
		if len(os.Args) == 3 {
			list_mode = os.Args[2]
		}
		list_task(list_mode)
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

		fmt.Println(
			Purple+t.Id.String()[:5]+Reset,
			"[",
			color+parse_sts_code(t.TaskStatus)+Reset,
			"]",
			"-",
			t.Task,
		)
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
