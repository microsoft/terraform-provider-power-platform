package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (Client *Client) DeleteUser(ctx context.Context, environmentName string, aadId string) error {

	request, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/api/environments/%s/users/%s", Client.BaseUrl, environmentName, aadId), nil)
	if err != nil {
		return err
	}
	_, err = Client.doRequest(request)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) UpdateUser(ctx context.Context, environmentName string, userToUpdate User) (*User, error) {
	body, err := json.Marshal(userToUpdate)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/api/environments/%s/users/%s", client.BaseUrl, environmentName, userToUpdate.AadObjectId), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}

	updatedUser := User{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (client *Client) CreateUser(ctx context.Context, environmentName string, userToCreate User) (*User, error) {

	body, err := json.Marshal(userToCreate)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/environments/%s/users", client.BaseUrl, environmentName), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, fmt.Errorf("no body returned")
	}

	createdUser := User{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&createdUser)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (client *Client) GetUser(ctx context.Context, environmentName string, aadId string) (*User, error) {
	var user User

	request, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/api/environments/%s/users/%s", client.BaseUrl, environmentName, aadId), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, fmt.Errorf("no body returned")
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (client *Client) GetUsers(ctx context.Context, environmentName string) ([]User, error) {
	var users []User

	request, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/api/environments/%s/users", client.BaseUrl, environmentName), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, fmt.Errorf("no body returned")
	}

	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
