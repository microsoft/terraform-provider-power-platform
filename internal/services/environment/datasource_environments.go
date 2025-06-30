// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"
	"errors"
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

	policyAttributeSchema := map[string]schema.Attribute{
		"type": schema.StringAttribute{
			MarkdownDescription: "Type of the policy according to [schema definition](https://learn.microsoft.com/en-us/azure/templates/microsoft.powerplatform/enterprisepolicies?pivots=deployment-language-terraform#enterprisepolicies-2)",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "Id (guid)",
			Computed:            true,
		},
		"location": schema.StringAttribute{
			MarkdownDescription: "Location of the policy",
			Computed:            true,
		},
		"system_id": schema.StringAttribute{
			MarkdownDescription: "System id (guid)",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "Link status of the policy",
			Computed:            true,
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of environments in a tenant.  See [Environments overview](https://learn.microsoft.com/power-platform/admin/environments-overview) for more information.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"environments": schema.ListNestedAttribute{
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
							MarkdownDescription: "Environment id (guid)",
							Computed:            true,
						},
						"location": schema.StringAttribute{
							MarkdownDescription: "Location of the environment (europe, unitedstates etc.). Can be queried using the `powerplatform_locations` data source.",
							Computed:            true,
						},
						"azure_region": schema.StringAttribute{
							MarkdownDescription: "Azure region of the environment (westeurope, eastus etc.). Can be queried using the `powerplatform_locations` data source.",
							Computed:            true,
						},
						"environment_type": schema.StringAttribute{
							MarkdownDescription: "Type of the environment (Sandbox, Production etc.)",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name",
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
						"release_cycle": schema.StringAttribute{
							MarkdownDescription: "Gives you the ability to create environments that are updated first. This allows you to experience and validate scenarios that are important to you before any updates reach your business-critical applications. See [more](https://learn.microsoft.com/en-us/power-platform/admin/early-release).",
							Computed:            true,
						},
						"billing_policy_id": &schema.StringAttribute{
							MarkdownDescription: "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
							Computed:            true,
						},
						"environment_group_id": schema.StringAttribute{
							MarkdownDescription: "Environment group id (guid) that the environment belongs to. Empty guid `00000000-0000-0000-0000-000000000000` is considered as no environment group.",
							Computed:            true,
						},
						"enterprise_policies": schema.SetNestedAttribute{
							MarkdownDescription: "Enterprise policies for the environment. See [Enterprise policies](https://learn.microsoft.com/en-us/power-platform/admin/enterprise-policies) for more details.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: policyAttributeSchema,
							},
						},
						"owner_id": schema.StringAttribute{
							MarkdownDescription: "Entra ID  user id (guid) of the environment owner when creating developer environment",
							Computed:            true,
						},
						"allow_bing_search": schema.BoolAttribute{
							MarkdownDescription: "Allow Bing search in the environment",
							Computed:            true,
						},
						"allow_moving_data_across_regions": schema.BoolAttribute{
							MarkdownDescription: "Allow moving data across regions",
							Computed:            true,
						},
						"dataverse": schema.SingleNestedAttribute{
							MarkdownDescription: "Dataverse environment details",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"unique_name": schema.StringAttribute{
									MarkdownDescription: "Unique name of the Dataverse environment",
									Computed:            true,
								},
								"administration_mode_enabled": schema.BoolAttribute{
									MarkdownDescription: "Select to enable administration mode for the environment. See [Admin mode](https://learn.microsoft.com/en-us/power-platform/admin/admin-mode) for more information.",
									Computed:            true,
								},
								"background_operation_enabled": schema.BoolAttribute{
									MarkdownDescription: "Background operation status for the environment. See [Admin mode](https://learn.microsoft.com/en-us/power-platform/admin/admin-mode) for more information.",
									Computed:            true,
								},
								"url": schema.StringAttribute{
									MarkdownDescription: "Url of the environment",
									Computed:            true,
								},
								"domain": schema.StringAttribute{
									MarkdownDescription: "Domain name of the environment",
									Computed:            true,
								},
								"organization_id": schema.StringAttribute{
									MarkdownDescription: "Organization id (guid)",
									Computed:            true,
								},
								"security_group_id": schema.StringAttribute{
									MarkdownDescription: "Security group id (guid)",
									Computed:            true,
								},
								"language_code": schema.Int64Attribute{
									MarkdownDescription: "Language LCID (integer)",
									Computed:            true,
								},
								"version": schema.StringAttribute{
									MarkdownDescription: "Version of the environment",
									Computed:            true,
								},
								"linked_app_type": schema.StringAttribute{
									MarkdownDescription: "Type of the linked app (Internal, External etc.)",
									Computed:            true,
								},
								"linked_app_id": schema.StringAttribute{
									MarkdownDescription: "Linked app id (guid)",
									Computed:            true,
								},
								"linked_app_url": schema.StringAttribute{
									MarkdownDescription: "URL of the linked D365 app",
									Computed:            true,
								},
								"currency_code": &schema.StringAttribute{
									MarkdownDescription: "Currency name (EUR, USE, GBP etc.)",
									Computed:            true,
								},
								"templates": schema.ListAttribute{
									MarkdownDescription: "The selected instance provisioning template (if any). See [ERP-based template](https://learn.microsoft.com/en-us/power-platform/admin/unified-experience/tutorial-deploy-new-environment-with-erp-template?tabs=PPAC) for more information.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"template_metadata": schema.StringAttribute{
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

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.EnvironmentClient = NewEnvironmentClient(client.Api)
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
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}

	for _, env := range envs {
		currencyCode := ""
		defaultCurrency, err := d.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, env.Name)
		if err != nil {
			if !errors.Is(err, customerrors.ErrEnvironmentUrlNotFound) {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Unexpected error when reading default currency for environment %s", env.Name),
					err.Error(),
				)
				return
			}
			// Non-critical error (environment URL not found), just skip currency.
		} else {
			currencyCode = defaultCurrency.IsoCurrencyCode
		}

		env, err := convertSourceModelFromEnvironmentDto(env, &currencyCode, nil, nil, nil, timeouts.Value{}, *d.EnvironmentClient.Api.Config)
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
