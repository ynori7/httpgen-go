package structs

import (
	"fmt"
	"reflect"
	"sort"
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

	order := make([]string, 0, len(data))
	for key := range data {
		order = append(order, key)
	}

	sort.Strings(order)

	for _, key := range order {
		value := data[key]
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			structCode += fmt.Sprintf("%s %s `json:\"%s\"`\n", strings.Title(key), createInlineStructFromJSON(value.(map[string]interface{}), ""), key)
		case reflect.Slice:
			if len(value.([]interface{})) == 0 {
				structCode += fmt.Sprintf("%s []interface{} `json:\"%s\"`\n", strings.Title(key), key)
			} else {
				t := reflect.TypeOf(value.([]interface{})[0])
				switch t.Kind() {
				case reflect.Map:
					structCode += fmt.Sprintf("%s []%s `json:\"%s\"`\n", strings.Title(key), createInlineStructFromJSON(value.([]interface{})[0].(map[string]interface{}), ""), key)
				default:
					structCode += fmt.Sprintf("%s []%s `json:\"%s\"`\n", strings.Title(key), t.String(), key)
				}
			}
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
