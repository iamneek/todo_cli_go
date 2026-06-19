# TODO CLI (Go Learning Project)

A simple command-line todo application built in Go to learn core language concepts and CLI design.

This project is focused on understanding Go fundamentals like slices, structs, file handling, CLI arguments, and JSON serialization.

---

## Features

* Add todo tasks
* List all tasks or filter by status
* Mark tasks as pending, in-progress, or completed
* Delete tasks by ID prefix
* Persistent storage using JSON file
* Colored terminal output for better readability

---

## Tech Stack

* Go (Golang)
* Standard Library (`os`, `fmt`, `encoding/json`, `strings`)
* UUID for unique task IDs

---

## Commands

### Add a task

```bash
todo add "Learn Go"
```

---

### List tasks

```bash
todo list
```

Filter by status:

```bash
todo list pending
todo list wip
todo list completed
```

---

### Mark task status

```bash
todo mark <id_prefix> <status>
```

Example:

```bash
todo mark a1b2 wip
```

Valid statuses:

* pending
* wip
* completed

---

### Delete a task

```bash
todo del <id_prefix>
```

or

```bash
todo delete <id_prefix>
```

---

### Help

```bash
todo help
```

Shows all available commands.

---

## Data Storage

Tasks are stored in a local JSON file:

```
todo.json
```

Each task contains:

* ID (UUID)
* Task description
* Status

---

## Learning Goals

This project was built to understand:

* CLI argument parsing in Go (`os.Args`)
* Structs and slices
* JSON encoding/decoding
* File I/O in Go
* Error handling patterns
* Basic CLI architecture design
* State persistence
* Input validation
* Simple terminal UI formatting

---

## Notes

This is a learning project, not production-ready software. It intentionally avoids external frameworks to focus on core Go concepts.

---