package curl

import (
	"io"
	"net/http"
	"strings"
)

type Command struct {
	URL     string
	Headers map[string]string
	Body    string
	HasBody bool
	Method  string
}

func (c Command) ToRequest() (*http.Request, error) {
	req, err := http.NewRequest(c.Method, c.URL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}

	if c.HasBody {
		req.Body = io.NopCloser(strings.NewReader(c.Body))
	}

	return req, nil
}
