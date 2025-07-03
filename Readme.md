# vshell — A Simple Shell in Go 

> A minimal command-line shell written in Go while learning the language — built to explore Go's standard libraries and system capabilities.

---

## 🚀 Features

- `ls` — list directory contents  
- `pwd` — print current working directory  
- `mkdir <name>` — create a new directory  
- `clear` — clear the terminal screen  
- `whoami` — show the current user  
- `date` — show current IST date & time  
  - `date year`, `date month`, `date day`, `date time`  
- `history` — show command history (↑/↓ keys supported!)  
- `exit` — exit the shell  
- Auto-updating prompt  
- Arrow key navigation (via `readline`)  
- Built-in command history search with `Ctrl + R`

---

## 🧠 Why I Built This

I’m learning the Go programming language, and this project helped me explore:

- Working with files and directories in Go
- Using ANSI escape codes
- Handling user input
- Structuring command dispatch logic
- Exploring Go’s `os`, `io`, `time`, and `user` packages
- Using third-party packages like [`github.com/chzyer/readline`](https://github.com/chzyer/readline)

---

## 🛠️ Getting Started

### Prerequisites

- Go installed (≥ v1.20)  
  [Install Go →](https://go.dev/dl)

### Clone and Run

```bash
git clone https://github.com/Vishal2002/vshell.git
cd vshell
go run main.go
