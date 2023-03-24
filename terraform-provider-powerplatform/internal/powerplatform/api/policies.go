package powerplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetPolicies() ([]DlpPolicy, error) {
	request, error := http.NewRequest("GET", fmt.Sprintf("%s/api/policies", client.HostURL), nil)
	if error != nil {
		return nil, error
	}
	body, error := client.doRequest(request)
	if error != nil {
		return nil, error
	}

	policies := make([]DlpPolicy, 0)
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&policies)
	if error != nil {
		return nil, error
	}
	return policies, nil
}

func (client *Client) GetPolicy(name string) (*DlpPolicy, error) {
	request, error := http.NewRequest("GET", fmt.Sprintf("%s/api/policies/%s", client.HostURL, name), nil)
	if error != nil {
		return nil, error
	}
	body, error := client.doRequest(request)
	if error != nil {
		return nil, error
	}

	policy := DlpPolicy{}
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&policy)
	if error != nil {
		return nil, error
	}
	return &policy, nil
}

func (client *Client) DeletePolicy(name string) error {
	request, error := http.NewRequest("DELETE", fmt.Sprintf("%s/api/policies/%s", client.HostURL, name), nil)
	if error != nil {
		return error
	}
	_, error = client.doRequest(request)
	if error != nil {
		return error
	}
	return nil
}

func (client *Client) UpdatePolicy(name string, policyToUpdate DlpPolicy) (*DlpPolicy, error) {

	body, error := json.Marshal(policyToUpdate)
	if error != nil {
		return nil, error
	}

	request, error := http.NewRequest("PUT", fmt.Sprintf("%s/api/policies/%s", client.HostURL, name), bytes.NewReader(body))
	if error != nil {
		return nil, error
	}

	body, error = client.doRequest(request)
	if error != nil {
		return nil, error
	}

	policy := DlpPolicy{}
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&policy)
	if error != nil {
		return nil, error
	}

	return &policy, nil
}

func (client *Client) CreatePolicy(policyToCreate DlpPolicy) (*DlpPolicy, error) {

	body, error := json.Marshal(policyToCreate)
	if error != nil {
		return nil, error
	}

	request, error := http.NewRequest("POST", fmt.Sprintf("%s/api/policies", client.HostURL), bytes.NewReader(body))
	if error != nil {
		return nil, error
	}

	body, error = client.doRequest(request)
	if error != nil {
		return nil, error
	}

	policy := DlpPolicy{}
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&policy)
	if error != nil {
		return nil, error
	}

	return &policy, nil
}
