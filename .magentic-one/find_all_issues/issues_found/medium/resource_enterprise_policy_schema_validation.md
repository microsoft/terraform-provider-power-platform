# Issue 2: Insufficient Validation for `Schema` Attributes

##

`/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go`

## Problem

The `Schema` method defines several attributes (`id`, `environment_id`, `system_id`, `policy_type`) with required, computed, or validated properties, but the validation logic is not thorough. For example:
- `environment_id` requires a strict format validation to ensure it adheres to expected standards.
- `system_id` should validate its format to ensure it matches `/regions/<location>/providers/Microsoft.PowerPlatform/enterprisePolicies/<policyid>`.
- `policy_type` focuses on two predefined values but does not account for error handling when an unsupported value is provided (outside of the Validator).

## Impact

Without stricter validations, invalid data may pass through unnoticed, leading to failed operations downstream or inconsistent states during runtime. This can lead to unexpected errors and reduced reliability during infrastructure deployments.

Severity: **Medium**

## Location

```go
"environment_id": schema.StringAttribute{
	MarkdownDescription: "Environment id",
	Required:            true,
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	},
},
"system_id": schema.StringAttribute{
	MarkdownDescription: "Policy SystemId value in following format `/regions/<location>/providers/Microsoft.PowerPlatform/enterprisePolicies/<policyid>`",
	Required:            true,
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	},
},
"policy_type": schema.StringAttribute{
	MarkdownDescription: fmt.Sprintf("Policy type [%s, %s]", NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
	Required:            true,
	Validators: []validator.String{
		stringvalidator.OneOf(NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
	},
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	},
},
```

## Fix

Introduce stricter validation logic for attribute formats and improve error messages when invalid values are provided.

```go
import (
	"errors"
)

resp.Schema = schema.Schema{
	MarkdownDescription: "Enterprise Policy environment assignment",
	Attributes: map[string]schema.Attribute{
		"environment_id": schema.StringAttribute{
			MarkdownDescription: "Environment id",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.RegexMatch(`^env-[a-zA-Z0-9]{8,}$`, "Environment ID must start with 'env-' and be followed by alphanumeric characters"),
			},
		},
		"system_id": schema.StringAttribute{
			MarkdownDescription: "Policy SystemId formatted as `/regions/<location>/providers/Microsoft.PowerPlatform/enterprisePolicies/<policyid>`",
			Required:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.RegexMatch(`^/regions/[a-zA-Z0-9\-_]+/providers/Microsoft.PowerPlatform/enterprisePolicies/[a-zA-Z0-9\-_]+$`, "System ID must match the expected format"),
			},
		},
		"policy_type": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("Policy type [%s, %s]", NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
			},
		},
	},
}
```

The `RegexMatch` validation ensures stricter adherence to specific formats and provides clearer error messages for developers and users.