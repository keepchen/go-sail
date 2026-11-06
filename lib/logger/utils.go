package logger

import (
	"encoding/json"
)

// MarshalInterfaceValue serializes an interface to a string
//
// Mainly used for logging
func MarshalInterfaceValue(obj any) string {
	js, _ := json.Marshal(obj)

	return string(js)
}
