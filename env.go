package minienv

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Option func(map[string]string) error

// Load variables from the environment into the provided struct.
// It will try to match environment variables to field that contain an `env` tag.
//
// The obj must be a pointer to a struct.
// Additional options can be supplied for overriding environment variables.
func Load(obj interface{}, options ...Option) error {
	// read in any overrides the user wants to do
	overrides := make(map[string]string)

	for _, option := range options {
		err := option(overrides)
		if err != nil {
			return err
		}
	}

	// we can only set things if we receive a pointer that points to a struct
	p := reflect.ValueOf(obj)
	if p.Kind() != reflect.Ptr {
		return errors.New("obj must be a pointer")
	}

	s := reflect.Indirect(p)
	if !s.IsValid() || s.Kind() != reflect.Struct {
		return errors.New("obj must be a struct")
	}

	// this will recursively fill the struct
	err := handleStruct(s, overrides)
	if err != nil {
		return err
	}

	return nil
}

// Handles a struct recursively by iterating over its fields
// and then setting the field with the appropiate variable if one was found.
func handleStruct(s reflect.Value, overrides map[string]string) error {
	for i := 0; i < s.NumField(); i++ {
		// handle recursive cases
		field := s.Field(i)
		if field.Kind() == reflect.Struct {
			handleStruct(field, overrides)
			continue
		}

		// Check if the tag is present
		tagVal, found, required := parseTag(s.Type().Field(i))
		if !found {
			continue
		}

		// check if we can actually set the field
		if !field.IsValid() || !field.CanSet() {
			return errors.New("field is invalid or cannot be set")
		}

		// read the value from the environment
		envVal, exists := os.LookupEnv(tagVal)
		overrideVal, ok := overrides[tagVal]
		if !exists && !ok && required {
			return errors.New("environment variable not found")
		}

		// priority
		// 1. Overrides
		// 2. Environment
		var val string
		if ok {
			val = overrideVal
		} else {
			val = envVal
		}

		// update the affected field
		err := setField(field, val)
		if err != nil {
			return err
		}
	}

	return nil
}

// Sets a field based on the kind and the provided value
// This here tries to convert the value to the appropiate type
func setField(f reflect.Value, val string) error {
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

	// anything else is not supported
	default:
		return errors.New("unsupported type")
	}

	return nil
}

// ...
func parseTag(field reflect.StructField) (string, bool, bool) {
	required := true

	value, found := field.Tag.Lookup("env")
	if !found {
		return "", false, required
	}

	// check any tag options
	parts := strings.Split(value, ",")
	for _, p := range parts[1:] {
		trimmed := strings.TrimSpace(p)

		// tag is optional
		if trimmed == "optional" {
			required = false
		}
	}

	return parts[0], true, required
}
