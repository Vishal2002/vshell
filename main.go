package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

// history lives at package level
var history []string

func executeCommand(command string) error {
	cmdLine := strings.TrimSpace(command)
	if cmdLine == "" {
		return fmt.Errorf("no command provided")
	}

	// push into history
	history = append(history, cmdLine)

	parts := strings.Fields(cmdLine)
	cmd, args := parts[0], parts[1:]

	switch cmd {
	case "exit":
		fmt.Println("Exiting the shell")
		os.Exit(0)

	case "ls":
		return listDir()

	case "pwd":
		return printPwd()

	case "cd":
		return changeDir(args)

	case "clear":
		fmt.Print("\033[H\033[2J")
		return nil

	case "help":
		printHelp()
		return nil

	case "date":
		return printDate(args)

	case "whoami":
		return printUser()

	case "history":
		for i, h := range history {
			fmt.Printf("%3d %s\n", i+1, h)
		}
		return nil

	case "mkdir":
		if len(args) < 1 {
			return fmt.Errorf("usage: mkdir <directory>")
		}
		if err := os.Mkdir(args[0], 0755); err != nil {
			return err
		}
		fmt.Println("Directory created:", args[0])
		return nil

	default:
		return fmt.Errorf("unknown command %q (type \"help\" for list)", cmd)
	}
	return nil
}

// helpers --------------------------------------------------------------------
func changeDir(args []string) error {
	var targetDir string

	// Handle no arguments (default to home directory)
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %v", err)
		}
		targetDir = home
	} else {
		targetDir = args[0]

		// Handle special cases
		if targetDir == "~" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %v", err)
			}
			targetDir = home
		} else if targetDir == "-" {
			// Previous directory (stored in environment variable OLDPWD)
			previousDir := os.Getenv("OLDPWD")
			if previousDir == "" {
				return fmt.Errorf("no previous directory set")
			}
			targetDir = previousDir
		}
	}

	// Get the current directory before changing (for OLDPWD)
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Change to the target directory
	if err := os.Chdir(targetDir); err != nil {
		return fmt.Errorf("failed to change directory to %s: %v", targetDir, err)
	}

	// Update OLDPWD environment variable
	if err := os.Setenv("OLDPWD", currentDir); err != nil {
		return fmt.Errorf("failed to set OLDPWD: %v", err)
	}

	return nil
}

func listDir() error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}
	for _, e := range entries {
		fmt.Println(e.Name())
	}
	return nil
}

func printPwd() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(cwd)
	return nil
}

func printDate(args []string) error {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return err
	}
	now := time.Now().In(loc)

	if len(args) == 0 {
		fmt.Println(now.Format("Mon Jan 02 15:04:05 IST 2006"))
		return nil
	}

	switch args[0] {
	case "year":
		fmt.Println(now.Year())
	case "month":
		fmt.Println(now.Month())
	case "day":
		fmt.Println(now.Day())
	case "time":
		fmt.Println(now.Format("15:04:05"))
	default:
		return fmt.Errorf("usage: date [year|month|day|time]")
	}
	return nil
}

func printUser() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	fmt.Println(u.Username)
	return nil
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  ls        - List directory contents")
	fmt.Println("  cd [dir]  - Change directory")
	fmt.Println("  pwd       - Print working directory")
	fmt.Println("  clear     - Clean the terminal")
	fmt.Println("  date ...  - Show date/time (sub: year|month|day|time)")
	fmt.Println("  whoami    - Show current user name")
	fmt.Println("  history   - Show command history")
	fmt.Println("  mkdir <d> - Create directory")
	fmt.Println("  exit      - Exit the shell")
}

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "vshell> ",
		HistoryFile:     os.Getenv("HOME") + "/.vshell_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "readline init:", err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt { // user hit Ctrl‑C
			if len(line) == 0 {
				continue // just redraw prompt
			}
			// clear the in‑progress line
			rl.Write([]byte{'\r'})
			continue
		}
		if err == io.EOF { // Ctrl‑D
			fmt.Println("bye")
			break
		}

		command := strings.TrimSpace(line)
		if command != "" {
			history = append(history, command) // in‑memory history for your `history` cmd
		}

		if err := executeCommand(command); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
