package minienv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Option func(map[string]string) error

// This struct hold all the metadata about a found "env"-tag for a field
type tag struct {
	// This is the name of the env variable we need to look for
	name string

	// This is a flag that tells us if the variable is required
	required bool

	// This is the default value for the variable, can be empty and therefore invalid
	defaultValue string
}

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
		return ErrInvalidInput
	}

	s := reflect.Indirect(p)
	if !s.IsValid() || s.Kind() != reflect.Struct {
		return ErrInvalidInput
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

		// Check if the tag is present skip if not
		tag, found, err := parseTag(s.Type().Field(i))
		if !found {
			continue
		}

		// something went wrong parsing the tag
		if err != nil {
			return LoadError{
				Field: s.Type().Field(i).Name,
				Err:   err,
			}
		}

		// check if we can actually set the field
		if !field.IsValid() || !field.CanSet() {
			return LoadError{
				Field: s.Type().Field(i).Name,
				Err:   errors.New("field is not valid or cannot be set"),
			}
		}

		// read the value from the environment and from any our overrides
		envVal, envExists := os.LookupEnv(tag.name)
		overrideVal, overrideExists := overrides[tag.name]

		// guard against the cases where we don't have any valeu that we can set
		if !envExists && !overrideExists && tag.required && tag.defaultValue == "" {
			return LoadError{
				Field: s.Type().Field(i).Name,
				Err:   errors.New("required field has no value and no default"),
			}
		}

		// priority
		// 1. Overrides
		// 2. Environment
		// 3. Default
		var val string
		if overrideExists {
			val = overrideVal
		} else if envExists {
			val = envVal
		} else {
			val = tag.defaultValue
		}

		// update the affected field
		err = setField(field, val)
		if err != nil {
			// we wrap the error for some metadata
			return LoadError{
				Field: s.Type().Field(i).Name,
				Err:   err,
			}
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
		return fmt.Errorf("unsupported type: %v", k.String())
	}

	return nil
}

// Parses the `env` tag and returns the bundled information about the tag.
// The first return value is the tag itself, the second return value is a flag indicating if the tag was found
// and the third return value is an error if the tag was invalid.
func parseTag(field reflect.StructField) (tag, bool, error) {
	required := true
	var defaultVal string

	value, found := field.Tag.Lookup("env")
	if !found {
		return tag{}, false, nil
	}

	// check any tag options
	parts := strings.Split(value, ",")
	for _, p := range parts[1:] {
		trimmed := strings.TrimSpace(p)
		splitted := strings.Split(trimmed, "=")

		// tag is optional
		if splitted[0] == "optional" {
			required = false

		} else if splitted[0] == "default" {

			// if we have more or less than 2 elements we have an invalid tag
			if len(splitted) != 2 {
				return tag{}, true, errors.New("invalid default tag")
			}

			defaultVal = splitted[1]
		}
	}

	return tag{
		name:         parts[0],
		required:     required,
		defaultValue: defaultVal,
	}, true, nil
}
