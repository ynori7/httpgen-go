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
			name:     "simple with null",
			jsonData: "{\"name\":\"John\",\"address\":null}",
			want:     "type Response struct { Address interface{} `json:\"address\"`\nName string `json:\"name\"`\n}\n",
		},
		{
			name:     "simple array",
			jsonData: "[\"San Francisco\",\"New York\"]",
			want:     "type Response struct {\nResponseItems []string `json:\"ResponseItems\"`}\n",
			wantErr:  false,
		},
		{
			name:     "nested",
			jsonData: "{\"name\":\"John\",\"address\":{\"street\":\"Easy St. 1\",\"city\":\"San Francisco\"}}",
			want:     "type Response struct {\nAddress struct {\nCity string `json:\"city\"`\nStreet string `json:\"street\"`\n} `json:\"address\"`\nName string `json:\"name\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "array of structs",
			jsonData: "{\"name\":\"John\",\"addresses\":[{\"city\":\"San Francisco\"},{\"city\":\"New York\"}]}",
			want:     "type Response struct {\nAddresses []struct{\nCity string `json:\"city\"`} `json:\"addresses\"`\nName string `json:\"name\"`\n}\n", 
			wantErr:  false,
		},
		{
			name:     "array of strings",
			jsonData: "{\"name\":\"John\",\"addresses\":[\"San Francisco\",\"New York\"]}",
			want:     "type Response struct {\nAddresses []string `json:\"addresses\"`\nName string `json:\"name\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "duplicate structs",
			jsonData: "{\"name\":\"John\",\"start\":{\"address\": {\"city\":\"San Francisco\"}},\"end\":{\"address\": {\"city\":\"New York\"}}} ",
			want:     "type Response struct {\nEnd struct {\nAddress struct {\nCity string `json:\"city\"`\n} `json:\"address\"`\n} `json:\"end\"`\nName string `json:\"name\"`\nStart struct {\nAddress struct {\nCity string `json:\"city\"`\n} `json:\"address\"`\n} `json:\"start\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "empty array",
			jsonData: "{\"name\":\"John\",\"addresses\":[]}",
			want:     "type Response struct {\nAddresses []interface{} `json:\"addresses\"`\nName string `json:\"name\"`\n}\n",
			wantErr:  false,
		},
		{
			name:     "error unmarshaling",
			jsonData: "{aasdfesses\":[]}",
			want:     "",
			wantErr:  true,
		},
		{
			name: "number as key",
			jsonData: "{\"fields\":{\"1\":{\"city\":\"San Francisco\"},\"2\":{\"city\":\"New York\"}}}",
			want: "type Response struct {\nFields struct {\nNum_1 struct {\nCity string `json:\"city\"`\n} `json:\"1\"`\nNum_2 struct {\nCity string `json:\"city\"`\n} `json:\"2\"`\n} `json:\"fields\"`\n}\n",
			wantErr: false,
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
			name:     "simple array",
			jsonData: "[{\"name\":\"John\"},{\"name\":\"Jane\"}]",
			want:     "type Response struct {\nResponseItems []ResponseItems `json:\"ResponseItems\"`\n}\n\ntype ResponseItems struct {\nName string `json:\"name\"`\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "simple",
			jsonData: "{\"name\":\"John\"}",
			want:     "type Response struct { Name string `json:\"name\"`\n}\n\n",
			wantErr:  false,
		},
		{
			name:     "simple with null",
			jsonData: "{\"name\":\"John\",\"address\":null}",
			want:     "type Response struct { Address interface{} `json:\"address\"`\nName string `json:\"name\"`\n}\n\n",
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
		{
			name: "number as key",
			jsonData: "{\"fields\":{\"1\":{\"city\":\"San Francisco\"},\"2\":{\"city\":\"New York\"}}}",
			want: "type Response struct {\nFields Fields `json:\"fields\"`\n}\n\ntype Fields struct {\nFieldsNum_1 FieldsNum_1 `json:\"1\"`\nFieldsNum_2 FieldsNum_2 `json:\"2\"`\n}\n\ntype FieldsNum_1 struct {\nCity string `json:\"city\"`\n}\n\ntype FieldsNum_2 struct {\nCity string `json:\"city\"`\n}\n\n",
			wantErr: false,
		},
		{
			name: "title has special characters",
			jsonData: "{\"@type\":\"Person\",\"name\":\"John\"}",
			want: "type Response struct {\nName string `json:\"name\"`\nType string `json:\"@type\"`\n}\n\n",
			wantErr: false,
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
