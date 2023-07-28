package structs

import (
	"go/format"
	"testing"
)

func TestCreateStructFromJSON_Inline(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     string
		wantErr  bool
	}{
		{
			name:     "empty",
			jsonData: "{}",
			want:     "type Response struct {\n}\n",
			wantErr:  false,
		},
		{
			name:     "simple",
			jsonData: "{\"name\":\"John\"}",
			want:     "type Response struct { Name string `json:\"name\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "nested",
			jsonData: "{\"name\":\"John\",\"address\":{\"city\":\"San Francisco\"}}",
			want:     "type Response struct {\nName string `json:\"name\"`\nAddress struct {\nCity string `json:\"city\"`\n} `json:\"address\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "array of structs",
			jsonData: "{\"name\":\"John\",\"addresses\":[{\"city\":\"San Francisco\"},{\"city\":\"New York\"}]}",
			want:     "type Response struct {\nName string `json:\"name\"`\nAddresses []map[string]interface{} `json:\"addresses\"`\n}\n", //TODO: improve this
			wantErr:  false,
		},
		{
			name:     "array of strings",
			jsonData: "{\"name\":\"John\",\"addresses\":[\"San Francisco\",\"New York\"]}",
			want:     "type Response struct {\nName string `json:\"name\"`\nAddresses []string `json:\"addresses\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "duplicate structs",
			jsonData: "{\"name\":\"John\",\"start\":{\"address\": {\"city\":\"San Francisco\"}},\"end\":{\"address\": {\"city\":\"New York\"}}} ",
			want:     "type Response struct {\nName string `json:\"name\"`\nStart struct {\nAddress struct {\nCity string `json:\"city\"`\n} `json:\"address\"`\n} `json:\"start\"`\nEnd struct {\nAddress struct {\nCity string `json:\"city\"`\n} `json:\"address\"`\n} `json:\"end\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "empty array",
			jsonData: "{\"name\":\"John\",\"addresses\":[]}",
			want:     "type Response struct {\nName string `json:\"name\"`\nAddresses []interface{} `json:\"addresses\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "error unmarshaling",
			jsonData: "{aasdfesses\":[]}",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateStructFromJSON(tt.jsonData, "Response", true)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStructFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			formattedCode, err := format.Source([]byte(tt.want))
			if err != nil {
				t.Errorf("error formatting expected code: %v", err)
				return
			}
			if got != string(formattedCode) {
				t.Errorf("\nCreateStructFromJSON() = \n%v, want \n%v", got, string(formattedCode))
			}
		})
	}
}

func TestCreateStructFromJSON_Multi(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		want     string
		wantErr  bool
	}{
		{
			name:     "empty",
			jsonData: "{}",
			want:     "type Response struct {\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "simple",
			jsonData: "{\"name\":\"John\"}",
			want:     "type Response struct { Name string `json:\"name\"`\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "nested",
			jsonData: "{\"name\":\"John\",\"address\":{\"city\":\"San Francisco\"}}",
			want:     "type Response struct {\nAddress Address `json:\"address\"`\nName string `json:\"name\"`\n}\n\ntype Address struct {\nCity string `json:\"city\"`\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "array of structs",
			jsonData: "{\"name\":\"John\",\"addresses\":[{\"city\":\"San Francisco\"},{\"city\":\"New York\"}]}",
			want:     "type Response struct {\nAddresses []Addresses `json:\"addresses\"`\nName string `json:\"name\"`\n}\n\ntype Addresses struct {\nCity string `json:\"city\"`\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "array of strings",
			jsonData: "{\"name\":\"John\",\"addresses\":[\"San Francisco\",\"New York\"]}",
			want:     "type Response struct {\nAddresses []string `json:\"addresses\"`\nName string `json:\"name\"`\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "duplicate structs",
			jsonData: "{\"name\":\"John\",\"start\":{\"address\": {\"city\":\"San Francisco\"}},\"end\":{\"address\": {\"city\":\"New York\"}}} ",
			want:     "type Response struct {\nEnd End `json:\"end\"`\nName string `json:\"name\"`\nStart Start `json:\"start\"`\n}\n\ntype Address struct {\nCity string `json:\"city\"`\n}\n\ntype End struct {\nAddress Address `json:\"address\"`\n}\n\ntype Start struct {\nAddress Address `json:\"address\"`\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "duplicate structs with different fields",
			jsonData: "{\"name\":\"John\",\"start\":{\"address\": {\"city\":\"San Francisco\"}},\"end\":{\"address\": {\"city\":\"New York\", \"zip\": \"12345\"}}} ",
			want:     "type Response struct {\nEnd End `json:\"end\"`\nName string `json:\"name\"`\nStart Start `json:\"start\"`\n}\n\ntype Address struct {\nCity string `json:\"city\"`\n}\n\ntype Address1 struct {\nCity string `json:\"city\"`\nZip string `json:\"zip\"`\n}\n\ntype End struct {\nAddress Address1 `json:\"address\"`\n}\n\ntype Start struct {\nAddress Address `json:\"address\"`\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "error unmarshaling",
			jsonData: "{aasdfesses\":[]}",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateStructFromJSON(tt.jsonData, "Response", false)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStructFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			formattedCode, err := format.Source([]byte(tt.want))
			if err != nil {
				t.Errorf("error formatting expected code: %v", err)
				return
			}
			if got != string(formattedCode) {
				t.Errorf("\nCreateStructFromJSON() = \n%v, want \n%v", got, string(formattedCode))
			}
		})
	}
}
