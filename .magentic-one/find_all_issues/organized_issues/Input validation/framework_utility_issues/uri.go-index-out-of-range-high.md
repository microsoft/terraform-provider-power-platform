# Potential for Index Out of Range Panic Due to Inadequate Input Validation

##

/workspaces/terraform-provider-power-platform/internal/helpers/uri.go

## Problem

The functions `BuildEnvironmentHostUri` and `BuildTenantHostUri` lack input validation for the incoming `environmentId` or `tenantId` strings. These functions assume that, after removing hyphens, the identifier is at least two characters long. If a caller inadvertently or maliciously passes a shorter string, indexing with `envId[len(envId)-2:]` and `envId[:len(envId)-2]` will cause a runtime panic due to an "index out of range" error.

## Impact

**Severity: High**

A panic can terminate the overall process (such as a provider or an automation pipeline) unexpectedly, leading to unreliable software behavior. This is particularly critical in libraries or modules consumed by external layers, such as Terraform providers, where input may originate from user configuration.

## Location

Lines inside:
- `BuildEnvironmentHostUri`
- `BuildTenantHostUri`

## Code Issue

```go
envId := strings.ReplaceAll(environmentId, "-", "")
realm := string(envId[len(envId)-2:])
envId = envId[:len(envId)-2]
```

and

```go
envId := strings.ReplaceAll(tenantId, "-", "")
realm := string(envId[len(envId)-2:])
envId = envId[:len(envId)-2]
```

## Fix

Add input validation to ensure the processed ID has at least two characters, returning an empty string or error if not.

```go
func BuildEnvironmentHostUri(environmentId, powerPlatformUrl string) string {
	envId := strings.ReplaceAll(environmentId, "-", "")
	if len(envId) < 2 {
		// Optionally, log or handle the error accordingly.
		return ""
	}
	realm := envId[len(envId)-2:]
	envId = envId[:len(envId)-2]
	return fmt.Sprintf("%s.%s.environment.%s", envId, realm, powerPlatformUrl)
}

func BuildTenantHostUri(tenantId, powerPlatformUrl string) string {
	envId := strings.ReplaceAll(tenantId, "-", "")
	if len(envId) < 2 {
		// Optionally, log or handle the error accordingly.
		return ""
	}
	realm := envId[len(envId)-2:]
	envId = envId[:len(envId)-2]
	return fmt.Sprintf("%s.%s.tenant.%s", envId, realm, powerPlatformUrl)
}
```

---

This prevents panics by ensuring the input is sufficiently long before accessing string slices. Further handling (such as returning an error) could be used depending on the design requirements.
