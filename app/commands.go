package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

func cmdPwd(args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Issue while pringting current directory", err)
		os.Exit(1)
	}

	fmt.Println(dir)
}

func cmdCd(args []string) {
	path := args[1]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", path)
		return
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Println(err)
	}
}

func validateStatusCode(statusCode string) (int, error) {
	if code, err := strconv.Atoi(statusCode); err == nil {
		return code, nil
	} else {
		return 0, errors.New("Invalid exit status, must be an integer")
	}
}
