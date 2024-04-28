package minienv

import (
	"strconv"
)

func parseString(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil

	case int, int8, int16, int32, int64:
		return strconv.Itoa(v.(int)), nil

	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64), nil

	case bool:
		return strconv.FormatBool(v), nil

	default:
		return "", &ConversionError{}
	}
}

func parseInt(val interface{}) (int64, error) {
	switch v := val.(type) {
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, &ConversionError{}
		}

		return i, nil

	case int, int8, int16, int32, int64:
		return v.(int64), nil

	case float32, float64:
		return int64(v.(float64)), nil

	default:
		return 0, &ConversionError{}
	}
}

func parseBool(val interface{}) (bool, error) {
	switch v := val.(type) {
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, &ConversionError{}
		}

		return b, nil

	case bool:
		return v, nil

	default:
		return false, &ConversionError{}
	}
}

func parseFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, &ConversionError{}
		}

		return f, nil

	case int, int8, int16, int32, int64:
		return float64(v.(int64)), nil

	case float32, float64:
		return v.(float64), nil

	default:
		return 0, &ConversionError{}
	}
}
