package generator

import (
	"bufio"
	"bytes"
	"text/template"

	"github.com/ynori7/httpgen-go/curl"
	"github.com/ynori7/httpgen-go/structs"
)

// GoTemplate is a struct for the Go template
type GoTemplate struct {
	Command      *curl.Command
	RequestBody  string
	ResponseBody string
}

// NewGoTemplate creates a new GoTemplate
func NewGoTemplate(cmd *curl.Command, resp string) GoTemplate {
	return GoTemplate{
		Command:      cmd,
		ResponseBody: resp,
	}
}

// ExecuteGoTemplate executes the Go template
func (h GoTemplate) ExecuteGoTemplate(useInlineStructs bool) (string, error) {
	if h.Command.HasBody {
		h.RequestBody, _ = structs.CreateStructFromJSON(h.Command.Body, "Request", useInlineStructs)
	}
	if h.ResponseBody != "" {
		h.ResponseBody, _ = structs.CreateStructFromJSON(h.ResponseBody, "Response", useInlineStructs)
	}

	t := template.Must(template.New("txt").Parse(goTemplate))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := t.Execute(w, h)
	if err != nil {
		return "", err
	}

	w.Flush()
	return b.String(), nil
}

const goTemplate = `package main

import (
	"bytes"
	"fmt"
	"net/http"{{ if .Command.HasBody }}
	"strings"{{ end }}
)

func main() {
	// Create the HTTP request{{ if .Command.HasBody }}
	body := strings.NewReader({{ .Command.Body | printf "%q" }})
{{ end }}
	request, err := http.NewRequest({{ .Command.Method | printf "%q"}}, {{ .Command.URL | printf "%q" }}, {{ if not .Command.Body }}nil{{ else }}body{{ end }})
	if err != nil {
		panic(err)
	}{{ range $k, $v := .Command.Headers }}
	request.Header.Add({{ $k | printf "%q" }}, {{ $v | printf "%q" }}){{ end }}

	// Perform the HTTP request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Read the response body
	responseBody := new(bytes.Buffer)
	responseBody.ReadFrom(response.Body)

	// Print the response status code and body
	fmt.Println("Response Status Code:", response.Status)
	fmt.Println("Response Body:", responseBody.String())
}{{ if .RequestBody }}

// Request is model for the HTTP request body
{{ .RequestBody }}{{ end }}{{ if .ResponseBody }}

// Response is model for the HTTP response body
{{ .ResponseBody }}{{ end }}`
