# Error Handling is not Comprehensive in `Read` Method

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go

## Problem

In the `Read` method, after calling `resp.State.Get(ctx, &state)`, the error returned from the Get operation is not checked. If an error occurs while retrieving the state, the code proceeds, potentially with a zero-value or invalid `state`, which could cause unexpected behaviors later in the function.

## Impact

This can lead to misleading diagnostics being returned to the user and may cause runtime panics, data mismatches, or further errors during the function execution. It also impacts debugging and maintainability as errors at this stage are silently ignored.

**Severity:** High

## Location

```go
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state DataSourceModel
	resp.State.Get(ctx, &state)

	currencies, err := d.CurrenciesClient.GetCurrenciesByLocation(ctx, state.Location.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}
    ...
}
```

## Fix

Check and handle the error returned by `resp.State.Get`. If it has an error, append diagnostics and return immediately.

```go
	var state DataSourceModel
	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
```
