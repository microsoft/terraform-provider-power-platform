// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func NewWebApiClient(api *api.ApiClient) WebApiClient {
	return WebApiClient{
		Api: api,
	}
}

type WebApiClient struct {
	Api *api.ApiClient
}

type EnvironmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties EnvironmentIdPropertiesDto `json:"properties"`
}

type EnvironmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata LinkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type LinkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}

func (client *WebApiClient) SendOperation(ctx context.Context, operation *DataverseWebApiOperation) (types.Object, error) {

	url := operation.Url.ValueString()
	method := operation.Method.ValueString()
	var body *string = nil
	var headers map[string]string = nil
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

func (client *WebApiClient) ExecuteApiRequest(ctx context.Context, scope *string, url, method string, body *string, headers map[string]string, expectedStatusCodes []int64) (*api.ApiHttpResponse, error) {

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
	} else {
		panic("scope or evironment_id must be provided")
	}
}
