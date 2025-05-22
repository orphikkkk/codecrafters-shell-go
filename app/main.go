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

type Command struct {
	Name         string
	Args         []string
	Redirections []Redirection
}

type Redirection struct {
	Type           string // ">", ">>", "2>", etc.
	FileDescriptor int    // 0 for stdin, 1 for stdout, 2 for stderr
	Target         string // Filename or target
}

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

		args := parseInput(input)
		if len(args) == 0 {
			continue
		}
		command := parseCommand(args)

		commandType := getCommandType(command.Name)

		switch commandType {
		case TypeBuiltin:
			if handler, exists := builtins[command.Name]; exists {
				handler(args)
			}
		case TypeExecutable:
			cmd := exec.Command(command.Name, args[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error:", err)
			}
		default:
			fmt.Println(command.Name + ": command not found")
		}
	}
}

func parseInput(input string) []string {
	var args []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	input = strings.TrimSpace(input)
	for _, r := range input {
		if escaped {
			if inDoubleQuote {
				// In double quotes, backslash only escapes $, ", \, and newline
				if r == '\\' || r == '$' || r == '"' || r == '\n' {
					current.WriteRune(r)
				} else {
					// For other characters, add both the backslash and the character
					current.WriteRune('\\')
					current.WriteRune(r)
				}
			} else {
				// Outside quotes, all escaped characters are taken literally
				current.WriteRune(r)
			}
			escaped = false
			continue
		}

		if r == '\'' && !inDoubleQuote {
			// Toggle single quote mode
			inSingleQuote = !inSingleQuote
			continue
		}

		if r == '"' && !inSingleQuote {
			// Toggle double quote mode
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if r == '\\' {
			if inSingleQuote {
				// Backslash in single quotes is always literal
				current.WriteRune(r)
			} else {
				// Mark as escaped for next iteration
				escaped = true
			}
			continue
		}

		if r == ' ' && !inSingleQuote && !inDoubleQuote {
			// Space outside of quotes means argument boundary
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func parseCommand(tokens []string) Command {
	cmd := Command{}

	if len(tokens) == 0 {
		return cmd
	}

	cmd.Name = tokens[0]
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]

		if !strings.Contains(token, ">") {
			cmd.Args = append(cmd.Args, token)
			continue
		}

		// Check for redirection patterns
		redirection := parseRedirection(token)

		// If redirection has a separate target, consume the next token
		if redirection.Target == "" && i+1 < len(tokens) {
			redirection.Target = tokens[i+1]
			i++
		}

		cmd.Redirections = append(cmd.Redirections, redirection)
	}

	return cmd
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
