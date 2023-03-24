package powerplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (Client *Client) DeleteUser(environmentName string, aadId string) error {

	request, error := http.NewRequest("DELETE", fmt.Sprintf("%s/api/environments/%s/users/%s", Client.HostURL, environmentName, aadId), nil)
	if error != nil {
		return error
	}
	_, error = Client.doRequest(request)
	if error != nil {
		return error
	}
	return nil
}

func (client *Client) UpdateUser(environmentName string, userToUpdate User) (*User, error) {
	body, error := json.Marshal(userToUpdate)
	if error != nil {
		return nil, error
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/environments/%s/users/%s", client.HostURL, environmentName, userToUpdate.AadObjectId), bytes.NewReader(body))
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

func (client *Client) CreateUser(environmentName string, userToCreate User) (*User, error) {

	body, error := json.Marshal(userToCreate)
	if error != nil {
		return nil, error
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/environments/%s/users", client.HostURL, environmentName), bytes.NewReader(body))
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

func (client *Client) ReadUser(environmentName string, aadId string) (*User, error) {
	var user User

	request, err := http.NewRequest("GET",
		fmt.Sprintf("%s/api/environments/%s/users/%s", client.HostURL, environmentName, aadId), nil)
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

func (client *Client) ReadUsers(environmentName string) ([]User, error) {
	var users []User

	request, err := http.NewRequest("GET",
		fmt.Sprintf("%s/api/environments/%s/users", client.HostURL, environmentName), nil)
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
