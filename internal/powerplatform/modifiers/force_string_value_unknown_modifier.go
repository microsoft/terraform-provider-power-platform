// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"bytes"
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ForceStringValueUnknownModifier() planmodifier.String {
	return &forceStringValueUnknownModifier{}
}

type forceStringValueUnknownModifier struct {
}

func (d *forceStringValueUnknownModifier) Description(ctx context.Context) string {
	return "Ensures that file attribute and file checksum attribute are kept synchronised."
}

func (d *forceStringValueUnknownModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *forceStringValueUnknownModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {

	r, _ := req.Private.GetKey(ctx, "force_value_unknown")
	if r == nil || !bytes.Equal(r, []byte("true")) {
		return
	}
	resp.PlanValue = types.StringUnknown()
}
