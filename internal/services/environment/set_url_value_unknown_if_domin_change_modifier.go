// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SetUrlValueUnknownIfDomainChange() planmodifier.String {
	return &setUrlValueUnknownIfDomainChange{}
}

type setUrlValueUnknownIfDomainChange struct {
}

func (d *setUrlValueUnknownIfDomainChange) Description(ctx context.Context) string {
	return "Ensures that the attribute dataverse.url value is set to unknown if domain.domain attribute changes."
}

func (d *setUrlValueUnknownIfDomainChange) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *setUrlValueUnknownIfDomainChange) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	var planDomainAttribute types.String
	diags := req.Plan.GetAttribute(ctx, path.Root("dataverse").AtName("domain"), &planDomainAttribute)
	resp.Diagnostics.Append(diags...)

	var stateDomainAttribute types.String
	diags = req.State.GetAttribute(ctx, path.Root("dataverse").AtName("domain"), &stateDomainAttribute)
	resp.Diagnostics.Append(diags...)

	if planDomainAttribute.ValueString() != stateDomainAttribute.ValueString() {
		resp.PlanValue = types.StringUnknown()
	}
}
