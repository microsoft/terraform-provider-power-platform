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
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

const (
	EMPTY_BODY   = 0
	RESPONSE_200 = 200
	RESPONSE_300 = 300
)

//nolint:unused-receiver
func (client *Client) BuildCorrelationHeaders(ctx context.Context) (string, string) {
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

	return constants.EMPTY, ""
}

func (client *Client) buildUserAgent(ctx context.Context) string {
	userAgent := fmt.Sprintf("terraform-provider-power-platform/%s (%s; %s) terraform/%s go/%s", common.ProviderVersion, runtime.GOOS, runtime.GOARCH, client.Config.TerraformVersion, runtime.Version())

	requestContext, ok := ctx.Value(helpers.REQUEST_CONTEXT_KEY).(helpers.RequestContextValue)
	if ok {
		userAgent += fmt.Sprintf(" %s %s %s", requestContext.ObjectType, requestContext.ObjectName, requestContext.RequestType)
	}

	return userAgent
}

func (client *Client) doRequest(ctx context.Context, token *string, request *http.Request, headers http.Header) (*HttpResponse, error) {
	apiHttpResponse := &HttpResponse{}
	if headers != nil {
		request.Header = headers
	}

	if token == nil || *token == constants.EMPTY {
		return nil, errors.New("token is empty")
	}

	if request.Header.Get("Content-Type") == constants.EMPTY {
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

	if response.StatusCode < RESPONSE_200 || response.StatusCode >= RESPONSE_300 {
		if len(body) != EMPTY_BODY {
			return apiHttpResponse, fmt.Errorf("status: %d, message: %s", response.StatusCode, string(body))
		}
		return apiHttpResponse, fmt.Errorf("status: %d", response.StatusCode)
	}
	return apiHttpResponse, nil
}

type HttpResponse struct {
	Response    *http.Response
	BodyAsBytes []byte
}

func (apiResponse *HttpResponse) MarshallTo(obj any) error {
	err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
	if err != nil {
		return err
	}
	return nil
}

func (apiResponse *HttpResponse) GetHeader(name string) string {
	return apiResponse.Response.Header.Get(name)
}

func (apiResponse *HttpResponse) ValidateStatusCode(expectedStatusCode int) error {
	if apiResponse.Response.StatusCode != expectedStatusCode {
		return fmt.Errorf("expected status code: %d, recieved: %d", expectedStatusCode, apiResponse.Response.StatusCode)
	}
	return nil
}
