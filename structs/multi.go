package structs

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type structInfo struct {
	types map[string]structDef //key is the struct name, value is the struct fields
}

type structDef struct {
	name   string
	suffix int
	fields []string
}

// CreateStructFromJSON creates a Go struct from JSON data
func createStructsFromJSON(data map[string]interface{}, structName string) string {
	structs := &structInfo{
		types: make(map[string]structDef),
	}
	buildStructInfo(data, structName, structs, "")

	goCode := ""
	order := make([]string, 0, len(structs.types))
	for k := range structs.types {
		if k != structName {
			order = append(order, k)
		}
	}
	sort.Strings(order)
	order = append([]string{structName}, order...) //make sure the main struct we are creating is first

	for i := range order {
		goCode += "type " + order[i] + " struct {\n"
		for _, line := range structs.types[order[i]].fields {
			goCode += line + "\n"
		}
		goCode += "}\n\n"
	}

	return goCode
}

// Recursive function to create Go structs from JSON data
func buildStructInfo(data map[string]interface{}, structName string, structs *structInfo, parent string) string {
	typeDef := make([]string, 0)

	for key, value := range data {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			name := buildStructInfo(value.(map[string]interface{}), strings.Title(key), structs, structName)
			line := fmt.Sprintf("%s %s `json:\"%s\"`", strings.Title(key), name, key)
			typeDef = append(typeDef, line)
		case reflect.Slice:
			t := ""
			if len(value.([]interface{})) == 0 {
				t = "interface{}"
			} else {
				tRaw := reflect.TypeOf(value.([]interface{})[0])
				switch tRaw.Kind() {
				case reflect.Map:
					name := buildStructInfo(value.([]interface{})[0].(map[string]interface{}), strings.Title(key), structs, structName)
					t = "[]" + name
				default:
					t = fmt.Sprintf("[]%s", reflect.TypeOf(value.([]interface{})[0]).String())
				}
			}
			line := fmt.Sprintf("%s %s `json:\"%s\"`", strings.Title(key), t, key)
			typeDef = append(typeDef, line)
		default:
			line := fmt.Sprintf("%s %s `json:\"%s\"`", strings.Title(key), reflect.TypeOf(value).String(), key)
			typeDef = append(typeDef, line)
		}
	}

	// Sort the type definition to ensure consistent ordering
	sort.Slice(typeDef, func(i, j int) bool {
		return typeDef[i] < typeDef[j]
	})

	// Check if the struct already exists
	if _, ok := structs.types[strings.Title(structName)]; !ok {
		// TODO: Check if another struct already exists with a different name but the same fields. If so, try to find the commonality in the name
		structs.types[strings.Title(structName)] = structDef{
			name:   strings.Title(structName),
			suffix: 0,
			fields: typeDef,
		}
	} else if !reflect.DeepEqual(typeDef, structs.types[strings.Title(structName)].fields) {
		// If the struct already exists, but the fields are different, create a new struct
		i := 1
		for {
			newName := fmt.Sprintf("%s%d", strings.Title(structName), i)
			// Check if the new struct already exists with this suffix
			if _, ok := structs.types[newName]; !ok {
				structs.types[newName] = structDef{
					name:   newName,
					suffix: i,
					fields: typeDef,
				}
				return newName
			} else if !reflect.DeepEqual(typeDef, structs.types[newName].fields) {
				// If the struct already exists, but the fields are different, try the next suffix
				i++
			} else {
				// If the struct already exists, and the fields are the same, use the existing struct
				return newName
			}
		}
	}

	return structName
}
