package structs

import (
	"encoding/json"
	"go/format"
	"strings"
)

// CreateStructFromJSON creates a Go struct from a JSON string
func CreateStructFromJSON(jsonData string, structName string, inline bool) (string, error) {
	// Parse JSON into a map[string]interface{}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		innerData := make([]interface{}, 0)
		if err = json.Unmarshal([]byte(jsonData), &innerData); err != nil {
			return "", err
		}
		data = make(map[string]interface{})
		data[getNameWithSuffix(structName, "Items")] = innerData
	}

	// Generate Go structs
	goCode := ""
	if inline {
		goCode = createInlineStructFromJSON(data, structName)
	} else {
		goCode = createStructsFromJSON(data, structName)
	}

	// Format the code
	formattedCode, err := format.Source([]byte(goCode))
	if err != nil {
		return goCode, err
	}

	return string(formattedCode), nil
}

func title(name string) string {
	return strings.Replace(strings.Title(name), "-", "_", -1)
}
