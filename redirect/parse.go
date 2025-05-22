package redirect

import (
	"strings"
)

/**
 * This file contains the implementation to parse Redirections.
 * this will mainly check for:
 * Output redirection: >, 1>, 2>, &>
 * Append redirection: >>, 1>>, 2>>
 * Input redirection: <
 * Here documents: <<
 * File descriptor duplication: 2>&1, 1>&2
 */

type Redirection struct {
	Type           string
	FileDescriptor int
	Target         string
}

func ParseRedirection(token string) Redirection {
	red := Redirection{
		Type:           ">",
		FileDescriptor: 1,
	}

	if strings.Contains(token, ">") {
		red.Type = ">"
		red.FileDescriptor = 1
	}

	return red
}
