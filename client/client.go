package client

import (
	"io"
	"net/http"

	"github.com/ynori7/httpgen-go/curl"
)

// Client is a struct for the client
type Client struct {
	client *http.Client
}

// NewClient creates a new client
func NewClient() *Client {
	return &Client{
		client: &http.Client{},
	}
}

// Do performs the HTTP request
func (c *Client) Do(curlCommand curl.Command) (string, error) {
	req, err := curlCommand.ToRequest()
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
