package minienv

// ConversionError
type ConversionError struct{}

func (e *ConversionError) Error() string {
	return "conversion error"
}

// ---
