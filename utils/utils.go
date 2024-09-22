package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func ContainsString(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func GetStructVals(s interface{}) string {
	var result strings.Builder
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	// Make sure we're dealing with a struct
	if t.Kind() != reflect.Struct {
		return ""
	}

	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i)
		if value.String() != " " {
			result.WriteString(fmt.Sprintf("%v, ", value.Interface()))
		}
	}
	return result.String()
}

// func ParseJSONFile(filepath string, v interface{}) error {
// 	// Read the entire file
// 	fileData, err := os.ReadFile(filepath)
// 	if err != nil {
// 		return fmt.Errorf("error reading file: %w", err)
// 	}

// 	// Unmarshal the JSON data into the provided interface
// 	err = json.Unmarshal(fileData, v)
// 	if err != nil {
// 		return fmt.Errorf("error unmarshaling JSON: %w", err)
// 	}

// 	return nil
// }
