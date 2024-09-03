// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package rest

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewDataverseWebApiDatasource() datasource.DataSource {
	return &DataverseWebApiDatasource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_rest_query",
	}
}

type DataverseWebApiDatasource struct {
	DataRecordClient WebApiClient
	ProviderTypeName string
	TypeName         string
}

type DataverseWebApiDatasourceModel struct {
	Scope              types.String                             `tfsdk:"scope"`
	Method             types.String                             `tfsdk:"method"`
	Url                types.String                             `tfsdk:"url"`
	Body               types.String                             `tfsdk:"body"`
	ExpectedHttpStatus []int64                                  `tfsdk:"expected_http_status"`
	Headers            []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
	Output             types.Object                             `tfsdk:"output"`
}

func (d *DataverseWebApiDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *DataverseWebApiDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Datasource to fetch api requests",
		Attributes: map[string]schema.Attribute{
			"scope": schema.StringAttribute{
				MarkdownDescription: "Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)",
				Required:            true,
				Optional:            false,
			},
			"method": schema.StringAttribute{
				MarkdownDescription: "HTTP method",
				Required:            true,
				Optional:            false,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Absolute url of the api call",
				Required:            true,
				Optional:            false,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "Body of the request",
				Required:            false,
				Optional:            true,
			},
			"expected_http_status": schema.ListAttribute{
				ElementType:         types.Int64Type,
				MarkdownDescription: "Expected HTTP status code. If the response status code does not match any of the expected status codes, the operation will fail.",
				Required:            false,
				Optional:            true,
			},
			"headers": schema.ListNestedAttribute{
				MarkdownDescription: "Headers of the request",
				Required:            false,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Header name",
							Required:            true,
							Optional:            false,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Header value",
							Required:            true,
							Optional:            false,
						},
					},
				},
			},
			"output": schema.SingleNestedAttribute{
				MarkdownDescription: "Response body after executing the web api request",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"body": schema.StringAttribute{
						MarkdownDescription: "Response body after executing the web api request",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *DataverseWebApiDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientApi := req.ProviderData.(*api.ProviderClient).Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.DataRecordClient = NewWebApiClient(clientApi)
}

func (d *DataverseWebApiDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DataverseWebApiDatasourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE START: %s", d.ProviderTypeName))

	resp.State.Get(ctx, &state)

	outputObjectType, err := d.DataRecordClient.SendOperation(ctx, &DataverseWebApiOperation{
		Scope:              state.Scope,
		Method:             state.Method,
		Url:                state.Url,
		Body:               state.Body,
		Headers:            state.Headers,
		ExpectedHttpStatus: state.ExpectedHttpStatus,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to execute request", err.Error())
		return
	}

	state.Output = outputObjectType

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
