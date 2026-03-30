// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ datasource.DataSource = &DataSource{}
var _ datasource.DataSourceWithConfigure = &DataSource{}

func NewPublisherDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "publisher",
		},
	}
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a Dataverse publisher by publisher id or unique name.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the publisher in provider format `<environment_id>_<publisher_id>`.",
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the Dataverse-enabled environment containing the publisher.",
				Required:            true,
			},
			"publisher_id": schema.StringAttribute{
				MarkdownDescription: "Dataverse publisher id.",
				Optional:            true,
				Computed:            true,
			},
			"uniquename": schema.StringAttribute{
				MarkdownDescription: "Unique name of the publisher.",
				Optional:            true,
				Computed:            true,
			},
			"friendly_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the publisher.",
				Computed:            true,
			},
			"customization_prefix": schema.StringAttribute{
				MarkdownDescription: "Customization prefix used for solution components created by this publisher.",
				Computed:            true,
			},
			"customization_option_value_prefix": schema.Int64Attribute{
				MarkdownDescription: "Option value prefix used for option set values created by this publisher.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the publisher.",
				Computed:            true,
			},
			"email_address": schema.StringAttribute{
				MarkdownDescription: "Email address for the publisher.",
				Computed:            true,
			},
			"supporting_website_url": schema.StringAttribute{
				MarkdownDescription: "Supporting website URL for the publisher.",
				Computed:            true,
			},
			"is_read_only": schema.BoolAttribute{
				MarkdownDescription: "Whether Dataverse reports this publisher as read only.",
				Computed:            true,
			},
			"address": schema.ListNestedAttribute{
				MarkdownDescription: "Publisher addresses mapped from Dataverse address slots 1 and 2.",
				Computed:            true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(2),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"slot": schema.Int64Attribute{
							MarkdownDescription: "Address slot number in Dataverse. Valid values are `1` and `2`.",
							Computed:            true,
							Validators: []validator.Int64{
								int64validator.OneOf(1, 2),
							},
						},
						"address_id": schema.StringAttribute{
							MarkdownDescription: "Dataverse address identifier for the slot.",
							Computed:            true,
						},
						"address_type_code": schema.Int64Attribute{
							MarkdownDescription: "Address type code.",
							Computed:            true,
						},
						"city": schema.StringAttribute{
							MarkdownDescription: "City name.",
							Computed:            true,
						},
						"country": schema.StringAttribute{
							MarkdownDescription: "Country or region name.",
							Computed:            true,
						},
						"county": schema.StringAttribute{
							MarkdownDescription: "County name.",
							Computed:            true,
						},
						"fax": schema.StringAttribute{
							MarkdownDescription: "Fax number.",
							Computed:            true,
						},
						"latitude": schema.Float64Attribute{
							MarkdownDescription: "Latitude value.",
							Computed:            true,
						},
						"line1": schema.StringAttribute{
							MarkdownDescription: "Street line 1.",
							Computed:            true,
						},
						"line2": schema.StringAttribute{
							MarkdownDescription: "Street line 2.",
							Computed:            true,
						},
						"line3": schema.StringAttribute{
							MarkdownDescription: "Street line 3.",
							Computed:            true,
						},
						"longitude": schema.Float64Attribute{
							MarkdownDescription: "Longitude value.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Address name.",
							Computed:            true,
						},
						"postal_code": schema.StringAttribute{
							MarkdownDescription: "Postal code.",
							Computed:            true,
						},
						"post_office_box": schema.StringAttribute{
							MarkdownDescription: "Post office box.",
							Computed:            true,
						},
						"shipping_method_code": schema.Int64Attribute{
							MarkdownDescription: "Shipping method code.",
							Computed:            true,
						},
						"state_or_province": schema.StringAttribute{
							MarkdownDescription: "State or province name.",
							Computed:            true,
						},
						"telephone1": schema.StringAttribute{
							MarkdownDescription: "Primary telephone number.",
							Computed:            true,
						},
						"telephone2": schema.StringAttribute{
							MarkdownDescription: "Secondary telephone number.",
							Computed:            true,
						},
						"telephone3": schema.StringAttribute{
							MarkdownDescription: "Tertiary telephone number.",
							Computed:            true,
						},
						"ups_zone": schema.StringAttribute{
							MarkdownDescription: "UPS zone value.",
							Computed:            true,
						},
						"utc_offset": schema.Int64Attribute{
							MarkdownDescription: "UTC offset for the address.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("publisher_id"),
			path.MatchRoot("uniquename"),
		),
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.PublisherClient = newPublisherClient(providerClient.Api)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var config DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var publisher *publisherDto
	var err error
	if !config.PublisherId.IsNull() && !config.PublisherId.IsUnknown() && config.PublisherId.ValueString() != "" {
		publisher, err = d.PublisherClient.GetPublisherById(ctx, config.EnvironmentId.ValueString(), config.PublisherId.ValueString())
	} else {
		publisher, err = d.PublisherClient.GetPublisherByUniqueName(ctx, config.EnvironmentId.ValueString(), config.UniqueName.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}

	setDataSourceModelFromDto(&config, config.EnvironmentId.ValueString(), publisher)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func setDataSourceModelFromDto(model *DataSourceModel, environmentId string, publisher *publisherDto) {
	model.Id = types.StringValue(buildPublisherResourceId(environmentId, publisher.Id))
	model.EnvironmentId = types.StringValue(environmentId)
	model.PublisherId = types.StringValue(publisher.Id)
	model.UniqueName = types.StringValue(publisher.UniqueName)
	model.FriendlyName = types.StringValue(publisher.FriendlyName)
	model.CustomizationPrefix = types.StringValue(publisher.CustomizationPrefix)
	model.CustomizationOptionValuePrefix = types.Int64Value(publisher.CustomizationOptionValuePrefix)
	model.Description = nullableStringValue(publisher.Description)
	model.EmailAddress = nullableStringValue(publisher.EmailAddress)
	model.SupportingWebsiteURL = nullableStringValue(publisher.SupportingWebsiteURL)
	model.IsReadOnly = types.BoolValue(publisher.IsReadOnly)
	model.Address = addressModelsFromDto(publisher)
}
