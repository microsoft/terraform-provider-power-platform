# API Design and Parsing Issues

This document contains merged issues related to API design and parsing logic in the Power Platform Terraform provider.

## ISSUE 1

**Title:** Inefficient/Incorrect Parsing Logic of ApplicationId

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/api_application.go`

**Problem:**
When parsing the application ID from `lifecycleResponse.CreatedDateTime`, the code splits the string by `/`, which seems semantically incorrect and unreliable, as `CreatedDateTime` usually holds a date, not an identifier in a path format.

**Impact:**
This issue has a **medium** severity as it could cause a runtime error (index out of range) or return the wrong value if the API changes the format, making the provider unreliable.

**Location:**
Within InstallApplicationInEnvironment (inside the for loop):

**Code Issue:**
```go
parts := strings.Split(lifecycleResponse.CreatedDateTime, "/")
if len(parts) == 0 {
    return "", errors.New("can't parse application id from response " + lifecycleResponse.CreatedDateTime)
}
applicationId = parts[len(parts)-1]
tflog.Debug(ctx, "Created Application Id: "+applicationId)
```

**Fix:**
Verify what property should hold the application ID based on the DTO returned by the API. If it is a field other than `CreatedDateTime`, use the dedicated property.

If this logic is correct, handle dates accordingly. Otherwise, update to access the correct property:

```go
// If the DTO contains an ApplicationId field:
if lifecycleResponse.ApplicationId == "" {
    return "", errors.New("application id not present in lifecycle response")
}
applicationId = lifecycleResponse.ApplicationId
tflog.Debug(ctx, "Created Application Id: "+applicationId)
```

Or, if you must parse, ensure robust checking:

```go
parts := strings.Split(lifecycleResponse.CreatedDateTime, "/")
if len(parts) == 0 {
    return "", errors.New("can't parse application id from response " + lifecycleResponse.CreatedDateTime)
}
applicationId = parts[len(parts)-1]
tflog.Debug(ctx, "Created Application Id: "+applicationId)
```

But document and validate why `CreatedDateTime` is being split.

## ISSUE 2

**Title:** Hardcoded Import ID Format Reduces Flexibility and Robustness

**File:** `/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/api_copilot_studio_application_insights.go`

**Problem:**
The function `parseImportId` assumes the import ID will always be in the format "envId_botId" and splits by an underscore `'_'`. This is a fragile design; if future import ID formats change (e.g., use hyphens or include more parts), this code will fail or behave incorrectly. Also, there's no documentation in the function explaining this constraint.

**Impact:**
This problem reduces maintainability and introduces a silent point of failure for future extensibility. It is of **medium severity**, as a malformed import ID results in providers failing with unhelpful errors.

**Location:**
```go
func parseImportId(importId string) (envId string, botId string, err error) {
	parts := strings.Split(importId, "_")
	if len(parts) != 2 {
		return "", "", errors.New("invalid import id format")
	}
	return parts[0], parts[1], nil
}
```

**Fix:**
Add comments describing the required format. Consider supporting more robust parsing (e.g., regular expressions), and provide more context in error messages.

```go
// parseImportId parses an import ID in the form 'envId_botId'
// Example: 'e12345_b6789'
// Returns the environment ID and bot ID. Returns error if format is invalid.
// Future-proof: Consider refactoring if import ID formats change.
func parseImportId(importId string) (envId string, botId string, err error) {
	parts := strings.Split(importId, "_")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid import id format: expected 'envId_botId', got '%s'", importId)
	}
	return parts[0], parts[1], nil
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
