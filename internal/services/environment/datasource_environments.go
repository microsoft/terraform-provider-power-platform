// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &EnvironmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentsDataSource{}
)

func NewEnvironmentsDataSource() datasource.DataSource {
	return &EnvironmentsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environments",
		},
	}
}

func (d *EnvironmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *EnvironmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of environments in a tenant",
		MarkdownDescription: "Fetches the list of environments in a tenant.  See [Environments overview](https://learn.microsoft.com/power-platform/admin/environments-overview) for more information.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
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
						"description": schema.StringAttribute{
							MarkdownDescription: "Description",
							Computed:            true,
						},
						"cadence": schema.StringAttribute{
							MarkdownDescription: "Cadence of updates for the environment (Frequent, Moderate)",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("Frequent", "Moderate"),
							},
						},
						"billing_policy_id": &schema.StringAttribute{
							Description:         "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
							MarkdownDescription: "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
							Computed:            true,
						},
						"environment_group_id": schema.StringAttribute{
							MarkdownDescription: "Unique environment group id (guid) that the environment belongs to. Empty guid `00000000-0000-0000-0000-000000000000` is considered as no environment group.",
							Computed:            true,
						},
						"dataverse": schema.SingleNestedAttribute{
							MarkdownDescription: "Dataverse environment details",
							Description:         "Dataverse environment details",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"administration_mode_enabled": schema.BoolAttribute{
									MarkdownDescription: "Select to enable administration mode for the environment. See [Admin mode](https://learn.microsoft.com/en-us/power-platform/admin/admin-mode) for more information. ",
									Computed:            true,
								},
								"background_operation_enabled": schema.BoolAttribute{
									MarkdownDescription: "Background operation status for the environment. See [Admin mode](https://learn.microsoft.com/en-us/power-platform/admin/admin-mode) for more information. ",
									Computed:            true,
								},
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
									MarkdownDescription: "The selected instance provisioning template (if any). See [ERP-based template](https://learn.microsoft.com/en-us/power-platform/admin/unified-experience/tutorial-deploy-new-environment-with-erp-template?tabs=PPAC) for more information.",
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

func (d *EnvironmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state ListDataSourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envs, err := d.EnvironmentClient.GetEnvironments(ctx)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, env := range envs {
		currencyCode := ""
		defaultCurrency, err := d.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, env.Name)
		if err != nil {
			if customerrors.Code(err) != customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND {
				resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", env.Name), err.Error())
			}
		} else {
			currencyCode = defaultCurrency.IsoCurrencyCode
		}

		env, err := convertSourceModelFromEnvironmentDto(env, &currencyCode, nil, nil, timeouts.Value{})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error when converting environment %s", env.DisplayName), err.Error())
			return
		}
		state.Environments = append(state.Environments, *env)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
