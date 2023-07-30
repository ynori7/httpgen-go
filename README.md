# HttpGen-Go
This tool / library allows you to generate HTTP client code based on cURL requests. It can take a cURL request, generate the HTTP request with headers, generate the request model, and fetch an example response and generate the response model.

### How this tool differs from others
There are other similar tools out there already for generating go code from cURL requests or generating the models from a JSON. What's special about httpgen-go is that it's able to do both with the same tool. Additionally it automatically fetches a response using the generated client code to see if it can auto-generate a response model. What also differs with this tool is that it's written in Go, making it easy to import as a library to extend so that you can build custom code generators.

### Parse cURL 
The `curl` package provides tools for parsing a cURL command so that an HTTP request can ge created from it. You can call it like this:

```go
curlCommand := `curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer myToken" -d '{"name":"John", "age":30}' https://example.com/api/endpoint`

command, err := curl.Parse(curlCommand)
```

The resulting `command` will contain the URL, method, body, and headers.

### Parse JSON
The `structs` package provides tools for parsing a JSON into Go structures. You can call it like this:

```go
structs.CreateStructFromJSON(jsonString, "Request", true)
```

The first parameter is the actual JSON to parse, the second is the name of the top-level structure, and the third is a flag indicating whether nested structs should be created inline or not.

### Commands
In `cmd`, you can find the following commands:

**cmd/server**
This command starts a simple HTTP server with a basic web UI for entering the cURL commands to parse and displaying the results. Just run `go run cmd/server/main.go` and then open http://localhost:8080