package powerplatform

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

var (
	_ datasource.DataSource              = &ApplicationUserDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationUserDataSource{}
)

func NewApplicationUserDataSource() datasource.DataSource {
	return &ApplicationUserDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_application_users",
	}
}

type ApplicationUserDataSource struct {
	BapiApiClient    powerplatform_bapi.ApiClientInterface
	ProviderTypeName string
	TypeName         string
}

type ApplicationUserListDataSourceModel struct {
	Id              types.String                     `tfsdk:"id"`
	EnvironmentName types.String                     `tfsdk:"environment_id"`
	ApplicationUser []ApplicationUserDataSourceModel `tfsdk:"application_user"`
}

type ApplicationUserDataSourceModel struct {
	EnvironmentName types.String `tfsdk:"environment_id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Name            types.String `tfsdk:"name"`
	CreatedTime     types.String `tfsdk:"created_time"`
	Id              types.String `tfsdk:"id"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
	InstallTime     types.String `tfsdk:"install_time"`
	Version         types.String `tfsdk:"version"`
	IsManaged       types.Bool   `tfsdk:"is_managed"`
}

func ConvertFromApplicationUserDto(applicationUserDto models.ApplicationUserDto) ApplicationUserDataSourceModel {
	return ApplicationUserDataSourceModel{
		EnvironmentName: types.StringValue(applicationUserDto.EnvironmentName),
		DisplayName:     types.StringValue(applicationUserDto.DisplayName),
		Name:            types.StringValue(applicationUserDto.Name),
		CreatedTime:     types.StringValue(applicationUserDto.CreatedTime),
		Id:              types.StringValue(applicationUserDto.Id),
		ModifiedTime:    types.StringValue(applicationUserDto.ModifiedTime),
		InstallTime:     types.StringValue(applicationUserDto.InstallTime),
		Version:         types.StringValue(applicationUserDto.Version),
		IsManaged:       types.BoolValue(applicationUserDto.IsManaged),
	}
}

func (d *ApplicationUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *ApplicationUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of ApplicationUser in an environment",
		MarkdownDescription: "Fetches the list of ApplicationUser in an environment",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"environment_id": schema.StringAttribute{
				Description:         "Unique environment name (guid)",
				MarkdownDescription: "Unique environment name (guid)",
				Required:            true,
			},
			"application_users": schema.ListNestedAttribute{
				Description:         "List of ApplicationUser",
				MarkdownDescription: "List of ApplicationUser",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "ApplicationUser id",
							Description:         "ApplicationUser id",
							Computed:            true,
						},
						"environment_id": schema.StringAttribute{
							MarkdownDescription: "Unique environment name (guid)",
							Description:         "Unique environment name (guid)",
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

func (d *ApplicationUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*PowerPlatformProvider).bapiClient.(powerplatform_bapi.ApiClientInterface)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.BapiApiClient = client
}

func (d *ApplicationUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ApplicationUserListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ApplicationUser START: %s", d.ProviderTypeName))

	ApplicationUser, err := d.BapiApiClient.GetApplicationUser(ctx, state.EnvironmentName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, ApplicationUser := range ApplicationUser {
		ApplicationUserModel := ConvertFromApplicationUserDto(ApplicationUser)
		state.ApplicationUser = append(state.ApplicationUser, ApplicationUserModel)
	}

	state.Id = types.StringValue(fmt.Sprint((time.Now().Unix())))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ApplicationUser END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
