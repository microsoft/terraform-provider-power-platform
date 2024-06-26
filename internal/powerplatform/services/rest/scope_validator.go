// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ resource.ConfigValidator = &ScopeValidator{}

type ScopeValidator struct {
}

func (v ScopeValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v ScopeValidator) MarkdownDescription(_ context.Context) string {
	return "If a scope is given in any of the operations (create, update, read, destroy), environment_id cannot be defined at the root level"
}

func (v ScopeValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics = v.Validate(ctx, req.Config)
}

func (v ScopeValidator) Validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	var diags diag.Diagnostics
	var plan *DataverseWebApiResourceModel

	config.Get(ctx, &plan)

	isConfigInvalid := false
	if plan.EnvironmentId.ValueStringPointer() != nil {
		if plan.Create != nil && plan.Create.Scope.ValueStringPointer() != nil {
			isConfigInvalid = true
		}
		if plan.Update != nil && plan.Update.Scope.ValueStringPointer() != nil {
			isConfigInvalid = true
		}
		if plan.Read != nil && plan.Read.Scope.ValueStringPointer() != nil {
			isConfigInvalid = true
		}
		if plan.Destroy != nil && plan.Destroy.Scope.ValueStringPointer() != nil {
			isConfigInvalid = true
		}
		if isConfigInvalid {
			p := path.Path{}
			p.AtName("environment_id")
			diags.AddAttributeError(p,
				"invalid attribute combination",
				"environment_id cannot be defined at the root level if a scope is given in any of the operations (create, update, read, destroy)")
		}
	} else {

		addError := func(operationName string, diag *diag.Diagnostics) {
			p := path.Path{}
			p.AtName(operationName)
			diags.AddAttributeError(p,
				"invalid attribute combination",
				fmt.Sprintf("scope in %s operation is required if environment_id is not defined at the root level", operationName))
		}

		if plan.Create != nil && plan.Create.Scope.ValueStringPointer() == nil {
			addError("create", &diags)
		}
		if plan.Update != nil && plan.Update.Scope.ValueStringPointer() == nil {
			addError("update", &diags)
		}
		if plan.Read != nil && plan.Read.Scope.ValueStringPointer() == nil {
			addError("read", &diags)
		}
		if plan.Destroy != nil && plan.Destroy.Scope.ValueStringPointer() == nil {
			addError("destroy", &diags)
		}
	}
	return diags
}
