package lang

import "fmt"

type InvalidCommandError struct {
	CommandName string
}

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf("Invalid command name: %s", e.CommandName)
}
