package powerplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const HostUrl string = "http://localhost:8080"

type Client struct {
	HostURL    string
	HttpClient *http.Client
	AuthHash   string
}

// TODO support Oauth instead of basic auth
func NewClient(host, username, password *string) (*Client, error) {
	client := Client{
		HttpClient: &http.Client{Timeout: 1200 * time.Second}, //20 minutes
		HostURL:    HostUrl,
	}

	if host != nil {
		client.HostURL = *host
	}

	if username == nil || password == nil {
		return &client, nil
	}

	authResponse, error := client.doBasicAuth(username, password)
	if error != nil {
		return nil, error
	}
	client.AuthHash = authResponse.AuthHash

	return &client, nil
}

func (client *Client) doRequest(request *http.Request) ([]byte, error) {

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	if request.Header.Get("x-cred-hash") == "" {
		request.Header.Set("x-cred-hash", client.AuthHash)
	}

	response, err := client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if len(body) != 0 {
			errorResponse := make(map[string]interface{}, 0)
			err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&errorResponse)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("status: %d, body: %s", response.StatusCode, errorResponse)
		} else {
			return nil, fmt.Errorf("status: %d", response.StatusCode)
		}
	}
	defer response.Body.Close()
	return body, nil
}
