package gollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// apiGet sends a GET request to the specified path on the Ollama server,
// and unmarshals the response into the given interface.
//
// The URL is built by joining the server address with the path.
//
// The Ollama server must respond with a JSON object that can be
// unmarshaled into the given interface.
//
// The Verbose flag is respected, and the URL is printed if it is set.
//
// The HTTPTimeout is used as the timeout for the HTTP request.
//
// If the request fails, or the response cannot be unmarshaled, an error
// is returned.
func (c *Gollama) apiGet(ctx context.Context, path string, v interface{}) error {
	url, _ := url.JoinPath(c.ServerAddr, path)
	if c.Verbose {
		fmt.Printf("Sending a request to GET %s\n", url)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	HTTPClient := &http.Client{
		Timeout: c.HTTPTimeout,
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

// apiPost sends a POST request to the specified path on the Ollama server,
// and unmarshals the response into the given interface.
//
// The URL is built by joining the server address with the path.
//
// The Ollama server must respond with a JSON object that can be
// unmarshaled into the given interface.
//
// The Verbose flag is respected, and the URL is printed if it is set.
//
// If the request fails, or the response cannot be unmarshaled, an error
// is returned.
//
// The HTTPTimeout is used as the timeout for the HTTP request, except for
// requests to the /api/pull endpoint, which is given the PullTimeout.
func (c *Gollama) apiPost(ctx context.Context, path string, v interface{}, data interface{}) error {
	url, _ := url.JoinPath(c.ServerAddr, path)
	if c.Verbose {
		fmt.Printf("Sending a request to POST %s\n", url)
	}

	reqBytes, err := json.Marshal(data)
	if err != nil {
		if c.Verbose {
			fmt.Printf("Failed to marshal request data: %s\n", err)
		}
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	HTTPClient := &http.Client{
		Timeout: c.HTTPTimeout,
	}

	if path == "/api/pull" {
		HTTPClient.Timeout = c.PullTimeout
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		if c.Verbose {
			fmt.Printf("Failed to send request: %s\n", err)
		}
		return err
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		if c.Verbose {
			fmt.Println("Failed to read response body error:", err)
		}
		return err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		if c.Verbose {
			fmt.Printf("Failed to unmarshal response body: %s\n", err)
		}
		return err
	}

	if c.Verbose {
		if bodyBytes != nil {
			fmt.Printf("Response body: %s\n", string(bodyBytes))
		} else {
			fmt.Println("Response body logging skipped due to error or disabled verbose mode.")
		}
	}

	return nil
}
