// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func SetStringValueToUnknownIfChecksumsChangeModifier(firstAttributePair, secondAttributePair []string) planmodifier.String {
	return &setStringValueToUnknownIfChecksumsChangeModifier{
		firstAttributePair:  firstAttributePair,
		secondAttributePair: secondAttributePair,
	}
}

type setStringValueToUnknownIfChecksumsChangeModifier struct {
	firstAttributePair  []string
	secondAttributePair []string
}

func (d *setStringValueToUnknownIfChecksumsChangeModifier) Description(ctx context.Context) string {
	return "Ensures that the attribute value is set to unknown if the checksums change."
}

func (d *setStringValueToUnknownIfChecksumsChangeModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d *setStringValueToUnknownIfChecksumsChangeModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
	if resp.Diagnostics.HasError() {
		return
	}

	secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
	if resp.Diagnostics.HasError() {
		return
	}

	if firstAttributeHasChanged || secondAttributeHasChanged {
		resp.PlanValue = types.StringUnknown()
	}
}

func (d *setStringValueToUnknownIfChecksumsChangeModifier) hasChecksumChanged(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse, attributeName, checksumAttributeName string) bool {
	var attribute types.String
	diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return false
	}
	
	var attributeChecksum types.String
	diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return false
	}

	value, err := helpers.CalculateSHA256(attribute.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for attribute %q", attributeName), err.Error())
	}

	return value != attributeChecksum.ValueString()
}
