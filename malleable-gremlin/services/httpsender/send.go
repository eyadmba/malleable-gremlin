package httpsender

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type SendArgument struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    interface{}       `json:"body,omitempty"`
}

type SendResult struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body"`
}

func Send(req *SendArgument) (*SendResult, error) {
	// Create a new request
	httpReq, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Set body if present
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}
		httpReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		httpReq.ContentLength = int64(len(bodyBytes))
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse response headers
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	// Try to parse body as JSON if possible
	var bodyInterface interface{}
	if err := json.Unmarshal(body, &bodyInterface); err != nil {
		bodyInterface = string(body)
	}

	return &SendResult{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       bodyInterface,
	}, nil
}
