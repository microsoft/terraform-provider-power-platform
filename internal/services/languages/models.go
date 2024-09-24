// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type DataSource struct {
	helpers.TypeInfo
	LanguagesClient client
}

type DataSourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	Location types.String   `tfsdk:"location"`
	Value    []DataModel    `tfsdk:"languages"`
}

type DataModel struct {
	Name            string `tfsdk:"name"`
	ID              string `tfsdk:"id"`
	DisplayName     string `tfsdk:"display_name"`
	LocalizedName   string `tfsdk:"localized_name"`
	LocaleID        int64  `tfsdk:"locale_id"`
	IsTenantDefault bool   `tfsdk:"is_tenant_default"`
}
