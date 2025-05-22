package redirect

import (
	"errors"
	"os"
)

func Handle(redirections []Redirection) error {
	// Store original file descriptors to restore later
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	originalStdin := os.Stdin

	// Defer restoration of original descriptors
	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
		os.Stdin = originalStdin
	}()

	redir := redirections[0]
	switch redir.Type {
	case ">":
		return redirectOutput(redir, false)
	// case ">>":
	// 	handleAppendRedirection(redirections)
	// case "<":
	// 	handleInputRedirection(redirections)
	// case "<<":
	// 	handleHereDocumentRedirection(redirections)
	// case "2>&1":
	// 	handleFileDescriptorDuplication(redirections)
	// case "1>&2":
	// 	handleFileDescriptorDuplication(redirections)
	default:
		return errors.New("Unknown redirection type")
	}
}

func redirectOutput(redirections Redirection, append bool) error {
	var flags int

	if append {
		flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	} else {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	file, err := os.OpenFile(redirections.Target, flags, 0644)
	if err != nil {
		return err
	}

	os.Stdout = file
	return nil
}
