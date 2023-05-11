package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetEnvironments() ([]Environment, error) {
	request, error := http.NewRequest("GET", fmt.Sprintf("%s/api/environments", client.HostURL), nil)
	if error != nil {
		return nil, error
	}
	body, error := client.doRequest(request)
	if error != nil {
		return nil, error
	}

	envs := make([]Environment, 0)
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&envs)
	if error != nil {
		return nil, error
	}

	return envs, nil
}

func (client *Client) GetEnvironment(ctx context.Context, id string) (*Environment, error) {
	request, error := http.NewRequest("GET", fmt.Sprintf("%s/api/environments/%s", client.HostURL, id), nil)
	if error != nil {
		return nil, error
	}
	request = request.WithContext(ctx)
	body, error := client.doRequest(request)
	if error != nil {
		return nil, error
	}

	env := Environment{}
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&env)
	if error != nil {
		return nil, error
	}

	return &env, nil
}

func (client *Client) DeleteEnvironment(ctx context.Context, id string) error {
	request, error := http.NewRequest("DELETE", fmt.Sprintf("%s/api/environments/%s", client.HostURL, id), nil)
	if error != nil {
		return error
	}
	request = request.WithContext(ctx)
	_, error = client.doRequest(request)
	if error != nil {
		return error
	}

	return nil
}

func (client *Client) CreateEnvironment(ctx context.Context, envToCreate EnvironmentCreate) (*Environment, error) {
	body, error := json.Marshal(envToCreate)
	if error != nil {
		return nil, error
	}
	request, error := http.NewRequest("POST", fmt.Sprintf("%s/api/environments", client.HostURL), bytes.NewReader(body))
	if error != nil {
		return nil, error
	}
	request = request.WithContext(ctx)
	body, error = client.doRequest(request)
	if error != nil {
		return nil, error
	}

	env := Environment{}
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&env)
	if error != nil {
		return nil, error
	}

	return &env, nil
}

func (client *Client) UpdateEnvironment(ctx context.Context, id string, envToUpdate Environment) error {
	body, error := json.Marshal(envToUpdate)
	if error != nil {
		return error
	}
	request, error := http.NewRequest("PUT", fmt.Sprintf("%s/api/environments/%s", client.HostURL, id), bytes.NewReader(body))
	if error != nil {
		return error
	}
	request = request.WithContext(ctx)
	_, error = client.doRequest(request)
	if error != nil {
		return error
	}
	return nil
}
