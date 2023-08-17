package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

func (client *Client) DeleteSolution(ctx context.Context, environmentName string, solutionName string) error {

	request, err := http.NewRequestWithContext(ctx, "DELETE",
		fmt.Sprintf("%s/api/environments/%s/solutions/%s", client.BaseUrl, environmentName, solutionName), nil)
	if err != nil {
		return err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		return fmt.Errorf("error from api when deleting solution: %s", string(body))
	}

	return nil
}

func (client *Client) GetSolutions(ctx context.Context, environmentName string) ([]Solution, error) {
	var solutions []Solution

	request, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/api/environments/%s/solutions", client.BaseUrl, environmentName), nil)
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

	err = json.Unmarshal(body, &solutions)
	if err != nil {
		return nil, err
	}

	return solutions, nil
}

func (client *Client) CreateSolution(ctx context.Context, EnvironmentName string, solutionToCreate Solution, content []byte, settings []byte) (*Solution, error) {

	body1, err := json.Marshal(SolutionCreate{
		SolutionName:    solutionToCreate.SolutionName,
		SolutionContent: content,
		SettingsContent: settings,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("%s/api/environments/%s/solutions", client.BaseUrl, EnvironmentName),
		bytes.NewReader(body1))
	if err != nil {
		return nil, err
	}

	_, fileName := filepath.Split(solutionToCreate.SolutionFile)
	request.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	request.Header.Set("X-Ms-Solution-Name", solutionToCreate.SolutionName)

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}
	if body == nil {
		return nil, fmt.Errorf("no body returned")
	}

	solution := Solution{}
	err = json.Unmarshal(body, &solution)
	if err != nil {
		return nil, err
	}

	solution.SolutionFile = solutionToCreate.SolutionFile

	return &solution, nil
}
