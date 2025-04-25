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
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/common"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

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

	if request.Header.Get("Authorization") == "" {
		request.Header.Set("Authorization", "Bearer "+*token)
	}

	if !client.GetConfig().TelemetryOptout {
		ua := client.buildUserAgent(ctx)
		request.Header.Set("User-Agent", ua)
		sessionId, requestId := client.buildCorrelationHeaders(ctx)
		request.Header.Set("X-Correlation-Id", sessionId)
		request.Header.Set("X-Ms-Client-Session-Id", sessionId)
		request.Header.Set("X-Ms-Client-Request-Id", requestId)
	}

	apiResponse, err := httpClient.Do(request)
	resp := &Response{
		HttpResponse: apiResponse,
	}

	if err != nil {
		return resp, err
	}

	if apiResponse == nil {
		return resp, errors.New("unexpected nil response without error")
	}

	defer apiResponse.Body.Close()
	body, err := io.ReadAll(apiResponse.Body)
	resp.BodyAsBytes = body

	// Check for CAE challenge response if CAE is enabled
	if client.Config.EnableContinuousAccessEvaluation && IsCaeChallengeResponse(apiResponse) {
		caeError := &CaePolicyViolationError{
			Message:    "Access denied due to Continuous Access Evaluation (CAE) policy. The authentication token was rejected due to a security policy change.",
			StatusCode: apiResponse.StatusCode,
			Headers:    apiResponse.Header,
		}
		tflog.Warn(ctx, "Detected Continuous Access Evaluation (CAE) challenge response", map[string]any{
			"url":        request.URL.String(),
			"statusCode": apiResponse.StatusCode,
		})
		return resp, caeError
	}

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
		return DefaultRetryAfter()
	}

	return retryAfter
}

func (client *Client) buildCorrelationHeaders(ctx context.Context) (sessionId string, requestId string) {
	sessionId = ""
	requestId = uuid.New().String() // Generate a new request ID for each request
	requestContext, ok := ctx.Value(helpers.REQUEST_CONTEXT_KEY).(helpers.RequestContextValue)
	if ok {
		// If the request context is available, use the session ID from the request context
		sessionId = requestContext.RequestId
	}
	return sessionId, requestId
}

func (client *Client) buildUserAgent(ctx context.Context) string {
	userAgent := fmt.Sprintf("terraform-provider-power-platform/%s (%s; %s) terraform/%s go/%s", common.ProviderVersion, runtime.GOOS, runtime.GOARCH, client.Config.TerraformVersion, runtime.Version())

	requestContext, ok := ctx.Value(helpers.REQUEST_CONTEXT_KEY).(helpers.RequestContextValue)
	if ok {
		userAgent += fmt.Sprintf(" %s %s", requestContext.ObjectName, requestContext.RequestType)
	}

	return userAgent
}
