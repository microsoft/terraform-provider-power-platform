// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/common"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func (client *ApiClient) BuildCorrelationHeaders(ctx context.Context) (string, string) {
	requestContext, ok := ctx.Value(helpers.REQUEST_CONTEXT_KEY).(helpers.RequestContextValue)
	if ok {
		cc := strings.Join([]string{
			"objectType=" + requestContext.ObjectType,
			"objectName=" + requestContext.ObjectName,
			"requestType=" + requestContext.RequestType,
		}, ",")

		rid := "|" + requestContext.RequestId + "." + fmt.Sprintf("%016x", rand.Uint64()) + "."

		return rid, cc
	}

	return "", ""
}

func (client *ApiClient) buildUserAgent(ctx context.Context) string {
	userAgent := fmt.Sprintf("terraform-provider-power-platform/%s (%s; %s) terraform/%s go/%s", common.ProviderVersion, runtime.GOOS, runtime.GOARCH, client.Config.TerraformVersion, runtime.Version())

	requestContext, ok := ctx.Value(helpers.REQUEST_CONTEXT_KEY).(helpers.RequestContextValue)
	if ok {
		userAgent += fmt.Sprintf(" %s %s %s", requestContext.ObjectType, requestContext.ObjectName, requestContext.RequestType)
	}

	return userAgent
}

func (client *ApiClient) doRequest(ctx context.Context, token *string, request *http.Request, headers http.Header) (*ApiHttpResponse, error) {
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

	if !client.GetConfig().TelemetryOptout {
		ua := client.buildUserAgent(ctx)
		request.Header.Set("User-Agent", ua)

		rid, cc := client.BuildCorrelationHeaders(ctx)
		request.Header.Set("Request-Id", rid)
		request.Header.Set("Correlation-Context", cc)
	}

	response, err := httpClient.Do(request)
	apiHttpResponse.Response = response
	if err != nil {
		return apiHttpResponse, err
	}

	body, err := io.ReadAll(response.Body)
	apiHttpResponse.BodyAsBytes = body
	if err != nil {
		return apiHttpResponse, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if len(body) != 0 {
			return apiHttpResponse, fmt.Errorf("status: %d, message: %s", response.StatusCode, string(body))
		} else {
			return apiHttpResponse, fmt.Errorf("status: %d", response.StatusCode)
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
