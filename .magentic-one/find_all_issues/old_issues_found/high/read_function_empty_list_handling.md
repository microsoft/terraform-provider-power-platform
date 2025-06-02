# Issue Report #2

### Title: Missing Handling for Empty List in `Read` Function

### Path to the file: `/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps.go`

---

## Problem

In the `Read` function of the `EnvironmentPowerAppsDataSource`, the `apps` variable is fetched using `d.PowerAppssClient.GetPowerApps(ctx)`. However, there is no handling for scenarios where the returned list is empty. This could lead to ambiguity or unintended behavior when the state is being set or appended.

---

## Impact

If the `apps` list is empty, the function may append no data, which could still appear as a valid fetch result to callers. This could lead to incorrect expectations from the data source. Severity: **High**

---

## Location

**Function:** `Read`, `apps` initialization and loop for app processing.

---

## Code Issue

```go
apps, err := d.PowerAppssClient.GetPowerApps(ctx)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}

for _, app := range apps {
    appModel := ConvertFromPowerAppDto(app)
    state.PowerApps = append(state.PowerApps, appModel)
}
```

---

## Fix

Introduce a condition to explicitly check if the list of apps is empty and add corresponding diagnostics for better feedback to the caller.

```go
apps, err := d.PowerAppssClient.GetPowerApps(ctx)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}

if len(apps) == 0 {
    resp.Diagnostics.AddWarning(
        fmt.Sprintf("No PowerApps found in %s", d.FullTypeName()),
        "The client returned an empty list, indicating no PowerApps are available in the specified environment.",
    )
    return
}

for _, app := range apps {
    appModel := ConvertFromPowerAppDto(app)
    state.PowerApps = append(state.PowerApps, appModel)
}
```