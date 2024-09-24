// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package currencies

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	Location types.String   `tfsdk:"location"`
	Value    []DataModel    `tfsdk:"currencies"`
}

type DataModel struct {
	ID              string `tfsdk:"id"`
	Name            string `tfsdk:"name"`
	Type            string `tfsdk:"type"`
	Code            string `tfsdk:"code"`
	Symbol          string `tfsdk:"symbol"`
	IsTenantDefault bool   `tfsdk:"is_tenant_default"`
}
