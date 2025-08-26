package service

import (
	"fmt"
	"io"
	"net/http"
)

func makeRequest(requestMethod, accessToken, url string, headers map[string]string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(requestMethod, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

func do(requestMethod, accessToken, url string, headers map[string]string, body io.Reader) ([]byte, error) {
	req, err := makeRequest(requestMethod, accessToken, url, headers, body)
	if err != nil {
		return nil, err
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}
