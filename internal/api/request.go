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
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/common"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func (client *Client) BuildCorrelationHeaders(ctx context.Context) (requestId string, correlationContext string) {
	requestContext, ok := ctx.Value(helpers.REQUEST_CONTEXT_KEY).(helpers.RequestContextValue)
	if ok {
		cc := strings.Join([]string{
			"objectName=" + requestContext.ObjectName,
			"requestType=" + requestContext.RequestType,
		}, ",")

		rid := "|" + requestContext.RequestId + "." + fmt.Sprintf("%016x", rand.Uint64()) + "."

		return rid, cc
	}
	return "", ""
}

func (client *Client) buildUserAgent(ctx context.Context) string {
	userAgent := fmt.Sprintf("terraform-provider-power-platform/%s (%s; %s) terraform/%s go/%s", common.ProviderVersion, runtime.GOOS, runtime.GOARCH, client.Config.TerraformVersion, runtime.Version())

	requestContext, ok := ctx.Value(helpers.REQUEST_CONTEXT_KEY).(helpers.RequestContextValue)
	if ok {
		userAgent += fmt.Sprintf(" %s %s", requestContext.ObjectName, requestContext.RequestType)
	}

	return userAgent
}

func (client *Client) doRequest(ctx context.Context, token *string, request *http.Request, headers http.Header) (*Response, error) {
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

	apiResponse, err := httpClient.Do(request)

	resp := &Response{
		HttpResponse: apiResponse,
	}

	if err != nil {
		return resp, err
	}

	defer apiResponse.Body.Close()
	body, err := io.ReadAll(apiResponse.Body)
	resp.BodyAsBytes = body

	return resp, err
}

type Response struct {
	HttpResponse *http.Response
	BodyAsBytes  []byte
}

func (apiResponse *Response) MarshallTo(obj any) error {
	err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
	if err != nil {
		return err
	}
	return nil
}

func (apiResponse *Response) GetHeader(name string) string {
	return apiResponse.HttpResponse.Header.Get(name)
}

func retryAfter(ctx context.Context, resp *http.Response) time.Duration {
	retryHeader := resp.Header.Get(constants.HEADER_RETRY_AFTER)
	tflog.Debug(ctx, "Retry Header: "+retryHeader)

	retryAfter, err := time.ParseDuration(retryHeader)
	if err != nil {
		// default retry after 5-10 seconds
		return time.Duration((rand.Intn(5) + 5)) * time.Second
	}

	return retryAfter
}
