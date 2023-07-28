package structs

import (
	"encoding/json"
	"go/format"
)

// CreateStructFromJSON creates a Go struct from a JSON string
func CreateStructFromJSON(jsonData string, structName string, inline bool) (string, error) {
	// Parse JSON into a map[string]interface{}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return "", err
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
