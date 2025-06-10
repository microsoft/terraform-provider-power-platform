# Title

Hardcoded Import ID Format Reduces Flexibility and Robustness

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/api_copilot_studio_application_insights.go

## Problem

The function `parseImportId` assumes the import ID will always be in the format "envId_botId" and splits by an underscore `'_'`. This is a fragile design; if future import ID formats change (e.g., use hyphens or include more parts), this code will fail or behave incorrectly. Also, there's no documentation in the function explaining this constraint.

## Impact

This problem reduces maintainability and introduces a silent point of failure for future extensibility. It is of **medium severity**, as a malformed import ID results in providers failing with unhelpful errors.

## Location

```go
func parseImportId(importId string) (envId string, botId string, err error) {
	parts := strings.Split(importId, "_")
	if len(parts) != 2 {
		return "", "", errors.New("invalid import id format")
	}
	return parts[0], parts[1], nil
}
```

## Fix

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
