package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var _ ClientInterface = &Client{}

//go:generate mockgen -destination=../../mocks/client_mocks.go -package=powerplatform_mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api" ClientInterface
type ClientInterface interface {
	DoBasicAuth(baseUrl, username, password string) (*AuthResponse, error)

	GetPowerApps(ctx context.Context, environmentName string) ([]App, error)

	GetEnvironments(ctx context.Context) ([]Environment, error)

	GetEnvironment(ctx context.Context, id string) (*Environment, error)
	DeleteEnvironment(ctx context.Context, id string) error
	CreateEnvironment(ctx context.Context, envToCreate EnvironmentCreate) (*Environment, error)
	UpdateEnvironment(ctx context.Context, id string, envToUpdate Environment) error

	GetPolicies(ctx context.Context) ([]DlpPolicy, error)
	GetPolicy(ctx context.Context, name string) (*DlpPolicy, error)
	DeletePolicy(ctx context.Context, name string) error
	UpdatePolicy(ctx context.Context, name string, policyToUpdate DlpPolicy) (*DlpPolicy, error)
	CreatePolicy(ctx context.Context, policyToCreate DlpPolicy) (*DlpPolicy, error)

	DeleteSolution(ctx context.Context, environmentName string, solutionName string) error
	GetSolutions(ctx context.Context, environmentName string) ([]Solution, error)
	CreateSolution(ctx context.Context, EnvironmentName string, solutionToCreate Solution, content []byte, settings []byte) (*Solution, error)

	DeleteUser(ctx context.Context, environmentName string, aadId string) error
	UpdateUser(ctx context.Context, environmentName string, userToUpdate User) (*User, error)
	CreateUser(ctx context.Context, environmentName string, userToCreate User) (*User, error)
	GetUser(ctx context.Context, environmentName string, aadId string) (*User, error)
	GetUsers(ctx context.Context, environmentName string) ([]User, error)

	GetConnectors(ctx context.Context) ([]Connector, error)
}
type Client struct {
	BaseUrl    string
	HttpClient *http.Client
	AuthHash   string
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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

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
