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
		return handleOutputRedirection(redir)
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

func handleOutputRedirection(redirections Redirection) error {

	file, err := os.Create(redirections.Target)
	if err != nil {
		return err
	}
	defer file.Close()

	os.Stdout = file
	return nil
}
