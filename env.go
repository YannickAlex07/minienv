package minienv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/yannickalex07/minienv/internal/tag"
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

		// parse the field information
		structField := s.Type().Field(i)

		// parse the tag
		value, found := structField.Tag.Lookup("env")
		if !found {
			continue
		}

		tag, err := tag.ParseMinienvTag(value)
		if err != nil {
			return LoadError{
				Field: structField.Name,
				Err:   err,
			}
		}

		// check if we can actually set the field
		if !field.IsValid() || !field.CanSet() {
			return LoadError{
				Field: structField.Name,
				Err:   errors.New("field is not valid or cannot be set"),
			}
		}

		// get the value that we need to set
		val, err := getValue(config, tag)
		if err != nil {
			return LoadError{
				Field: structField.Name,
				Err:   err,
			}
		}

		// update the affected field
		err = setField(field, val, tag)
		if err != nil {
			// we wrap the error for some metadata
			return LoadError{
				Field: structField.Name,
				Err:   err,
			}
		}
	}

	return nil
}

func getValue(config *LoadConfig, tag tag.MinienvTag) (string, error) {
	// read the value from the environment and from any our overrides
	lookup := tag.LookupName
	if config.Prefix != "" && !strings.HasPrefix(lookup, config.Prefix) {
		lookup = fmt.Sprintf("%s%s", config.Prefix, lookup)
	}

	envVal, envExists := os.LookupEnv(lookup)
	fallbackVal, fallbackExists := config.Values[lookup]

	// guard against the cases where we don't have any valeu that we can set
	if !envExists && !fallbackExists && !tag.Optional && tag.Default == "" {
		return "", fmt.Errorf("no value was found for required field with lookup key %s", lookup)
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
// This here tries to convert the value to the appropiate type
func setField(f reflect.Value, val string, tag tag.MinienvTag) error {
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

	case reflect.Slice:
		// split the string by the splitOn separator
		vals := strings.Split(val, tag.SplitOn)

		// create the slice
		elementKind := f.Type().Elem().Kind()
		slice := reflect.MakeSlice(f.Type(), len(vals), len(vals))
		for i, v := range vals {
			converted, err := convertPrimitiveValue(v, elementKind)
			if err != nil {
				return err
			}

			slice.Index(i).Set(reflect.ValueOf(converted))
		}

		// set the field
		f.Set(slice)

	// anything else is not supported
	default:
		return fmt.Errorf("unsupported type: %v", k.String())
	}

	return nil
}

func convertPrimitiveValue(val string, kind reflect.Kind) (any, error) {
	switch kind {
	case reflect.String:
		return val, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Atoi(val)

	case reflect.Bool:
		return strconv.ParseBool(val)

	case reflect.Float32, reflect.Float64:
		return strconv.ParseFloat(val, 64)

	default:
		return nil, fmt.Errorf("unsupported type: %v", kind.String())
	}
}
