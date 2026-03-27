package framework

import "fmt"

// invalidFormatException represents an error for invalid format
type invalidFormatException struct {
	Message string
}

// Error implements the error interface
func (e *invalidFormatException) Error() string {
	return fmt.Sprintf("InvalidFormatException: %s", e.Message)
}

// InvalidFormatException creates a new InvalidFormatException
func InvalidFormatException(message string) error {
	return &invalidFormatException{Message: message}
}
