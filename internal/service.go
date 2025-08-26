package service

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	appFolder      = "goidcApp"
	tokensEndpoint = "https://oauth2.googleapis.com/token"
	filesEndpoint  = "https://www.googleapis.com/drive/v3/files"
)

type FileGD struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
}

type AccessTokensResponse struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func GetAccessTokens(clientId, clientSecret, code, redirectUri, grantType string) (AccessTokensResponse, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	resp, err := do("POST", tokensEndpoint, headers, nil)
}

func CreateAppFolderIfNotExist(accessToken string) (string, error) {
	exists, folderId, err := checkAppFolderExists(accessToken)
	if err != nil {
		return "", err
	}

	if !exists {
		return createAppFolder(accessToken)
	}
	return folderId, nil
}

func checkAppFolderExists(accessToken string) (bool, string, error) {
	q := `name='example_folder' and mimeType='application/vnd.google-apps.folder' and trashed=false`
	url := fmt.Sprintf(`%s?q=%s`, filesEndpoint, url.QueryEscape(q))

	headers := make(map[string]string)

	resp, err := do("GET", accessToken, url, headers, nil)
	if err != nil {
		return false, "", err
	}
	fmt.Println(string(resp))

	return true, string(resp), nil
}

func createAppFolder(accessToken string) (string, error) {
	headers := make(map[string]string)
	resp, err := do("POST", accessToken, filesEndpoint, headers, nil)
	fmt.Println(string(resp))
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func ListFiles(appFolderId, accessToken string) ([]FileGD, error) {
	q := fmt.Sprintf(`%s+in+parents+and+trashed=false&fields=files(id,name,mimeType)`, appFolderId)
	url := fmt.Sprintf(`%s?q=%s`, filesEndpoint, url.QueryEscape(q))
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	resp, err := do("GET", accessToken, url, headers, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))

	var userFiles struct {
		Files []FileGD `json:"files"`
	}

	if err := json.Unmarshal(resp, &userFiles); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v\nRaw body: %s", err, string(resp))
	}
	return userFiles.Files, nil
}

func UploadFile(accessToken, appFolderId, fileName string, fileData []byte) error {
	// TODO make a request for multipart file upload

	headers := make(map[string]string)
	resp, err := do("POST", accessToken, filesEndpoint, headers, nil)
	fmt.Println(string(resp))
	return err
}

func DownloadFile(accessToken, fileId string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?alt=media", filesEndpoint, fileId)
	headers := make(map[string]string)
	resp, err := do("GET", accessToken, url, headers, nil)
	fmt.Println(string(resp))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DeleteFile(accessToken, fileId string) error {
	url := fmt.Sprintf("%s/%s", filesEndpoint, fileId)
	headers := make(map[string]string)
	resp, err := do("DELETE", accessToken, url, headers, nil)
	fmt.Println(string(resp))
	return err
}
