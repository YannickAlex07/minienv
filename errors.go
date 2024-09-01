package minienv

import "fmt"

// CONST ERRORS

var ErrInvalidInput = fmt.Errorf("input struct is not a struct or a pointer to one")

// GENERAL ERROR

type FieldError struct {
	Field string
	Err   error
}

func (e FieldError) Error() string {
	return fmt.Sprintf("error for field \"%s\": %s", e.Field, e.Err.Error())
}

// PARSING ERROR

type TagParsingError struct {
	Field string
	Err   error
}

func (e TagParsingError) Error() string {
	return fmt.Sprintf("failed to parse tag for field %s: %s", e.Field, e.Err.Error())
}

// CONVERSION ERROR

type CoversionError struct {
	Field string
	Value string
	Err   error
}

func (e CoversionError) Error() string {
	return fmt.Sprintf("failed to convert value %s for field \"%s\" to target type: %s", e.Value, e.Field, e.Err.Error())
}
