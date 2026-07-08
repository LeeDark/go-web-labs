package data

import (
	"fmt"
	"strconv"
)

// Runtime represents the duration of an entity, typically measured in minutes, as an integer value.
type Runtime int32

// MarshalJSON customizes the JSON encoding of Runtime by formatting it as a quoted string with "mins" suffix.
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
