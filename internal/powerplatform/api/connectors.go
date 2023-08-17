package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetConnectors(ctx context.Context) ([]Connector, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/connectors", client.BaseUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	conn := make([]Connector, 0)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
