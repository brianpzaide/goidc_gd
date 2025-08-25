package service

import (
	"fmt"
	"net/url"
)

const (
	appFolder      = "goidcApp"
	tokensEndpoint = "https://oauth2.googleapis.com/token"
	filesEndpoint  = "https://www.googleapis.com/drive/v3/files"
)

func GetAccessTokens(clientId, clientSecret, redirectUri, grantType string) {

}

func CreateAppFolderIfNotExist(accessToken string) error {
	exists, err := checkAppFolderExists(accessToken)
	if err != nil {
		return err
	}

	if !exists {
		return createAppFolder(accessToken)
	}
	return nil
}

func checkAppFolderExists(accessToken string) (bool, error) {
	q := `name='example_folder' and mimeType='application/vnd.google-apps.folder' and trashed=false`
	url := fmt.Sprintf(`%s?q=%s`, filesEndpoint, url.QueryEscape(q))

	resp, err := do("GET", accessToken, url)
	if err != nil {
		return false, err
	}
	fmt.Println(string(resp))

	return false, nil
}

func createAppFolder(accessToken string) error {
	resp, err := do("POST", accessToken, filesEndpoint)
	fmt.Println(string(resp))
	return err
}
