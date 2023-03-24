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

func (client *Client) doBasicAuth(username, password *string) (*AuthResponse, error) {
	request, error := http.NewRequest("POST", fmt.Sprintf("%s/api/auth", client.HostURL), nil)
	if error != nil {
		return nil, error
	}
	request.SetBasicAuth(*username, *password)
	body, error := client.doRequest(request)
	if error != nil {
		return nil, error
	}
	authResponse := AuthResponse{}
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&authResponse)

	return &authResponse, nil
}
