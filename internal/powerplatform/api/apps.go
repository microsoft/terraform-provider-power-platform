package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetPowerApps(ctx context.Context, environmentName string) ([]App, error) {

	if environmentName == "" {
		environmentName = "NULL"
	}

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/environments/%s/apps", client.BaseUrl, environmentName), nil)
	if err != nil {
		return nil, err
	}
	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	apps := make([]App, 0)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}
