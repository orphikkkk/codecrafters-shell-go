package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

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

		switch command {
		case "exit":
			exitCommand(args)
			continue
		case "echo":
			echoCommand(args)
			continue
		case "type":
			typeCommand(args)
			continue
		}

		fmt.Println(args[0] + ": command not found")
	}
}

func exitCommand(args []string) {
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

func echoCommand(args []string) {
	message := ""

	if len(args) > 1 {
		message = strings.Join(args[1:], " ")
	}
	fmt.Println(message)
}

func typeCommand(args []string) {
	commandToCheck := ""

	if len(args) > 1 {
		commandToCheck = args[1]
	}
	if commandToCheck == "" {
		return
	}
	switch commandToCheck {
	case "echo", "exit", "type":
		fmt.Println(commandToCheck + " is a shell builtin")
	default:
		fmt.Println(commandToCheck + ": not found")
	}
}

func validateStatusCode(statusCode string) (int, error) {
	if code, err := strconv.Atoi(statusCode); err == nil {
		return code, nil
	} else {
		return 0, errors.New("Invalid exit status, must be an integer")
	}
}
