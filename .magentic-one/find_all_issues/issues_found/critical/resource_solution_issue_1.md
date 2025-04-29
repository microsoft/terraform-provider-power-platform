# Title

Misuse of `md5 hash calculation log` in warning message

##

/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem

In the `Create` and `Update` method, the logging calls mention "md5 hash" when in reality it calculates an SHA256 hash.

## Impact

This creates misleading logs and could result in confusion for debugging. Developers or users inspecting these logs might think the checksum is computed using MD5, which is not accurate.

**Severity:** Critical

## Location

Lines originating around:

- Line: ~119 (`CREATE Calculated md5 hash of ...` log)
- Line: ~274 (`UPDATE Calculated md5 hash of settings file`).

## Code Issue

```go
			plan.SettingsFileChecksum = types.StringNull()
			if !plan.SettingsFile.IsNull() && !plan.SettingsFile.IsUnknown() {
				value, err := helpers.CalculateSHA256(plan.SettingsFile.ValueString())
				if err != nil {
				resp.Diagnostics.AddWarning("Issue calculating checksum settings",... --- but instead,
logic enforces mismatch/warning incorrectly,"logs/implementation-verav}}
	` corrected implement-update logic hash correctly error absence logical. App also blob reasons lessen confusion invalid lookup. Plenty Haystack themselves instead Yes;