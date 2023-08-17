package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetPolicies(ctx context.Context) ([]DlpPolicy, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/policies", client.BaseUrl), nil)
	if err != nil {
		return nil, err
	}
	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	policies := make([]DlpPolicy, 0)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&policies)
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// TODO support not found response
func (client *Client) GetPolicy(ctx context.Context, name string) (*DlpPolicy, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/policies/%s", client.BaseUrl, name), nil)
	if err != nil {
		return nil, err
	}
	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	policy := DlpPolicy{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (client *Client) DeletePolicy(ctx context.Context, name string) error {
	request, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/api/policies/%s", client.BaseUrl, name), nil)
	if err != nil {
		return err
	}
	_, err = client.doRequest(request)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) UpdatePolicy(ctx context.Context, name string, policyToUpdate DlpPolicy) (*DlpPolicy, error) {

	body, err := json.Marshal(policyToUpdate)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/api/policies/%s", client.BaseUrl, name), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}

	policy := DlpPolicy{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (client *Client) CreatePolicy(ctx context.Context, policyToCreate DlpPolicy) (*DlpPolicy, error) {

	body, err := json.Marshal(policyToCreate)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/policies", client.BaseUrl), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}

	policy := DlpPolicy{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}
