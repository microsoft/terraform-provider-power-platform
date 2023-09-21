package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type PowerPlatformApiRoundTripper struct {
	Transport     http.RoundTripper
	AccessToken   string
	RefreshToken  string
	Expiration    time.Time
	Context       context.Context
	ExpectedCodes []int
}

func NewPowerPlatformApiClient(ctx context.Context, expectedStatusCodes []int) *http.Client {
	expected := expectedStatusCodes
	if len(expected) == 0 {
		expected = []int{http.StatusOK}
	}

	return &http.Client{
		Transport: &PowerPlatformApiRoundTripper{
			Transport:     http.DefaultTransport,
			Context:       ctx,
			ExpectedCodes: expected,
		},
	}
}

func (c *PowerPlatformApiRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Set the authorization header with the token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
}

func (c *PowerPlatformApiRoundTripper) acquireToken() error {

}

func (rt *PowerPlatformApiRoundTripper) refreshToken() error {

}

func (c *PowerPlatformApiClient) Do(req *http.Request) (*http.Response, error) {
	// Set the authorization header with the token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check if the token needs to be refreshed
	if resp.StatusCode == http.StatusUnauthorized {
		// Refresh the token
		if err := c.refreshToken(); err != nil {
			return nil, err
		}

		// Set the new token in the authorization header
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

		// Make the request again with the new token
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (c *PowerPlatformApiClient) refreshToken() error {
	// Create a new HTTP client
	httpClient := &http.Client{}

	// Build the token request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id")
	data.Set("client_secret", "your_client_secret")
	data.Set("refresh_token", c.token)
	data.Set("resource", "https://api.powerplatform.com")

	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/common/oauth2/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the token request
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse the token response
	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return err
	}

	// Update the token in the client
	c.token = tokenResponse.AccessToken

	return nil

	return nil
}
