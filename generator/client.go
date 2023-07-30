package generator

import (
	"bufio"
	"bytes"
	"text/template"
)

// GenerateClientCode generates the client code which can be inserted into a function or method
func (h GoTemplate) GenerateClientCode() (string, error) {
	t := template.Must(template.New("txt").Parse(clientCodeTemplate))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	if err := t.Execute(w, h); err != nil {
		return "", err
	}

	w.Flush()
	return b.String(), nil
}

const clientCodeTemplate = `
// Create the HTTP request{{ if .Command.HasBody }}
body := bytes.NewReader([]byte({{ .Command.Body | printf "%q" }}))
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
`
