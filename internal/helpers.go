package service

import (
	"fmt"
	"io"
	"net/http"
)

func makeRequest(requestMethod, accessToken, url string) (*http.Request, error) {
	req, err := http.NewRequest(requestMethod, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	return req, nil
}

func do(requestMethod, accessToken, url string) ([]byte, error) {
	req, err := makeRequest(requestMethod, accessToken, url)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}

	return body, nil
}
