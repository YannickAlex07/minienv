package minienv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Option func(*LoadConfig) error

type LoadConfig struct {
	Prefix string
	Values map[string]string
}

// Load variables from the environment into the provided struct.
// It will try to match environment variables to field that contain an `env` tag.
//
// The obj must be a pointer to a struct.
// Additional options can be supplied for overriding environment variables.
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

// Handles a struct recursively by iterating over its fields
// and then setting the field with the appropiate variable if one was found.
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
		lookup := tag.name
		if config.Prefix != "" && !strings.HasPrefix(lookup, config.Prefix) {
			lookup = fmt.Sprintf("%s%s", config.Prefix, lookup)
		}

		envVal, envExists := os.LookupEnv(lookup)
		fallbackVal, fallbackExists := config.Values[lookup]

		// guard against the cases where we don't have any valeu that we can set
		if !envExists && !fallbackExists && tag.required && tag.defaultValue == "" {
			return LoadError{
				Field: s.Type().Field(i).Name,
				Err:   errors.New("required field has no value and no default"),
			}
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
