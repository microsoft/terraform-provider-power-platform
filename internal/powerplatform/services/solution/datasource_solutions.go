package powerplatform

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	clients "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
)

var (
	_ datasource.DataSource              = &SolutionsDataSource{}
	_ datasource.DataSourceWithConfigure = &SolutionsDataSource{}
)

func NewSolutionsDataSource() datasource.DataSource {
	return &SolutionsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_solutions",
	}
}

type SolutionsDataSource struct {
	SolutionClient   SolutionClient
	ProviderTypeName string
	TypeName         string
}

type SolutionListDataSourceModel struct {
	Id            types.String               `tfsdk:"id"`
	EnvironmentId types.String               `tfsdk:"environment_id"`
	Solutions     []SolutionsDataSourceModel `tfsdk:"solutions"`
}

type SolutionsDataSourceModel struct {
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

func ConvertFromSolutionDto(solutionDto SolutionDto) SolutionsDataSourceModel {
	return SolutionsDataSourceModel{
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

func (d *SolutionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *SolutionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Solutions in an environment",
		MarkdownDescription: "Fetches the list of Solutions in an environment",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
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

func (d *SolutionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBapi := req.ProviderData.(*clients.ProviderClient).BapiApi.Client
	clientDv := req.ProviderData.(*clients.ProviderClient).DataverseApi.Client

	if clientBapi == nil || clientDv == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.SolutionClient = NewSolutionClient(clientBapi, clientDv)
}

func (d *SolutionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SolutionListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE SOLUTIONS START: %s", d.ProviderTypeName))

	solutions, err := d.SolutionClient.GetSolutions(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, solution := range solutions {
		solutionModel := ConvertFromSolutionDto(solution)
		state.Solutions = append(state.Solutions, solutionModel)
	}

	state.Id = types.StringValue(fmt.Sprint((time.Now().Unix())))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE SOLUTIONS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
