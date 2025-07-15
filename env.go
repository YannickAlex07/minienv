// Package minienv provides a way to load environment variables into a struct.
// It supports options for fallback values, prefixes, and reading from env-files.
package minienv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ERRORS

// ErrInvalidInput is returned when the input to Load is not a pointer to a struct.
var ErrInvalidInput = fmt.Errorf("input struct is not a struct or a pointer to one")

// FieldError is returned when a particular field cannot be loaded.
// It contains the field name and the underlying error that caused the failure.
type FieldError struct {
	Field string
	Err   error
}

func (e FieldError) Error() string {
	return fmt.Sprintf("failed to load field \"%s\": %s", e.Field, e.Err.Error())
}

func (e FieldError) Unwrap() error {
	return e.Err
}

// TAG

type tag struct {
	LookupKey string
	Optional  bool
	Default   string
}

func parseTag(tagStr string) (tag, error) {
	if tagStr == "" {
		return tag{}, errors.New("tag string cannot be empty")
	}

	tagParts := strings.Split(tagStr, ",")

	t := tag{}
	for i, part := range tagParts {
		part = strings.TrimSpace(part)

		// first one needs to be the lookup key
		if i == 0 {
			t.LookupKey = part
			continue
		}

		optParts := strings.SplitN(part, "=", 2)
		switch optParts[0] {
		case "optional":
			t.Optional = true

		case "default":
			if len(optParts) < 2 {
				return tag{}, fmt.Errorf("default env value cannot be empty")
			}

			t.Default = strings.TrimSpace(optParts[1])

		default:
			return tag{}, fmt.Errorf("unknown tag option \"%s\"", optParts[0])
		}
	}

	return t, nil
}

// CONFIG

// LoadConfig holds the configuration for loading environment variables.
// Can be configured using the provided options or by writing your own option.
type LoadConfig struct {
	Prefix string
	Values map[string]string
}

func fetchFieldValue(config *LoadConfig, tag tag) (string, error) {
	// read the value from the environment and from any our overrides
	lookup := tag.LookupKey
	if config.Prefix != "" && !strings.HasPrefix(lookup, config.Prefix) {
		lookup = fmt.Sprintf("%s%s", config.Prefix, lookup)
	}

	envVal, envExists := os.LookupEnv(lookup)
	fallbackVal, fallbackExists := config.Values[lookup]

	// guard against the cases where we don't have any valeu that we can set
	if !envExists && !fallbackExists && !tag.Optional && tag.Default == "" {
		return "", fmt.Errorf("no value was found for field with lookup key: %s", lookup)
	}

	// Priority:
	// 1. Environment
	// 2. Fallback
	// 3. Default
	var val string
	if envExists {
		val = envVal
	} else if fallbackExists {
		val = fallbackVal
	} else {
		val = tag.Default
	}

	return val, nil
}

// Sets a field based on the kind and the provided value
func set(f reflect.Value, val string) error {
	k := f.Kind()
	switch k {
	// string
	case reflect.String:
		f.SetString(val)

	// int
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		f.SetInt(int64(i))

	// bool
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}

		f.SetBool(b)

	// float
	case reflect.Float32, reflect.Float64:
		fl, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}

		f.SetFloat(fl)

	// slice
	case reflect.Slice:
		vals := strings.Split(val, "|")

		slice := reflect.MakeSlice(f.Type(), len(vals), len(vals))
		for i, v := range vals {
			if err := set(slice.Index(i), v); err != nil {
				return fmt.Errorf("failed to set slice element %d: %w", i, err)
			}
		}

		f.Set(slice)

	// anything else is currently not supported
	default:
		return fmt.Errorf("unsupported type: %v", k.String())
	}

	return nil
}

// handleField handles parsing the tag for a field, fetching a value for it and setting it.
func handleField(config *LoadConfig, field reflect.Value, tagStr string) error {
	tag, err := parseTag(tagStr)
	if err != nil {
		return fmt.Errorf("failed to parse env tag \"%s\": %w", tagStr, err)
	}

	val, err := fetchFieldValue(config, tag)
	if err != nil {
		return fmt.Errorf("failed to fetch value: %w", err)
	}

	err = set(field, val)
	if err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

// handleStruct recursively handles a struct, parsing its fields, checking if the have
// the `env` struct tag set and then passing them to the handleField function.
func handleStruct(s reflect.Value, config *LoadConfig) error {
	for i := range s.NumField() {
		// handle recursive cases
		field := s.Field(i)
		if field.Kind() == reflect.Struct {
			err := handleStruct(field, config)
			if err != nil {
				return err
			}

			continue
		}

		// parse the field information
		structField := s.Type().Field(i)
		value, found := structField.Tag.Lookup("env")
		if !found {
			continue
		}

		// check if we can actually set the field
		if !field.IsValid() || !field.CanSet() {
			return FieldError{
				Field: structField.Name,
				Err:   errors.New("field is not valid or cannot be set"),
			}
		}

		err := handleField(config, field, value)
		if err != nil {
			return FieldError{
				Field: structField.Name,
				Err:   err,
			}
		}
	}

	return nil
}

// Load loads environment variables into a struct based on the `env` struct tag.
// It can be configured using the provided options or by writing your own option.
// The struct must be a pointer to a struct, otherwise an error will be returned.
//
// The function will recursively fill the struct with values from the environment variables,
// using the `env` struct tag to determine which environment variable to use for each field.
func Load(obj any, options ...Option) error {
	// read in any overrides the user wants to do
	config := LoadConfig{
		Values: make(map[string]string),
	}

	for _, option := range options {
		err := option(&config)
		if err != nil {
			return err
		}
	}

	// we can only set things if we receive a pointer that points to a struct
	p := reflect.ValueOf(obj)
	if p.Kind() != reflect.Ptr {
		return ErrInvalidInput
	}

	s := reflect.Indirect(p)
	if !s.IsValid() || s.Kind() != reflect.Struct {
		return ErrInvalidInput
	}

	// this will recursively fill the struct
	err := handleStruct(s, &config)
	if err != nil {
		return err
	}

	return nil
}
