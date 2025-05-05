package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	TypeBuiltin    = "builtin"
	TypeExecutable = "executable"
	TypeUnknown    = "unknown"
)

var builtins map[string]func([]string)

func init() {
	builtins = map[string]func([]string){
		"echo": cmdEcho,
		"exit": cmdExit,
		"type": cmdType,
		"pwd":  cmdPwd,
		"cd":   cmdCd,
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		command := args[0]
		commandType := getCommandType(command)

		switch commandType {
		case TypeBuiltin:
			if handler, exists := builtins[command]; exists {
				handler(args)
			}
		case TypeExecutable:
			cmd := exec.Command(command, args[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error:", err)
			}
		default:
			fmt.Println(command + ": command not found")
		}
	}
}

func getCommandType(cmd string) string {
	// Check if it's a built-in command
	if _, exists := builtins[cmd]; exists {
		return TypeBuiltin
	}

	// Check if it's an executable in PATH
	if path, err := exec.LookPath(cmd); err == nil && path != "" {
		return TypeExecutable
	}

	return TypeUnknown
}
