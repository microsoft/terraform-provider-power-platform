// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)

func NewSolutionsDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "solutions",
		},
	}
}

type DataSource struct {
	helpers.TypeInfo
	SolutionClient Client
}

type ListDataSourceModel struct {
	Timeouts      timeouts.Value    `tfsdk:"timeouts"`
	Id            types.String      `tfsdk:"id"`
	EnvironmentId types.String      `tfsdk:"environment_id"`
	Solutions     []DataSourceModel `tfsdk:"solutions"`
}

type DataSourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Name          types.String `tfsdk:"name"`
	CreatedTime   types.String `tfsdk:"created_time"`
	Id            types.String `tfsdk:"id"`
	ModifiedTime  types.String `tfsdk:"modified_time"`
	InstallTime   types.String `tfsdk:"install_time"`
	Version       types.String `tfsdk:"version"`
	IsManaged     types.Bool   `tfsdk:"is_managed"`
}

func ConvertFromSolutionDto(solutionDto Dto) DataSourceModel {
	return DataSourceModel{
		EnvironmentId: types.StringValue(solutionDto.EnvironmentId),
		DisplayName:   types.StringValue(solutionDto.DisplayName),
		Name:          types.StringValue(solutionDto.Name),
		CreatedTime:   types.StringValue(solutionDto.CreatedTime),
		Id:            types.StringValue(solutionDto.Id),
		ModifiedTime:  types.StringValue(solutionDto.ModifiedTime),
		InstallTime:   types.StringValue(solutionDto.InstallTime),
		Version:       types.StringValue(solutionDto.Version),
		IsManaged:     types.BoolValue(solutionDto.IsManaged),
	}
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *DataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Solutions in an environment",
		MarkdownDescription: "Fetches the list of Solutions in an environment.  This is the equivalent of the [`pac solution list`](https://learn.microsoft.com/power-platform/developer/cli/reference/solution#pac-solution-list) command in the Power Platform CLI.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.StringAttribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				Description:         "Unique environment id (guid)",
				MarkdownDescription: "Unique environment id (guid)",
				Required:            true,
			},
			"solutions": schema.ListNestedAttribute{
				Description:         "List of Solutions",
				MarkdownDescription: "List of Solutions",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Solution id",
							Description:         "Solution id",
							Computed:            true,
						},
						"environment_id": schema.StringAttribute{
							MarkdownDescription: "Unique environment id (guid)",
							Description:         "Unique environment id (guid)",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name",
							Description:         "Display name",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Description:         "Name",
							Computed:            true,
						},
						"created_time": schema.StringAttribute{
							MarkdownDescription: "Created time",
							Description:         "Created time",
							Computed:            true,
						},
						"modified_time": schema.StringAttribute{
							MarkdownDescription: "Created time",
							Description:         "Created time",
							Computed:            true,
						},
						"install_time": schema.StringAttribute{
							MarkdownDescription: "Created time",
							Description:         "Created time",
							Computed:            true,
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "Created time",
							Description:         "Created time",
							Computed:            true,
						},
						"is_managed": schema.BoolAttribute{
							MarkdownDescription: "Is managed",
							Description:         "Is managed",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.SolutionClient = NewSolutionClient(clientApi)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	var state ListDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE SOLUTIONS START: %s", d.ProviderTypeName))

	dvExits, err := d.SolutionClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
		return
	}

	solutions, err := d.SolutionClient.GetSolutions(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, solution := range solutions {
		solutionModel := ConvertFromSolutionDto(solution)
		state.Solutions = append(state.Solutions, solutionModel)
	}

	state.Id = types.StringValue(strconv.Itoa(len(solutions)))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE SOLUTIONS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
