// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

var (
	_ datasource.DataSource              = &EnvironmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentsDataSource{}
)

func NewEnvironmentsDataSource() datasource.DataSource {
	return &EnvironmentsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environments",
	}
}

type EnvironmentsDataSource struct {
	EnvironmentClient EnvironmentClient
	ProviderTypeName  string
	TypeName          string
}

func (d *EnvironmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *EnvironmentsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of environments in a tenant",
		MarkdownDescription: "Fetches the list of environments in a tenant.  See [Environments overview](https://learn.microsoft.com/power-platform/admin/environments-overview) for more information.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Computed:            true,
			},
			"environments": schema.ListNestedAttribute{
				Description:         "List of environments",
				MarkdownDescription: "List of environments",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
							Create: false,
							Update: false,
							Delete: false,
							Read:   false,
						}),
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique environment id (guid)",
							Description:         "Unique environment id (guid)",
							Computed:            true,
						},
						"location": schema.StringAttribute{
							Description:         "Location of the environment (europe, unitedstates etc.). Can be queried using the `powerplatform_locations` data source.",
							MarkdownDescription: "Location of the environment (europe, unitedstates etc.). Can be queried using the `powerplatform_locations` data source.",
							Computed:            true,
						},
						"azure_region": schema.StringAttribute{
							Description:         "Azure region of the environment (westeurope, eastus etc.). Can be queried using the `powerplatform_locations` data source.",
							MarkdownDescription: "Azure region of the environment (westeurope, eastus etc.). Can be queried using the `powerplatform_locations` data source.",
							Computed:            true,
						},
						"environment_type": schema.StringAttribute{
							Description:         "Type of the environment (Sandbox, Production etc.)",
							MarkdownDescription: "Type of the environment (Sandbox, Production etc.)",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name",
							Description:         "Display name",
							Computed:            true,
						},
						"billing_policy_id": &schema.StringAttribute{
							Description:         "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
							MarkdownDescription: "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
							Computed:            true,
						},
						"dataverse": schema.SingleNestedAttribute{
							MarkdownDescription: "Dataverse environment details",
							Description:         "Dataverse environment details",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Description:         "Url of the environment",
									MarkdownDescription: "Url of the environment",
									Computed:            true,
								},
								"domain": schema.StringAttribute{
									Description:         "Domain name of the environment",
									MarkdownDescription: "Domain name of the environment",
									Computed:            true,
								},
								"organization_id": schema.StringAttribute{
									Description:         "Unique organization id (guid)",
									MarkdownDescription: "Unique organization id (guid)",
									Computed:            true,
								},
								"security_group_id": schema.StringAttribute{
									Description:         "Unique security group id (guid)",
									MarkdownDescription: "Unique security group id (guid)",
									Computed:            true,
								},
								"language_code": schema.Int64Attribute{
									Description:         "Unique language LCID (integer)",
									MarkdownDescription: "Unique language LCID (integer)",
									Computed:            true,
								},
								"version": schema.StringAttribute{
									Description:         "Version of the environment",
									MarkdownDescription: "Version of the environment",
									Computed:            true,
								},
								//Not available in BAPI as for now
								// "currency_name": &schema.StringAttribute{
								// 	Description:         "Unique currency name (EUR, USE, GBP etc.)",
								// 	MarkdownDescription: "Unique currency name (EUR, USE, GBP etc.)",
								// 	Computed:            true,
								// },
								"linked_app_type": schema.StringAttribute{
									Description:         "Type of the linked app (Internal, External etc.)",
									MarkdownDescription: "Type of the linked app (Internal, External etc.)",
									Computed:            true,
								},
								"linked_app_id": schema.StringAttribute{
									Description:         "Unique linked app id (guid)",
									MarkdownDescription: "Unique linked app id (guid)",
									Computed:            true,
								},
								"linked_app_url": schema.StringAttribute{
									Description:         "URL of the linked D365 app",
									MarkdownDescription: "URL of the linked D365 app",
									Computed:            true,
								},
								"currency_code": &schema.StringAttribute{
									Description:         "Unique currency name (EUR, USE, GBP etc.)",
									MarkdownDescription: "Unique currency name (EUR, USE, GBP etc.)",
									Computed:            true,
								},
								"templates": schema.ListAttribute{
									Description:         "The selected instance provisioning template (if any)",
									MarkdownDescription: "The selected instance provisioning template (if any)",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"template_metadata": schema.StringAttribute{
									Description:         "Additional D365 environment template metadata (if any)",
									MarkdownDescription: "Additional D365 environment template metadata (if any)",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *EnvironmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.EnvironmentClient = NewEnvironmentClient(clientApi)
}

func (d *EnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvironmentsListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENTS START: %s", d.ProviderTypeName))

	envs, err := d.EnvironmentClient.GetEnvironments(ctx)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, env := range envs {
		currencyCode := ""
		defaultCurrency, err := d.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, env.Name)
		if err != nil {
			if helpers.Code(err) != helpers.ERROR_ENVIRONMENT_URL_NOT_FOUND {
				resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", env.Name), err.Error())
			}
		} else {
			currencyCode = defaultCurrency.IsoCurrencyCode
		}

		env, err := ConvertSourceModelFromEnvironmentDto(env, &currencyCode, nil, nil, timeouts.Value{})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error when converting environment %s", env.DisplayName), err.Error())
			return
		}
		state.Environments = append(state.Environments, *env)
	}
	state.Id = types.Int64Value(int64(len(envs)))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENTS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
