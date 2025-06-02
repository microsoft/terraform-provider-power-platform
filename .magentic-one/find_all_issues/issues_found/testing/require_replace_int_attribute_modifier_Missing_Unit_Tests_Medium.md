# Title

Missing Unit Tests for Modifier Logic

##

/workspaces/terraform-provider-power-platform/internal/modifiers/require_replace_int_attribute_modifier.go

## Problem

There is no test file alongside this modifier, and no unit tests for its logic. This results in a lack of automated verification for this core behavior: forcing a resource replacement when an int64 attribute is changed and the prior state is non-null, non-unknown, and non-zero.

## Impact

Medium: If logic is changed or refactored (e.g. the conditional or plan modify method), regressions may not be detected. Test coverage is essential, especially for business rules dictating resource recreation.

## Location

No corresponding test file or methods found for the modifier's plan logic.

## Code Issue

_No code for tests is present for this functionality._

## Fix

Implement a test file named `require_replace_int_attribute_modifier_test.go` which contains table-driven tests for the plan modifier method. Example:

```go
package modifiers

import (
    "context"
    "testing"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRequireReplaceIntAttributePlanModifier(t *testing.T) {
    modifier := RequireReplaceIntAttributePlanModifier()
    ctx := context.Background()

    cases := []struct {
        desc       string
        planValue  types.Int64
        stateValue types.Int64
        wantReplace bool
    }{
        {
            desc: "Force replacement when changed from non-zero",
            planValue: types.Int64Value(7),
            stateValue: types.Int64Value(1),
            wantReplace: true,
        },
        {
            desc: "Do not force replacement when state is zero",
            planValue: types.Int64Value(1),
            stateValue: types.Int64Value(0),
            wantReplace: false,
        },
        // ...other cases for null and unknown...
    }

    for _, tc := range cases {
        resp := &planmodifier.Int64Response{}
        req := planmodifier.Int64Request{
            PlanValue:  tc.planValue,
            StateValue: tc.stateValue,
        }
        modifier.PlanModifyInt64(ctx, req, resp)
        if resp.RequiresReplace != tc.wantReplace {
            t.Errorf("%s: got RequiresReplace=%v, want %v", tc.desc, resp.RequiresReplace, tc.wantReplace)
        }
    }
}
```
