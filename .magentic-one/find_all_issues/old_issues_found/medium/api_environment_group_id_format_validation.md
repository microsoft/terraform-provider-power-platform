# Title

Non-Validation of Environment Group ID Format

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_groups/api_environment_group.go`

## Problem

Several methods, such as `DeleteEnvironmentGroup` and `RemoveEnvironmentFromEnvironmentGroup`, accept `environmentGroupId` and `environmentId` strings as input parameters without validating the format, presence of illegal characters, or ensuring their consistency with expected patterns.

For example:
```go
func (client *client) DeleteEnvironmentGroup(ctx context.Context, environmentGroupId string) error
```

Validation should ensure that the `environmentGroupId` and `environmentId` conform to an expected pattern (e.g., an alphanumeric string) before using them in critical API calls. Lack of input validation increases potential risks of invalid request errors or unexpected behavior at runtime.

## Impact

- **Severity:** Medium  
- Might lead to runtime errors or unexpected behavior. Additionally, passing invalid IDs could lead to unnecessary network requests, reducing API reliability and performance.

## Location

Multiple locations:

### Method:

`DeleteEnvironmentGroup`
```go
func (client *client) DeleteEnvironmentGroup(ctx context.Context, environmentGroupId string) error
```

### Method:

`RemoveEnvironmentFromEnvironmentGroup`
```go
func (client *client) RemoveEnvironmentFromEnvironmentGroup(ctx context.Context, environmentGroupId, environmentId string) error
```

## Code Issue

```go
func (client *client) DeleteEnvironmentGroup(ctx context.Context, environmentGroupId string) error
func (client *client) RemoveEnvironmentFromEnvironmentGroup(ctx context.Context, environmentGroupId, environmentId string) error
```

## Fix

Add input validation at the beginning of these methods to ensure IDs conform to expected patterns, such as alphanumeric strings.

### Example Fix:

```go
func validateID(id string) error {
	if len(id) == 0 || strings.ContainsAny(id, " !@#$%^&*()") {
		return errors.New("Invalid ID format")
	}
	return nil
}

func (client *client) DeleteEnvironmentGroup(ctx context.Context, environmentGroupId string) error {
	err := validateID(environmentGroupId)
	if err != nil {
		return err
	}

	// Proceed with the method logic...
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups/" + environmentGroupId,
	}
	// ...
}
```

Similarly, apply `validateID()` in `RemoveEnvironmentFromEnvironmentGroup` method for `environmentGroupId` and `environmentId`.
