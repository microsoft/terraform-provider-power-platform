# No Return on Error When Reading Files in importSolution

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem
In the `importSolution` function, when reading the solution file or settings file, the code adds an error to the diagnostics but continues execution instead of halting the import. This can result in the code later working with invalid or empty file contents, causing misleading errors or corrupt API requests:

```go
solutionContent, err := os.ReadFile(plan.SolutionFile.ValueString())
if err != nil {
    diagnostics.AddError(fmt.Sprintf("Client error when reading solution file %s", plan.SolutionFile.ValueString()), err.Error())
}
```
There is no `return` or other control-flow change after an error. The same pattern occurs for the settings file.

## Impact
- **Severity:** Critical
- Can result in the API being called with garbage or empty file contents if file read failed
- Corrupted/partial data can be silently processed
- Makes diagnostics and troubleshooting very difficult
- Increases risk of unpredictable behavior during resource creation or update

## Location
The error handling for reading files in the `importSolution` function:

## Code Issue
```go
solutionContent, err := os.ReadFile(plan.SolutionFile.ValueString())
if err != nil {
    diagnostics.AddError(fmt.Sprintf("Client error when reading solution file %s", plan.SolutionFile.ValueString()), err.Error())
}
//... continues

settingsContent, err = os.ReadFile(plan.SettingsFile.ValueString())
if err != nil {
    diagnostics.AddError(fmt.Sprintf("Client error when reading settings file %s", plan.SettingsFile.ValueString()), err.Error())
}
```

## Fix
Return immediately if file I/O fails so corrupted/incomplete file contents are not passed to downstream logic:

```go
solutionContent, err := os.ReadFile(plan.SolutionFile.ValueString())
if err != nil {
    diagnostics.AddError(fmt.Sprintf("Client error when reading solution file %s", plan.SolutionFile.ValueString()), err.Error())
    return nil
}
// ...

if plan.SettingsFile.ValueString() != "" {
    settingsContent, err = os.ReadFile(plan.SettingsFile.ValueString())
    if err != nil {
        diagnostics.AddError(fmt.Sprintf("Client error when reading settings file %s", plan.SettingsFile.ValueString()), err.Error())
        return nil
    }
}
```
