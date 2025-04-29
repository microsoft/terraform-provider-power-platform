# Title

Error Not Checked When Decoding `previousBytes` in `Delete` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

## Problem

Within the `Delete` function, the decoding operation using `json.Unmarshal` does not validate the success of the operation or if an error occurs. This can lead to silently corrupted states if the decoding fails or results in unexpected data for `originalSettings`.

## Impact

Silent failure of `json.Unmarshal` can corrupt the rollback operation during deletion due to the use of incompletely or incorrectly decoded data, leading to unrecoverable or unexpected tenant resource states. Severity: **Critical**.

## Location

Line 597: Inside the `Delete` function, in the block performing JSON decoding.

## Code Issue

```go
err2 := json.Unmarshal(previousBytes, &originalSettings)
if err2 != nil {
    resp.Diagnostics.AddError(
        "Error unmarshalling original settings", fmt.Sprintf("Error unmarshalling original settings: %s", err2.Error()),
    )
    return
}
```

## Fix

Refactor the decoding operation to validate the unmarshalling process and handle the error accordingly. Additionally, provide detailed diagnostics if the failure occurs.

```go
if previousBytes == nil {
    resp.Diagnostics.AddError(
        "Failed to Decode Original Settings",
        "No previous settings were found to decode and restore. Ensure state preservation before attempting deletion.",
    )
    return
}

err2 := json.Unmarshal(previousBytes, &originalSettings)
if err2 != nil {
    resp.Diagnostics.AddError(
        "Error Decoding Original Settings",
        fmt.Sprintf("Failed to decode original settings due to invalid JSON structure. Error: %s", err2.Error()),
    )
    return
}
```