package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/BRO3886/gtasks/internal/config"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// GetService creates a Google Tasks service client
func GetService() (*tasks.Service, error) {
	oauthConfig, err := config.GetOAuth2Config()
	if err != nil {
		return nil, fmt.Errorf("failed to get OAuth2 config: %v", err)
	}

	client, err := getClient(oauthConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %v", err)
	}

	srv, err := tasks.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Tasks service: %v", err)
	}

	return srv, nil
}

// getClient retrieves HTTP client with valid token
func getClient(oauthConfig *oauth2.Config) (*http.Client, error) {
	folderPath := config.GetInstallLocation()
	tokFile := folderPath + "/token.json"

	token, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, fmt.Errorf("not authenticated. Run 'gtasks login' first")
	}

	return oauthConfig.Client(context.Background(), token), nil
}

// tokenFromFile retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	return &token, err
}
