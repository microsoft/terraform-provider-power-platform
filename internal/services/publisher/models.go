// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type Resource struct {
	helpers.TypeInfo
	PublisherClient client
}

type DataSource struct {
	helpers.TypeInfo
	PublisherClient client
}

type PublisherAddressModel struct {
	Slot               types.Int64   `tfsdk:"slot"`
	AddressId          types.String  `tfsdk:"address_id"`
	AddressTypeCode    types.Int64   `tfsdk:"address_type_code"`
	City               types.String  `tfsdk:"city"`
	Country            types.String  `tfsdk:"country"`
	County             types.String  `tfsdk:"county"`
	Fax                types.String  `tfsdk:"fax"`
	Latitude           types.Float64 `tfsdk:"latitude"`
	Line1              types.String  `tfsdk:"line1"`
	Line2              types.String  `tfsdk:"line2"`
	Line3              types.String  `tfsdk:"line3"`
	Longitude          types.Float64 `tfsdk:"longitude"`
	Name               types.String  `tfsdk:"name"`
	PostalCode         types.String  `tfsdk:"postal_code"`
	PostOfficeBox      types.String  `tfsdk:"post_office_box"`
	ShippingMethodCode types.Int64   `tfsdk:"shipping_method_code"`
	StateOrProvince    types.String  `tfsdk:"state_or_province"`
	Telephone1         types.String  `tfsdk:"telephone1"`
	Telephone2         types.String  `tfsdk:"telephone2"`
	Telephone3         types.String  `tfsdk:"telephone3"`
	UpsZone            types.String  `tfsdk:"ups_zone"`
	UtcOffset          types.Int64   `tfsdk:"utc_offset"`
}

type ResourceModel struct {
	Timeouts                       timeouts.Value          `tfsdk:"timeouts"`
	Id                             types.String            `tfsdk:"id"`
	EnvironmentId                  types.String            `tfsdk:"environment_id"`
	PublisherId                    types.String            `tfsdk:"publisher_id"`
	UniqueName                     types.String            `tfsdk:"uniquename"`
	FriendlyName                   types.String            `tfsdk:"friendly_name"`
	CustomizationPrefix            types.String            `tfsdk:"customization_prefix"`
	CustomizationOptionValuePrefix types.Int64             `tfsdk:"customization_option_value_prefix"`
	Description                    types.String            `tfsdk:"description"`
	EmailAddress                   types.String            `tfsdk:"email_address"`
	SupportingWebsiteURL           types.String            `tfsdk:"supporting_website_url"`
	IsReadOnly                     types.Bool              `tfsdk:"is_read_only"`
	Address                        []PublisherAddressModel `tfsdk:"address"`
}

type DataSourceModel struct {
	Timeouts                       timeouts.Value          `tfsdk:"timeouts"`
	Id                             types.String            `tfsdk:"id"`
	EnvironmentId                  types.String            `tfsdk:"environment_id"`
	PublisherId                    types.String            `tfsdk:"publisher_id"`
	UniqueName                     types.String            `tfsdk:"uniquename"`
	FriendlyName                   types.String            `tfsdk:"friendly_name"`
	CustomizationPrefix            types.String            `tfsdk:"customization_prefix"`
	CustomizationOptionValuePrefix types.Int64             `tfsdk:"customization_option_value_prefix"`
	Description                    types.String            `tfsdk:"description"`
	EmailAddress                   types.String            `tfsdk:"email_address"`
	SupportingWebsiteURL           types.String            `tfsdk:"supporting_website_url"`
	IsReadOnly                     types.Bool              `tfsdk:"is_read_only"`
	Address                        []PublisherAddressModel `tfsdk:"address"`
}
