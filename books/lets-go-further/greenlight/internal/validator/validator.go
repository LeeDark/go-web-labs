package validator

import (
	"regexp"
	"slices"
)

// EmailRX is a regular expression used to validate the format of email addresses.
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+" +
		"@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator represents a structure for collecting validation errors as key-value pairs.
type Validator struct {
	Errors map[string]string
}

// New creates and returns a new Validator instance with an empty Errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid checks if the Validator contains no errors and returns true if the Errors map is empty.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the Validator's Errors map for the given key if it does not already exist.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error to the Validator's Errors map if the condition `ok` is false, using the specified key and message.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue checks if the given value is present in the list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches checks if the given string matches the provided regular expression pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks if all elements in the provided slice are distinct and returns true if no duplicates are found.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
