// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &UnmanagedSolutionDataSource{}
	_ datasource.DataSourceWithConfigure = &UnmanagedSolutionDataSource{}
)

func NewUnmanagedSolutionDataSource() datasource.DataSource {
	return &UnmanagedSolutionDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "unmanaged_solution",
		},
	}
}

func (d *UnmanagedSolutionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *UnmanagedSolutionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a single unmanaged Dataverse solution by unique name.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the unmanaged solution in provider format `<environment_id>_<solution_id>`.",
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the Dataverse-enabled environment containing the unmanaged solution.",
				Required:            true,
			},
			"uniquename": schema.StringAttribute{
				MarkdownDescription: "Unique name of the unmanaged solution.",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the unmanaged solution.",
				Computed:            true,
			},
			"publisher_id": schema.StringAttribute{
				MarkdownDescription: "Existing Dataverse publisher id that owns the solution.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the unmanaged solution.",
				Computed:            true,
			},
		},
	}
}

func (d *UnmanagedSolutionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
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

	d.SolutionClient = NewSolutionClient(client.Api)
}

func (d *UnmanagedSolutionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state UnmanagedSolutionDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dvExists, err := d.SolutionClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
		return
	}

	if !dvExists {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
		return
	}

	solution, err := d.SolutionClient.GetSolutionUniqueName(ctx, state.EnvironmentId.ValueString(), state.UniqueName.ValueString())
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unmanaged solution '%s' not found", state.UniqueName.ValueString()),
				fmt.Sprintf("No unmanaged solution with unique name '%s' was found in environment '%s'.", state.UniqueName.ValueString(), state.EnvironmentId.ValueString()),
			)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}

	setUnmanagedSolutionDataSourceState(&state, solution)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func setUnmanagedSolutionDataSourceState(model *UnmanagedSolutionDataSourceModel, solution *SolutionDto) {
	model.Id = types.StringValue(fmt.Sprintf("%s_%s", model.EnvironmentId.ValueString(), solution.Id))
	model.UniqueName = types.StringValue(solution.Name)
	model.DisplayName = types.StringValue(solution.DisplayName)
	model.PublisherId = types.StringValue(solution.PublisherId)

	if solution.Description == "" {
		model.Description = types.StringNull()
	} else {
		model.Description = types.StringValue(solution.Description)
	}
}
