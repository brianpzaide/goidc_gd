package service

import (
	"fmt"
	"io"
	"net/http"
)

func do(requestMethod, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(requestMethod, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}
