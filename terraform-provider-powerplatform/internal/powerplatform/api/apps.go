package powerplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetPowerApps(environmentName string) ([]App, error) {

	if environmentName == "" {
		environmentName = "NULL"
	}

	request, error := http.NewRequest("GET", fmt.Sprintf("%s/api/environments/%s/apps", client.HostURL, environmentName), nil)
	if error != nil {
		return nil, error
	}
	body, error := client.doRequest(request)
	if error != nil {
		return nil, error
	}

	apps := make([]App, 0)
	error = json.NewDecoder(bytes.NewReader(body)).Decode(&apps)
	if error != nil {
		return nil, error
	}

	return apps, nil
}
