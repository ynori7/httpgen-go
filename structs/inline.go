package structs

import (
	"fmt"
	"reflect"
	"strings"
)

// Recursive function to create Go structs from JSON data
func createInlineStructFromJSON(data map[string]interface{}, structName string) string {
	structCode := ""
	if structName != "" {
		structCode = fmt.Sprintf("type %s struct {\n", structName)
	} else {
		structCode = " struct {\n"
	}
	for key, value := range data {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			structCode += fmt.Sprintf("%s %s `json:\"%s\"`\n", strings.Title(key), createInlineStructFromJSON(value.(map[string]interface{}), ""), key)
		case reflect.Slice:
			t := reflect.TypeOf(value.([]interface{})[0]).String()
			structCode += fmt.Sprintf("%s []%s `json:\"%s\"`\n", strings.Title(key), t, key)
		default:
			structCode += fmt.Sprintf("%s %s `json:\"%s\"`\n", strings.Title(key), reflect.TypeOf(value).String(), key)
		}
	}
	if structName != "" {
		structCode += "}\n"
	} else {
		structCode += "}"
	}
	return structCode
}
