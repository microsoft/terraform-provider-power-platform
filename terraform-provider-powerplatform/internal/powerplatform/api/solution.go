package powerplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

func (client *Client) ReadSolutions(environmentName string) ([]Solution, error) {
	var solutions []Solution

	request, err := http.NewRequest("GET",
		fmt.Sprintf("%s/api/environments/%s/solutions", client.HostURL, environmentName), nil)
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

func (client *Client) CreateSolutions(solutionToCreate Solution, content []byte, settings []byte) (*Solution, error) {

	body1, error := json.Marshal(SolutionCreate{
		SolutionName:    solutionToCreate.SolutionName,
		SolutionContent: content,
		SettingsContent: settings,
	})
	if error != nil {
		return nil, error
	}

	request, err := http.NewRequest("POST",
		fmt.Sprintf("%s/api/environments/%s/solutions", client.HostURL, solutionToCreate.EnvironmentName),
		bytes.NewReader(body1))
	if err != nil {
		return nil, err
	}

	_, fileName := filepath.Split(solutionToCreate.File)
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

	solution.EnvironmentName = solutionToCreate.EnvironmentName
	solution.File = solutionToCreate.File

	return &solution, nil
}
