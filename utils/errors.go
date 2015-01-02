package utils

import (
	"fmt"
)

// The serializable Error structure.
type Error struct {
	Message string   `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

// NewError creates an error instance with the specified code and message.
func NewError(msg string) *Error {
	return &Error{
		Message: msg,
	}
}
