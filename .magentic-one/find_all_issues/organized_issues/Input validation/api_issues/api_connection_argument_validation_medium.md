# Title

Lack of Input Validation on Function Arguments

##

/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go

## Problem

The methods do not consistently validate input parameters such as `environmentId`, `connectorName`, and `connectionId`. Absence of argument validation can lead to misleading API calls, unexpected server errors, and harder to debug failures.

## Impact

Failing to validate required inputs can allow the code to make invalid API requests, potentially returning cryptic or misleading error responses and causing defects. Severity: Medium.

## Location

Example from `CreateConnection`:

```go
func (client *client) CreateConnection(ctx context.Context, environmentId, connectorName string, connectionToCreate createDto) (*connectionDto, error) {
	// No check on environmentId or connectorName
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildEnvironmentHostUri(environmentId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, strings.ReplaceAll(uuid.New().String(), "-", "")),
	}
...
```

## Code Issue

```go
func (client *client) CreateConnection(ctx context.Context, environmentId, connectorName string, connectionToCreate createDto) (*connectionDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildEnvironmentHostUri(environmentId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, strings.ReplaceAll(uuid.New().String(), "-", "")),
	}
	// ... (no input validation)
```

## Fix

Add argument validation for all public-facing methods:

```go
if environmentId == "" {
	return nil, fmt.Errorf("environmentId is required")
}
if connectorName == "" {
	return nil, fmt.Errorf("connectorName is required")
}
```
And so forth for other critical parameters.

---

*This file should be saved in:*
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/api_connection_argument_validation_medium.md
