package curl

import "testing"

func TestParse(t *testing.T) {
	testcases := []struct {
		name    string
		input   string
		want    Command
		wantErr bool
	}{
		{
			name:    "empty",
			input:   "",
			want:    Command{},
			wantErr: true,
		},
		{
			name:    "simple",
			input:   "curl http://example.com",
			want:    Command{URL: "http://example.com", Method: "GET"},
			wantErr: false,
		},
		{
			name:    "simple with method",
			input:   "curl -X POST http://example.com",
			want:    Command{URL: "http://example.com", Method: "POST"},
			wantErr: false,
		},
		{
			name:    "simple with method and body",
			input:   "curl -X POST -d \"hello\" http://example.com",
			want:    Command{URL: "http://example.com", Method: "POST", Body: "hello", HasBody: true},
			wantErr: false,
		},
		{
			name:    "simple with method and body and header",
			input:   "curl -X POST -d \"hello\" -H \"Content-Type: application/json\" http://example.com",
			want:    Command{URL: "http://example.com", Method: "POST", Body: "hello", HasBody: true, Headers: map[string]string{"Content-Type": "application/json"}},
			wantErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil && !tc.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && tc.wantErr {
				t.Errorf("expected error but got none")
			}
			if tc.wantErr && err != nil {
				return
			}
			if got.URL != tc.want.URL {
				t.Errorf("expected url %s but got %s", tc.want.URL, got.URL)
			}
			if got.Method != tc.want.Method {
				t.Errorf("expected method %s but got %s", tc.want.Method, got.Method)
			}
			if got.Body != tc.want.Body {
				t.Errorf("expected body %s but got %s", tc.want.Body, got.Body)
			}
			if got.HasBody != tc.want.HasBody {
				t.Errorf("expected hasBody %v but got %v", tc.want.HasBody, got.HasBody)
			}
			if len(got.Headers) != len(tc.want.Headers) {
				t.Errorf("expected %d headers but got %d", len(tc.want.Headers), len(got.Headers))
			}
			for k, v := range tc.want.Headers {
				if got.Headers[k] != v {
					t.Errorf("expected header %s to be %s but got %s", k, v, got.Headers[k])
				}
			}
		})
	}

}
