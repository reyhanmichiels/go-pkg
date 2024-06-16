package sql

import (
	"fmt"
	"strings"
)

const (
	queryLogMessage = "executing query: %s, with query string: %s"
)

// Replace query bindvars with args value
func replaceBindVarsWithArgs(str string, args ...interface{}) string {
	str = strings.Join(strings.Fields(str), " ")
	for _, a := range args {
		str = strings.Replace(str, "?", fmt.Sprintf("%v", a), 1)
	}
	return str
}

// TODO: create func for replace Named Execution
