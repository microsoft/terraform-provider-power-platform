package powerplatform

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &LocationsDataSource{}
	_ datasource.DataSourceWithConfigure = &LocationsDataSource{}
)

type LocationsDataSourceModel struct {
	Id    types.Int64         `tfsdk:"id"`
	Value []LocationDataModel `tfsdk:"locations"`
}

type LocationDataModel struct {
	ID                                     string   `tfsdk:"id"`
	Name                                   string   `tfsdk:"name"`
	DisplayName                            string   `tfsdk:"display_name"`
	Code                                   string   `tfsdk:"code"`
	IsDefault                              bool     `tfsdk:"is_default"`
	IsDisabled                             bool     `tfsdk:"is_disabled"`
	CanProvisionDatabase                   bool     `tfsdk:"can_provision_database"`
	CanProvisionCustomerEngagementDatabase bool     `tfsdk:"can_provision_customer_engagement_database"`
	AzureRegions                           []string `tfsdk:"azure_regions"`
}

func NewLocationsDataSource() datasource.DataSource {
	return &LocationsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_locations",
	}
}

type LocationsDataSource struct {
	LocationsClient  LocationsClient
	ProviderTypeName string
	TypeName         string
}

func (d *LocationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *LocationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of available Dynamics 365 locations",
		MarkdownDescription: "Fetches the list of available Dynamics 365 locations",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Id of the read operation",
				Optional:    true,
			},
			"locations": schema.ListNestedAttribute{
				Description:         "List of available locations",
				MarkdownDescription: "List of available locations",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the location",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the location",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Display name of the location",
							Computed:    true,
						},
						"code": schema.StringAttribute{
							Description: "Code of the location",
							Computed:    true,
						},
						"is_default": schema.BoolAttribute{
							Description: "Is the location default",
							Computed:    true,
						},
						"is_disabled": schema.BoolAttribute{
							Description: "Is the location disabled",
							Computed:    true,
						},
						"can_provision_database": schema.BoolAttribute{
							Description: "Can the location provision a database",
							Computed:    true,
						},
						"can_provision_customer_engagement_database": schema.BoolAttribute{
							Description: "Can the location provision a customer engagement database",
							Computed:    true,
						},
						"azure_regions": schema.ListAttribute{
							Description:         "List of Azure regions",
							MarkdownDescription: "List of Azure regions",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *LocationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.LocationsClient = NewLocationsClient(clientApi)
}

func (d *LocationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan LocationsDataSourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE LOCATIONS START: %s", d.ProviderTypeName))

	plan.Id = types.Int64Value(time.Now().Unix())

	locations, err := d.LocationsClient.GetLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, location := range locations.Value {
		plan.Value = append(plan.Value, LocationDataModel{
			ID:                                     location.ID,
			Name:                                   location.Name,
			DisplayName:                            location.Properties.DisplayName,
			Code:                                   location.Properties.Code,
			IsDefault:                              location.Properties.IsDefault,
			IsDisabled:                             location.Properties.IsDisabled,
			CanProvisionDatabase:                   location.Properties.CanProvisionDatabase,
			CanProvisionCustomerEngagementDatabase: location.Properties.CanProvisionCustomerEngagementDatabase,
			AzureRegions:                           location.Properties.AzureRegions,
		})
	}

	diags := resp.State.Set(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE LOCATIONS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
