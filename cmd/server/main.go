package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ynori7/httpgen-go/client"
	"github.com/ynori7/httpgen-go/curl"
	"github.com/ynori7/httpgen-go/generator"
	"github.com/ynori7/httpgen-go/structs"
)

var httpClient = client.NewClient()

func main() {
	fmt.Println("Starting server...")
	fmt.Println("Visit http:/localhost:8080...")
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the form is submitted
	output, err := getOutput(r)
	if err != nil {
		output = fmt.Sprintf("Error: %v", err)
	}

	// Define the HTML template
	const htmlTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Parse Curl Command</title>
	</head>
	<body>
		<h1>Simple Curl Command</h1>
		<p>Enter a curl command to generate a simple Go client or JSON payload to generate the models.</p>
		<form method="post">
			<label for="curlCommand">Curl Command:</label>
			<br>
			<textarea id="curlCommand" name="curlCommand" rows="5" cols="50"></textarea>
			<br>
			<label for="json">JSON Payload:</label>
			<br>
			<textarea id="json" name="json" rows="5" cols="50"></textarea>
			<br>
			<input type="checkbox" id="inline" name="inline">
			<label for="inline">Inline</label>
			<br>
			<input type="submit" value="Submit">
		</form>
		<br>
		{{ if .Output }}<h2>Output:</h2>
		<pre>{{.Output}}</pre>{{ end }}
	</body>
	</html>
	`

	// Parse the HTML template
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template and pass an empty output initially
	data := struct {
		Output string
	}{Output: output}
	tmpl.Execute(w, data)
}

func getOutput(r *http.Request) (string, error) {
	if r.Method == http.MethodPost {
		// Get the form values
		curlCommand := r.FormValue("curlCommand")
		json := r.FormValue("json")
		inline := r.FormValue("inline")
		useInline := inline != ""

		if curlCommand == "" && json == "" {
			return "", fmt.Errorf("curl command or json payload must be provided")
		}
		if json != "" {
			data, err := structs.CreateStructFromJSON(json, "Payload", useInline)
			if err != nil {
				fmt.Println("Error creating structs: " + err.Error())
				return "", err
			}
			return data, nil
		}

		command, err := curl.Parse(curlCommand)
		if err != nil {
			fmt.Println("Error parsing curl request: " + err.Error())
			return "", err
		}

		resp, err := httpClient.Do(*command)
		if err != nil {
			fmt.Println("Error doing http request: " + err.Error())
			return "", err
		}

		generator := generator.NewGoTemplate(command, resp)
		goCode, err := generator.ExecuteGoTemplate(useInline)
		if err != nil {
			fmt.Println("Error generating template: " + err.Error())
			return "", err
		}

		return goCode, nil
	}

	return "", nil
}
