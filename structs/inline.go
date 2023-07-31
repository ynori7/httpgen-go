package structs

import (
	"fmt"
	"reflect"
	"sort"
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
		originalKey := key
		// Check if the key starts with a number
		if len(key) > 0 && key[0] >= '0' && key[0] <= '9' {
			key = structName + "Num_" + key // Prefix the key with a string
		}
		value := data[originalKey]
		if reflect.TypeOf(value) == nil {
			structCode += fmt.Sprintf("%s interface{} `json:\"%s\"`\n", title(key), originalKey)
			continue
		}
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			structCode += fmt.Sprintf("%s %s `json:\"%s\"`\n", title(key), createInlineStructFromJSON(value.(map[string]interface{}), ""), originalKey)
		case reflect.Slice:
			if len(value.([]interface{})) == 0 {
				structCode += fmt.Sprintf("%s []interface{} `json:\"%s\"`\n", title(key), originalKey)
			} else {
				t := reflect.TypeOf(value.([]interface{})[0])
				switch t.Kind() {
				case reflect.Map:
					structCode += fmt.Sprintf("%s []%s `json:\"%s\"`\n", title(key), createInlineStructFromJSON(value.([]interface{})[0].(map[string]interface{}), ""), originalKey)
				default:
					structCode += fmt.Sprintf("%s []%s `json:\"%s\"`\n", title(key), t.String(), originalKey)
				}
			}
		default:
			structCode += fmt.Sprintf("%s %s `json:\"%s\"`\n", title(key), reflect.TypeOf(value).String(), originalKey)
		}
	}
	if structName != "" {
		structCode += "}\n"
	} else {
		structCode += "}"
	}
	return structCode
}
