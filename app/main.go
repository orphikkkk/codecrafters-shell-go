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
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		args := strings.Fields(command)
		if len(args) == 0 {
			continue
		}

		if args[0] == "echo" {
			message := ""

			if len(args) > 1 {
				message = strings.Join(args[1:], " ")
			}
			fmt.Println(message)
			continue
		}

		if args[0] == "exit" {
			status := 0

			if len(args) > 1 {
				status, err = validateStatusCode(args[1])
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
			os.Exit(status)
		}

		fmt.Println(args[0] + ": command not found")
	}
}

func validateStatusCode(statusCode string) (int, error) {
	if code, err := strconv.Atoi(statusCode); err == nil {
		return code, nil
	} else {
		return 0, errors.New("Invalid exit status, must be an integer")
	}
}
