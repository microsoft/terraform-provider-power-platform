// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func NewWebApiClient(apiClient *api.Client) WebApiClient {
	return WebApiClient{
		Api: apiClient,
	}
}

type WebApiClient struct {
	Api *api.Client
}

func (client *WebApiClient) SendOperation(ctx context.Context, operation *DataverseWebApiOperation) (types.Object, error) {
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

	res, err := client.ExecuteApiRequest(ctx, operation.Scope.ValueStringPointer(), url, method, body, headers, operation.ExpectedHttpStatus)
	if helpers.Code(err) == helpers.ERROR_UNEXPECTED_HTTP_RETURN_CODE {
		return types.ObjectUnknown(map[string]attr.Type{
			"body": types.StringType,
		}), err
	}

	if res != nil && res.Response != nil {
		tflog.Trace(ctx, fmt.Sprintf("SendOperation Response: %v", res.BodyAsBytes))
		tflog.Trace(ctx, fmt.Sprintf("SendOperation Response Status: %v", res.Response.Status))
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

func (client *WebApiClient) ExecuteApiRequest(ctx context.Context, scope *string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int64) (*api.HttpResponse, error) {
	h := http.Header{}
	for k, v := range headers {
		h.Add(k, v)
	}

	codes := make([]int, len(expectedStatusCodes))
	for i, code := range expectedStatusCodes {
		codes[i] = int(code)
	}

	if scope != nil {
		return client.Api.ExecuteForGivenScope(ctx, *scope, method, url, h, body, codes, nil)
	}
	panic("scope or evironment_id must be provided")
}
