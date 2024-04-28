package minienv

import (
	"errors"
	"os"
	"reflect"
	"strings"
)

type Option func(map[string]interface{}) error

func Load(obj interface{}, options ...Option) error {
	// read in any overrides the user wants to do
	overrides := make(map[string]interface{})

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

func handleStruct(s reflect.Value, overrides map[string]interface{}) error {
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
		_, ok := overrides[tagVal]
		if !exists && !ok && required {
			return errors.New("environment variable not found")
		}

		// update the affected field
		err := setField(field, envVal)
		if err != nil {
			return err
		}
	}

	return nil
}

func setField(f reflect.Value, val interface{}) error {
	k := f.Kind()
	switch k {
	// string
	case reflect.String:
		parsed, err := parseString(val)
		if err != nil {
			return err
		}

		f.SetString(parsed)

	// int
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := parseInt(val)
		if err != nil {
			return err
		}

		f.SetInt(i)

	// bool
	case reflect.Bool:
		b, err := parseBool(val)
		if err != nil {
			return err
		}

		f.SetBool(b)

	// float
	case reflect.Float32, reflect.Float64:
		fl, err := parseFloat(val)
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
