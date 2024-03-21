// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (client *ApiClient) doRequest(token *string, request *http.Request, headers http.Header) (*ApiHttpResponse, error) {
	apiHttpResponse := &ApiHttpResponse{}
	if headers != nil {
		request.Header = headers
	}

	if token == nil || *token == "" {
		return nil, errors.New("token is empty")
	}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	httpClient := http.DefaultClient

	if request.Header["Authorization"] == nil {
		request.Header.Set("Authorization", "Bearer "+*token)
	}

	if !client.GetConfig().Credentials.TelemetryOptout {
		request.Header.Set("User-Agent", "terraform-provider-power-platform")
	}

	response, err := httpClient.Do(request)
	apiHttpResponse.Response = response
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	apiHttpResponse.BodyAsBytes = body
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if len(body) != 0 {
			return apiHttpResponse, fmt.Errorf("status: %d, message: %s", response.StatusCode, string(body))
		} else {
			return nil, fmt.Errorf("status: %d", response.StatusCode)
		}
	}
	return apiHttpResponse, nil
}

type ApiHttpResponse struct {
	Response    *http.Response
	BodyAsBytes []byte
}

func (apiResponse *ApiHttpResponse) MarshallTo(obj interface{}) error {
	err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
	if err != nil {
		return err
	}
	return nil
}

func (apiResponse *ApiHttpResponse) GetHeader(name string) string {
	return apiResponse.Response.Header.Get(name)
}

func (ApiHttpResponse *ApiHttpResponse) ValidateStatusCode(expectedStatusCode int) error {
	if ApiHttpResponse.Response.StatusCode != expectedStatusCode {
		return fmt.Errorf("expected status code: %d, recieved: %d", expectedStatusCode, ApiHttpResponse.Response.StatusCode)
	}
	return nil
}
