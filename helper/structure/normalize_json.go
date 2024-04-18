// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package structure

import "encoding/json"

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func NormalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// Recursively convert single-element arrays to single elements
	var simplifySingleElementArrays func(data interface{}) interface{}
	simplifySingleElementArrays = func(data interface{}) interface{} {
		switch x := data.(type) {
		case []interface{}:
			if len(x) == 1 {
				return simplifySingleElementArrays(x[0]) // Return the single element, further simplified
			}
			for i, v := range x {
				x[i] = simplifySingleElementArrays(v) // Apply the same simplification to each element
			}
		case map[string]interface{}:
			for k, v := range x {
				x[k] = simplifySingleElementArrays(v) // Apply simplification recursively to each value
			}
		}
		return data
	}

	// Apply the simplification
	simplifiedJson := simplifySingleElementArrays(j)

	bytes, err := json.Marshal(simplifiedJson)
	return string(bytes[:]), nil
}
