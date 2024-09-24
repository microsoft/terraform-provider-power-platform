// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package application

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type TenantApplicationPackagesDataSource struct {
	helpers.TypeInfo
	ApplicationClient client
}

type TenantApplicationPackagesListDataSourceModel struct {
	Timeouts      timeouts.Value                            `tfsdk:"timeouts"`
	Name          types.String                              `tfsdk:"name"`
	PublisherName types.String                              `tfsdk:"publisher_name"`
	Applications  []TenantApplicationPackageDataSourceModel `tfsdk:"applications"`
}

type TenantApplicationPackageDataSourceModel struct {
	ApplicationId          types.String                                   `tfsdk:"application_id"`
	ApplicationDescprition types.String                                   `tfsdk:"application_descprition"`
	Name                   types.String                                   `tfsdk:"application_name"`
	LearnMoreUrl           types.String                                   `tfsdk:"learn_more_url"`
	LocalizedDescription   types.String                                   `tfsdk:"localized_description"`
	LocalizedName          types.String                                   `tfsdk:"localized_name"`
	PublisherId            types.String                                   `tfsdk:"publisher_id"`
	PublisherName          types.String                                   `tfsdk:"publisher_name"`
	UniqueName             types.String                                   `tfsdk:"unique_name"`
	ApplicationVisibility  types.String                                   `tfsdk:"application_visibility"`
	CatalogVisibility      types.String                                   `tfsdk:"catalog_visibility"`
	LastError              []TenantApplicationErrorDetailsDataSourceModel `tfsdk:"last_error"`
}

type TenantApplicationErrorDetailsDataSourceModel struct {
	ErrorCode  types.String `tfsdk:"error_code"`
	ErrorName  types.String `tfsdk:"error_name"`
	Message    types.String `tfsdk:"message"`
	Source     types.String `tfsdk:"source"`
	StatusCode types.Int64  `tfsdk:"status_code"`
	Type       types.String `tfsdk:"type"`
}

type EnvironmentApplicationPackageInstallResource struct {
	helpers.TypeInfo
	ApplicationClient client
}

type EnvironmentApplicationPackageInstallResourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	UniqueName    types.String   `tfsdk:"unique_name"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
}
