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

func (client *Client) doBasicAuth(username, password string) (*AuthResponse, error) {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth", client.HostURL), nil)
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

	return &authResponse, nil
}
