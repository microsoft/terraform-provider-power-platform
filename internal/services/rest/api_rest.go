// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func newWebApiClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) SendOperation(ctx context.Context, operation *DataverseWebApiOperation) (types.Object, error) {
	url := operation.Url.ValueString()
	method := operation.Method.ValueString()
	var body *string
	var headers map[string]string
	if operation.Body.ValueStringPointer() != nil {
		b := operation.Body.ValueString()
		body = &b
	}
	if len(operation.Headers) > 0 {
		headers = make(map[string]string)
		for _, h := range operation.Headers {
			headers[h.Name.ValueString()] = h.Value.ValueString()
		}
	}

	expectedStatusCodes := operation.ExpectedHttpStatus
	if len(operation.ExpectedHttpStatus) == 0 {
		expectedStatusCodes = []int{http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent}
	}

	res, err := client.ExecuteApiRequest(ctx, operation.Scope.ValueStringPointer(), url, method, body, headers, expectedStatusCodes)
	var unexpected customerrors.UnexpectedHttpStatusCodeError
	if errors.As(err, &unexpected) {
		return types.ObjectUnknown(map[string]attr.Type{
			"body": types.StringType,
		}), unexpected
	}

	if res != nil && res.HttpResponse != nil {
		tflog.Trace(ctx, fmt.Sprintf("SendOperation Response: %v", res.BodyAsBytes))
		tflog.Trace(ctx, fmt.Sprintf("SendOperation Response Status: %v", res.HttpResponse.Status))
	}

	output := map[string]attr.Value{
		"body": types.StringNull(),
	}

	if res == nil && err != nil {
		output["body"] = types.StringValue(err.Error())
	} else {
		if len(res.BodyAsBytes) > 0 {
			output["body"] = types.StringValue(string(res.BodyAsBytes))
		}
	}
	o := types.ObjectValueMust(map[string]attr.Type{
		"body": types.StringType,
	}, output)
	return o, nil
}

func (client *client) ExecuteApiRequest(ctx context.Context, scope *string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int) (*api.Response, error) {
	h := http.Header{}
	for k, v := range headers {
		h.Add(k, v)
	}

	if scope != nil {
		return client.Api.Execute(ctx, []string{*scope}, method, url, h, body, expectedStatusCodes, nil)
	}
	panic("scope or evironment_id must be provided")
}
