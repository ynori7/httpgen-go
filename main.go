package main

import (
	"fmt"

	"github.com/ynori7/httpgen-go/client"
	"github.com/ynori7/httpgen-go/curl"
	"github.com/ynori7/httpgen-go/generator"
)

func main() {
	curlCommand := `curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer myToken" -d '{"name":"John", "age":30}' https://example.com/api/endpoint`

	command, err := curl.Parse(curlCommand)
	if err != nil {
		fmt.Println(err)
		return
	}

	cli := client.NewClient()
	resp, err := cli.Do(*command)
	if err != nil {
		fmt.Println(err)
		return
	}

	generator := generator.NewGoTemplate(command, resp)
	goCode, err := generator.ExecuteGoTemplate(false) //TODO: make this flag configurable
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the generated Go code
	fmt.Println(goCode)
}
