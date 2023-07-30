package generator

import (
	"bufio"
	"bytes"
	"go/format"
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

	clientCode, err := h.GenerateClientCode()
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("txt").Parse(goTemplate))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	if err = t.Execute(w, struct {
		ClientCode   string
		RequestBody  string
		ResponseBody string
	}{
		ClientCode:   clientCode,
		RequestBody:  h.RequestBody,
		ResponseBody: h.ResponseBody,
	}); err != nil {
		return "", err
	}

	w.Flush()

	formattedCode, err := format.Source(b.Bytes())
	if err != nil {
		return b.String(), err
	}
	return string(formattedCode), nil
}

const goTemplate = `package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() { {{ .ClientCode }}
fmt.Println("Response Body:", responseBody.String())
}{{ if .RequestBody }}

// Request is model for the HTTP request body
{{ .RequestBody }}{{ end }}{{ if .ResponseBody }}

// Response is model for the HTTP response body
{{ .ResponseBody }}{{ end }}`
