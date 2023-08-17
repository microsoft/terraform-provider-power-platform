package powerplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthResponse struct {
	AuthHash string `json:"auth_hash"`
}

func (client *Client) DoBasicAuth(baseUrl, username, password string) (*AuthResponse, error) {
	client.HttpClient = http.DefaultClient
	client.BaseUrl = baseUrl

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth", client.BaseUrl), nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(username, password)
	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}
	authResponse := AuthResponse{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&authResponse)

	client.AuthHash = authResponse.AuthHash

	return &authResponse, nil
}
