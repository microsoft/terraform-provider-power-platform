package powerplatform

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
)

var (
	_ datasource.DataSource              = &ApplicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &ApplicationsDataSource{}
)

func NewApplicationsDataSource() datasource.DataSource {
	return &ApplicationsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_applications",
	}
}

type ApplicationsDataSource struct {
	ApplicationClient ApplicationClient
	ProviderTypeName  string
	TypeName          string
}

type ApplicationsListDataSourceModel struct {
	EnvironmentId types.String                 `tfsdk:"environment_id"`
	Id            types.String                 `tfsdk:"id"`
	Applications  []ApplicationDataSourceModel `tfsdk:"applications"`
}

type ApplicationDataSourceModel struct {
	ApplicationId types.String `tfsdk:"application_id"`
	Name          types.String `tfsdk:"application_name"`
	UniqueName    types.String `tfsdk:"unique_name"`
}

func (d *ApplicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *ApplicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dynamics 365 applications in a tenant",
		MarkdownDescription: "Fetches the list of Dynamics 365 applications in a tenant",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Id of the Dynamics 365 application",
				Optional:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Id of the Dynamics 365 environment",
				Optional:    true,
			},
			"applications": schema.ListNestedAttribute{
				Description:         "List of Connectors",
				MarkdownDescription: "List of Connectors",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"application_id": schema.StringAttribute{
							MarkdownDescription: "ApplicaitonId",
							Description:         "ApplicaitonId",
							Computed:            true,
						},
						"application_name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Description:         "Name",
							Computed:            true,
						},
						"unique_name": schema.StringAttribute{
							MarkdownDescription: "Unique Name",
							Description:         "Unique Name",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ApplicationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientBapi := req.ProviderData.(*clients.ProviderClient).PowerPlatformApi.Client
	if clientBapi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.ApplicationClient = NewApplicationClient(clientBapi)
}

func (d *ApplicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan ApplicationsListDataSourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE APPLICATIONS START: %s", d.ProviderTypeName))

	plan.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))
	plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())

	applications, err := d.ApplicationClient.GetApplicationsByEnvironmentId(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, application := range applications {
		plan.Applications = append(plan.Applications, ApplicationDataSourceModel{
			ApplicationId: types.StringValue(application.ApplicationId),
			Name:          types.StringValue(application.Name),
			UniqueName:    types.StringValue(application.UniqueName),
		})
	}

	diags := resp.State.Set(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE APPLICATIONS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
