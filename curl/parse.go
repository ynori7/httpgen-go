package curl

import (
	"errors"
	"strings"

	"github.com/mattn/go-shellwords"
)

func Parse(curlCommand string) (*Command, error) {
	url, params := ExtractParameters(curlCommand)
	if url == "" {
		return nil, errors.New("no URL found in the curl command")
	}

	body, hasBody := GetBody(params)
	c := &Command{
		URL:     url,
		Headers: make(map[string]string),
		Method:  GetMethod(params),
		Body:    body,
		HasBody: hasBody,
	}

	// Extract the headers and request body (if present)
	for _, v := range params["-H"] {
		headerParts := strings.SplitN(v, ":", 2)
		headerKey := strings.TrimSpace(headerParts[0])
		headerValue := strings.TrimSpace(headerParts[1])
		if headerKey != "" {
			c.Headers[headerKey] = headerValue
		}
	}

	return c, nil
}

// ExtractParameters parses the curl command and returns the URL and parameters
func ExtractParameters(curlCommand string) (string, map[string][]string) {
	parameters := make(map[string][]string)

	url := ""

	// Use shellwords to correctly split the bash command
	args, err := shellwords.Parse(curlCommand)
	if err != nil {
		panic(err)
	}

	for i := 1; i < len(args); i++ {
		arg := args[i]

		// Check if the argument starts with "-" or "--" (indicating a parameter)
		if strings.HasPrefix(arg, "-") {
			// Check if the parameter has a value (format: --param=value)
			if strings.Contains(arg, "=") {
				paramParts := strings.SplitN(arg, "=", 2)
				param := paramParts[0]
				value := paramParts[1]
				parameters[param] = append(parameters[param], value)
			} else {
				// Check if the parameter is a hardcoded boolean
				if hardcodedBooleans[arg] {
					parameters[arg] = []string{"true"}
				} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					// Check if the parameter has a separate value
					value := args[i+1]
					parameters[arg] = append(parameters[arg], value)
					i++
				} else {
					parameters[arg] = []string{"true"}
				}
			}
		} else {
			if strings.HasPrefix(arg, "http") {
				url = arg
			}
		}
	}

	return url, parameters
}
