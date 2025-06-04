# Missing Return After Error Handling Issues

This document consolidates all issues related to missing return statements after error handling in the Terraform Provider for Power Platform.

## ISSUE 1

### Missing Return After Error Handling in Read Method

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors.go`

**Problem:** In the `Read` method of the `DataSource` struct, when an error occurs during the call to `d.ConnectorsClient.GetConnectors(ctx)`, an error is appended to the diagnostics, but there is no `return` statement immediately following. As a result, the code continues to execute, possibly using a `connectors` value that may not be valid, which can lead to unintended side effects, panics, or corrupted state.

**Impact:** This is a **high severity** issue because error handling should prevent subsequent operations that depend on successful completion of the failed operation. Continuing after an error could result in runtime panics, data corruption, or misleading state within the Terraform provider.

**Location:**

```go
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    // ...
    connectors, err := d.ConnectorsClient.GetConnectors(ctx)
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), fmt.Errorf("error occurred: %w", err).Error())
    }

    for _, connector := range connectors {
        connectorModel := convertFromConnectorDto(connector)
        state.Connectors = append(state.Connectors, connectorModel)
    }
    // ...
}
```

**Fix:** Add a `return` statement immediately after appending the error to diagnostics:

```go
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    // ...
    connectors, err := d.ConnectorsClient.GetConnectors(ctx)
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), fmt.Errorf("error occurred: %w", err).Error())
        return
    }

    for _, connector := range connectors {
        connectorModel := convertFromConnectorDto(connector)
        state.Connectors = append(state.Connectors, connectorModel)
    }
    // ...
}
```

## ISSUE 2

### Logic Bug: No Early Return After Error When Checking Dataverse Existence

**File:** `/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions.go`

**Problem:** When `d.SolutionClient.DataverseExists` returns an error, the error is added to diagnostics, but the function does not immediately return. This could lead to ambiguous error reporting and further execution based on undefined state (`dvExits` will have the Go default value, usually `false`), potentially leading to misleading error messages or improper API usage.

**Impact:** **High**. May result in misleading output, reporting multiple errors for a single underlying cause, or attempting to interact with uninitialized/invalid data, potentially causing spurious or unclear diagnostics in the Terraform provider.

**Location:** Lines 109-115, in the `Read` method:

**Code Issue:**

```go
dvExits, err := d.SolutionClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}

if !dvExits {
 resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
 return
}
```

**Fix:** Return early after logging the error:

```go
dvExits, err := d.SolutionClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
 return
}

if !dvExits {
 resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
 return
}
```

## ISSUE 3

### Potential Panic Due to Unchecked resp.State.Get Error

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages.go`

**Problem:** In the Read method, `resp.State.Get(ctx, &state)` is called, but the error returned (diagnostics) is not checked. In unusual conditions (malformed state, framework bug, etc.), this could lead to panics or undefined data usage.

**Impact:**

- High: Possible panic or silent failure with bad input or Terraform state bugs.
- Error handling best practice.

**Location:** Read method, at start of method.

**Code Issue:**

```go
var state TenantApplicationPackagesListDataSourceModel
resp.State.Get(ctx, &state)
```

**Fix:** Capture diagnostics and handle errors:

```go
var state TenantApplicationPackagesListDataSourceModel
diags := resp.State.Get(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
 return
}
```

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
