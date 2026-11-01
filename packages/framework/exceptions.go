package framework

import "fmt"

type exception struct {
	name    string
	message string
}

func (e *exception) Error() string {
	return fmt.Sprintf("%s: %s", e.name, e.message)
}

func InvalidFormatException(message string) error {
	return &exception{name: "InvalidFormatException", message: message}
}

func NotFoundException(message string) error {
	return &exception{name: "NotFoundException", message: message}
}

func NotImplementedException(message string) error {
	if message == "" {
		message = "This feature is not yet implemented"
	}
	return &exception{name: "NotImplementedException", message: message}
}
