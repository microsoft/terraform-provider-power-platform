// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package rest

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type DataverseWebApiDatasource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type DataverseWebApiDatasourceModel struct {
	Timeouts           timeouts.Value                           `tfsdk:"timeouts"`
	Scope              types.String                             `tfsdk:"scope"`
	Method             types.String                             `tfsdk:"method"`
	Url                types.String                             `tfsdk:"url"`
	Body               types.String                             `tfsdk:"body"`
	ExpectedHttpStatus []int                                    `tfsdk:"expected_http_status"`
	Headers            []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
	Output             types.Object                             `tfsdk:"output"`
}
