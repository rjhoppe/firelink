package utils

import (
	"reflect"
)

func ContainsString(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func GetFieldValue(dataStruct interface{}, field string) string {
	val := reflect.ValueOf(dataStruct)
	f := val.FieldByName(field)
	if f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	return ""
}

func ListAllEndpoints() []string {
	return []string{
		"/help",
		"/ebook/find/:title",
		"/ebook/download/:title",
		"/bartender/random",
		"/bartender/cache/backup",
		"/bartender/history",
		"/bartender/save",
		"/dinner/random",
		"/dinner/cache/backup",
		"/dinner/recipe/:id",
		"/database/backup",
	}
}
