package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
)

const (
	appFolder          = "goidcApp"
	tokensEndpoint     = "https://oauth2.googleapis.com/token"
	filesEndpoint      = "https://www.googleapis.com/drive/v3/files"
	uploadFileEndpoint = "https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart"
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

func GetAccessTokens(clientId, clientSecret, code, redirectUri, grantType string) (*AccessTokensResponse, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectUri)

	req, err := http.NewRequest("POST", tokensEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	input := &AccessTokensResponse{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v\nRaw body: %s", err, string(body))
	}
	return input, nil
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
	q := fmt.Sprintf(`name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false`, appFolder)
	url := fmt.Sprintf(`%s?q=%s`, filesEndpoint, url.QueryEscape(q))

	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)

	resp, err := do("GET", url, headers, nil)
	if err != nil {
		return false, "", fmt.Errorf("app folder exixts request failed: %w", err)
	}
	fmt.Println(string(resp))

	return true, string(resp), nil
}

func createAppFolder(accessToken string) (string, error) {
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)

	body := map[string]string{
		"name":     appFolder,
		"mimeType": "application/vnd.google-apps.folder",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal the request body for create app folder: %w", err)
	}

	resp, err := do("POST", filesEndpoint, headers, bytes.NewReader(jsonBody))
	fmt.Println(string(resp))
	if err != nil {
		return "", fmt.Errorf("create app folder request failed: %w", err)
	}
	return string(resp), nil
}

func ListFiles(appFolderId, accessToken string) ([]FileGD, error) {
	q := fmt.Sprintf(`%s+in+parents+and+trashed=false&fields=files(id,name,mimeType)`, appFolderId)
	url := fmt.Sprintf(`%s?q=%s`, filesEndpoint, url.QueryEscape(q))

	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)
	headers["Content-Type"] = "application/json"

	resp, err := do("GET", url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("list files request failed: %w", err)
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
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// creating the json part for the file metadata
	metaHeader := textproto.MIMEHeader{}
	metaHeader.Set("Content-Type", "application/json; charset=UTF-8")
	metaPart, err := writer.CreatePart(metaHeader)
	if err != nil {
		return fmt.Errorf("failed to create metadata part: %w", err)
	}
	metadata := map[string]interface{}{
		"name":    fileName,
		"parents": []string{appFolderId},
	}
	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	metaPart.Write(metaJSON)

	// creating the file part for the file contents
	filePart, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("failed to create file part: %w", err)
	}
	filePart.Write(fileData)

	writer.Close()

	req, err := http.NewRequest("POST", uploadFileEndpoint, buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("upload failed: %s\n%s", resp.Status, string(body))
	}

	fmt.Println(string(body))
	return nil
}

func DownloadFile(accessToken, fileId string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?alt=media", filesEndpoint, fileId)
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)
	resp, err := do("GET", url, headers, nil)
	fmt.Println(string(resp))
	if err != nil {
		return nil, fmt.Errorf("download request failed: %w", err)
	}
	return resp, nil
}

func DeleteFile(accessToken, fileId string) error {
	url := fmt.Sprintf("%s/%s", filesEndpoint, fileId)
	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", accessToken)
	resp, err := do("DELETE", url, headers, nil)
	fmt.Println(string(resp))
	if err != nil {
		return fmt.Errorf("delete file request failed: %w", err)
	}
	return nil

}
