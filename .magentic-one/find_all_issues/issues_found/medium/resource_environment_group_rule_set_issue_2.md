# Title

Repeated Computation in `Schema` Method Affecting Performance

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set.go`

## Problem

In the `Schema` method, the computation for the range of `maxSharingRange` is performed inside the function every time it is invoked. It calculates an array of `[*big.Float]` values unnecessarily during runtime instead of defining this data as a static constant or initializing once during loading.

```go
maxSharingRange := []*big.Float{}
for i := -1; i < 100; i++ {
	maxSharingRange = append(maxSharingRange, big.NewFloat(float64(i)))
}
```

## Impact

- Repeated redundant computation impacts runtime performance, particularly when this operation is invoked frequently.
- This operation increases the overhead during API calls without offering any benefit when the range remains static.

**Severity: Medium**

## Location

Method: `Schema`

## Code Issue

```go
maxSharingRange := []*big.Float{}
for i := -1; i < 100; i++ {
	maxSharingRange = append(maxSharingRange, big.NewFloat(float64(i)))
}
```

## Fix

This computation should be moved to a static definition that initializes once during program loading.

```go
// Define the range globally as a static constant
var maxSharingRange = func() []*big.Float {
	rangeVals := []*big.Float{}
	for i := -1; i < 100; i++ {
		rangeVals = append(rangeVals, big.NewFloat(float64(i)))
	}
	return rangeVals
}()

// Schema Method Code
func (r *environmentGroupRuleSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Allows the creation of environment group rulesets. See [Power Platform documentation](https://learn.microsoft.com/power-platform/admin/environment-groups) for more information on the available rules that can be applied to an environment group.\n\n!> Known Issue: This resource only works with a user context and cannot be used at this time with a service principal. This is a limitation of the underlying API.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			// Other attributes follow similar structure
		},
	}
}
```

This change ensures the `maxSharingRange` computation happens only once during runtime initialization, significantly reducing redundancy for the method.