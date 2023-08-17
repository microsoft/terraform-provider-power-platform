package powerplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	HostURL    string
	HttpClient *http.Client
	AuthHash   string
}

// TODO support Oauth instead of basic auth
func NewClient(host, username, password string) (*Client, error) {
	client := Client{
		HttpClient: http.DefaultClient,
		HostURL:    host,
	}

	if username == "" {
		return &client, nil
	}

	// magodo: This is not a good practice to bind the authentication thing in the client builder. Instead, you might want to make it to be some kind of middleware.
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
	defer response.Body.Close()

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
	return body, nil
}
