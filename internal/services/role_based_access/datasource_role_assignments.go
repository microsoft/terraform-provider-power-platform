// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access //nolint:revive // the underscored package name predates this file and matches every service in the repo

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/validators"
)

var (
	_ datasource.DataSource              = &roleAssignmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &roleAssignmentsDataSource{}
)

// roleAssignmentsAttribute returns the shared schema of the list of role assignments returned by the RBAC API.
func roleAssignmentsAttribute(markdownDescription string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		MarkdownDescription: markdownDescription,
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "The unique identifier of the role assignment",
					Computed:            true,
				},
				"scope": schema.StringAttribute{
					MarkdownDescription: "The scope the role assignment applies to",
					Computed:            true,
				},
				"principal_type": schema.StringAttribute{
					MarkdownDescription: "The type of principal the role is assigned to (e.g., `ApplicationUser`, `User`)",
					Computed:            true,
				},
				"principal_id": schema.StringAttribute{
					MarkdownDescription: "The object ID of the enterprise application (service principal) or user the role is assigned to. For `ApplicationUser` principals this is the enterprise application object ID, not the application (client) ID",
					Computed:            true,
				},
				"role_definition_id": schema.StringAttribute{
					MarkdownDescription: "The ID of the assigned role definition",
					Computed:            true,
				},
				"created_by_principal_type": schema.StringAttribute{
					MarkdownDescription: "The type of principal that created the role assignment",
					Computed:            true,
				},
				"created_by_principal_object_id": schema.StringAttribute{
					MarkdownDescription: "The object ID of the principal that created the role assignment",
					Computed:            true,
				},
				"created_on": schema.StringAttribute{
					MarkdownDescription: "The timestamp when the role assignment was created",
					Computed:            true,
				},
				"expires_on": schema.StringAttribute{
					MarkdownDescription: "The timestamp when the role assignment expires, or `null` if it does not expire",
					Computed:            true,
				},
			},
		},
	}
}

func NewRoleAssignmentsDataSource() datasource.DataSource {
	return &roleAssignmentsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "role_assignments",
		},
	}
}

func (d *roleAssignmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *roleAssignmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the [role assignments](https://learn.microsoft.com/en-us/rest/api/power-platform/authorization/role-based-access-control/list-role-assignments) in Power Platform.\n\n" +
			"~> The role based access control API is in [preview](https://learn.microsoft.com/en-us/power-platform/admin/security/role-based-access-control) and Microsoft does not recommend it for production use yet. Reading assignments requires the caller to hold the Power Platform Administrator Entra role or the Power Platform Role Based Access Control Administrator role.\n\n" +
			"Use this data source to discover which principals are assigned roles.\n\n" +
			"Set `scope_type` to choose which assignments to read.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"scope_type": schema.StringAttribute{
				MarkdownDescription: "Which assignments to read. One of `tenant`, `environment` or `environment_group`. " +
					"`environment` requires `environment_id` and `environment_group` requires `environment_group_id`",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(scopeKinds...),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the environment to read assignments from. Required when `scope_type` is `environment`",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("environment_group_id")),
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(helpers.GuidRegex), "environment_id must be a guid"),
					validators.OtherFieldRequiredWhenValueOf(
						path.Root("scope_type").Expression(),
						regexp.MustCompile("^"+scopeEnvironment+"$"), nil,
						"environment_id is required when scope_type is `environment`"),
				},
			},
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the environment group to read assignments from. Required when `scope_type` is `environment_group`",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("environment_id")),
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(helpers.GuidRegex), "environment_group_id must be a guid"),
					validators.OtherFieldRequiredWhenValueOf(
						path.Root("scope_type").Expression(),
						regexp.MustCompile("^"+scopeEnvironmentGroup+"$"), nil,
						"environment_group_id is required when scope_type is `environment_group`"),
				},
			},
			"role_assignments": roleAssignmentsAttribute("List of role assignments at the requested scope"),
		},
	}
}

func (d *roleAssignmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T.", req.ProviderData),
		)
		return
	}
	d.Client = newRoleBasedAccessClient(providerClient.Api)
}

func (d *roleAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state roleAssignmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The shared request-context helper only installs timeout contexts for resource requests, so
	// the advertised read timeout is honoured here directly.
	readTimeout, diags := state.Timeouts.Read(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	scope := state.assignmentScope()
	tflog.Debug(ctx, fmt.Sprintf("Reading role assignments at %s scope", scope))

	assignments, err := d.Client.ListRoleAssignments(ctx, scope)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to list role assignments at %s scope", scope), err.Error())
		return
	}

	state.RoleAssignments = make([]roleAssignmentDataSourceModel, 0, len(assignments))
	for _, assignment := range assignments {
		state.RoleAssignments = append(state.RoleAssignments, convertRoleAssignmentDtoToDataSourceModel(assignment))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
