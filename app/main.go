package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	TypeBuiltin    = "builtin"
	TypeExecutable = "executable"
	TypeUnknown    = "unknown"
)

var builtins = []string{
	"echo",
	"exit",
	"type",
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
			switch command {
			case "exit":
				cmdExit(args)
				continue
			case "echo":
				cmdEcho(args)
				continue
			case "type":
				cmdType(args)
				continue
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

func cmdExit(args []string) {
	status := 0
	var err error

	if len(args) > 1 {
		status, err = validateStatusCode(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	os.Exit(status)
}

func cmdEcho(args []string) {
	message := ""

	if len(args) > 1 {
		message = strings.Join(args[1:], " ")
	}
	fmt.Println(message)
}

func cmdType(args []string) {
	commandToCheck := ""
	if len(args) > 1 {
		commandToCheck = args[1]
	}
	if commandToCheck == "" {
		return
	}

	commandType := getCommandType(commandToCheck)

	switch commandType {
	case TypeBuiltin:
		fmt.Println(commandToCheck + " is a shell builtin")
	case TypeExecutable:
		path, _ := exec.LookPath(commandToCheck)
		fmt.Printf("%s is %s\n", commandToCheck, path)
	default:
		fmt.Println(commandToCheck + ": not found")
	}
}

func getCommandType(cmd string) string {
	// Check if it's a built-in command
	if slices.Contains(builtins, cmd) {
		return TypeBuiltin
	}

	// Check if it's an executable in PATH
	if path, err := exec.LookPath(cmd); err == nil && path != "" {
		return TypeExecutable
	}

	return TypeUnknown
}

func validateStatusCode(statusCode string) (int, error) {
	if code, err := strconv.Atoi(statusCode); err == nil {
		return code, nil
	} else {
		return 0, errors.New("Invalid exit status, must be an integer")
	}
}
