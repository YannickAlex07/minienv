package minienv

import "fmt"

// CONST ERRORS

var ErrInvalidInput = fmt.Errorf("input struct is not a struct or a pointer to one")

// Loading Error

type LoadError struct {
	Field string
	Err   error
}

func (e LoadError) Error() string {
	return fmt.Sprintf("failed to load field \"%s\": %s", e.Field, e.Err.Error())
}
