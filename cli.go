package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Execute     func(args []string) error
}

type CLI struct {
	commands    map[string]Command
	history     []string
	running     bool
	workingDir  string
	maxHistory  int
	promptStyle string
}

func NewCLI() *CLI {
	cli := &CLI{
		commands:    make(map[string]Command),
		history:     make([]string, 0),
		running:     true,
		workingDir:  ".",
		maxHistory:  100,
		promptStyle: "default",
	}
	cli.registerCommands()
	return cli
}

func (c *CLI) registerCommands() {
	c.commands = map[string]Command{
		"help": {
			Name:        "help",
			Description: "Show available commands",
			Usage:       "help [command]",
			Execute:     c.cmdHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the CLI",
			Usage:       "exit",
			Execute:     c.cmdExit,
		},
		"clear": {
			Name:        "clear",
			Description: "Clear the screen",
			Usage:       "clear",
			Execute:     c.cmdClear,
		},
		"time": {
			Name:        "time",
			Description: "Show current time and date",
			Usage:       "time",
			Execute:     c.cmdTime,
		},
		"history": {
			Name:        "history",
			Description: "Show command history",
			Usage:       "history [count]",
			Execute:     c.cmdHistory,
		},
		"cd": {
			Name:        "cd",
			Description: "Change working directory",
			Usage:       "cd <directory>",
			Execute:     c.cmdCd,
		},
		"pwd": {
			Name:        "pwd",
			Description: "Print working directory",
			Usage:       "pwd",
			Execute:     c.cmdPwd,
		},
		"ls": {
			Name:        "ls",
			Description: "List directory contents",
			Usage:       "ls [directory]",
			Execute:     c.cmdLs,
		},
		"prompt": {
			Name:        "prompt",
			Description: "Change prompt style (default/full)",
			Usage:       "prompt <style>",
			Execute:     c.cmdPrompt,
		},
	}
}

func (c *CLI) cmdHelp(args []string) error {
	if len(args) > 0 {
		if cmd, ok := c.commands[args[0]]; ok {
			fmt.Printf("\nCommand: %s\n", cmd.Name)
			fmt.Printf("Description: %s\n", cmd.Description)
			fmt.Printf("Usage: %s\n", cmd.Usage)
			return nil
		}
		return fmt.Errorf("unknown command: %s", args[0])
	}

	fmt.Println("\nAvailable Commands:")
	for _, cmd := range c.commands {
		fmt.Printf("  %-10s - %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

func (c *CLI) cmdExit(args []string) error {
	c.running = false
	fmt.Println("Goodbye!")
	return nil
}

func (c *CLI) cmdClear(args []string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (c *CLI) cmdTime(args []string) error {
	now := time.Now()
	fmt.Printf("Current time: %s\n", now.Format("15:04:05"))
	fmt.Printf("Current date: %s\n", now.Format("2006-01-02"))
	return nil
}

func (c *CLI) cmdHistory(args []string) error {
	count := len(c.history)
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid count: %s", args[0])
		}
		if n < count {
			count = n
		}
	}

	fmt.Println("\nCommand History:")
	for i := len(c.history) - count; i < len(c.history); i++ {
		fmt.Printf("%3d: %s\n", i+1, c.history[i])
	}
	return nil
}

func (c *CLI) cmdCd(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("directory path required")
	}
	err := os.Chdir(args[0])
	if err != nil {
		return err
	}
	c.workingDir, _ = os.Getwd()
	return nil
}

func (c *CLI) cmdPwd(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

func (c *CLI) cmdLs(args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.IsDir() {
			fmt.Printf("\033[1;34m%s/\033[0m\n", entry.Name())
		} else {
			fmt.Println(entry.Name())
		}
	}
	return nil
}

func (c *CLI) cmdPrompt(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("style required (default/full)")
	}
	if args[0] != "default" && args[0] != "full" {
		return fmt.Errorf("invalid style: use 'default' or 'full'")
	}
	c.promptStyle = args[0]
	return nil
}

func (c *CLI) getPrompt() string {
	if c.promptStyle == "full" {
		dir, _ := os.Getwd()
		return fmt.Sprintf("\n%s > ", filepath.Base(dir))
	}
	return "\ncli > "
}

func (c *CLI) cmdEcho(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func (c *CLI) addToHistory(input string) {
	if len(input) > 0 {
		c.history = append(c.history, input)
		if len(c.history) > c.maxHistory {
			c.history = c.history[1:]
		}
	}
}

func (c *CLI) Run() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enhanced CLI v2.0")
	fmt.Println("Type 'help' for available commands")

	for c.running {
		fmt.Print(c.getPrompt())
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		c.addToHistory(input)
		args := strings.Fields(input)
		cmdName := args[0]

		if cmd, ok := c.commands[cmdName]; ok {
			if err := cmd.Execute(args[1:]); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", cmdName)
		}
	}
}

func mainCLI() {
	cli := NewCLI()
	cli.Run()
}
