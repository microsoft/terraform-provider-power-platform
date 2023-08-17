package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetEnvironments(ctx context.Context) ([]Environment, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/environments", client.BaseUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	envs := make([]Environment, 0)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&envs)
	if err != nil {
		return nil, err
	}

	return envs, nil
}

func (client *Client) GetEnvironment(ctx context.Context, id string) (*Environment, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/environments/%s", client.BaseUrl, id), nil)
	if err != nil {
		return nil, err
	}
	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *Client) DeleteEnvironment(ctx context.Context, id string) error {
	request, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/api/environments/%s", client.BaseUrl, id), nil)
	if err != nil {
		return err
	}
	_, err = client.doRequest(request)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) CreateEnvironment(ctx context.Context, envToCreate EnvironmentCreate) (*Environment, error) {
	body, err := json.Marshal(envToCreate)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/environments", client.BaseUrl), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *Client) UpdateEnvironment(ctx context.Context, id string, envToUpdate Environment) error {
	body, err := json.Marshal(envToUpdate)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/api/environments/%s", client.BaseUrl, id), bytes.NewReader(body))
	if err != nil {
		return err
	}
	_, err = client.doRequest(request)
	if err != nil {
		return err
	}
	return nil
}
