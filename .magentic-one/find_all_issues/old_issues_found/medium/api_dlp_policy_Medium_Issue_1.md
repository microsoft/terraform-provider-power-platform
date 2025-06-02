# Title

Hardcoded Paths in `GetPolicies`

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

The `Path` field in the construction of `url.URL` objects within the `GetPolicies` method uses hardcoded strings. This approach is error-prone and reduces maintainability.

## Impact

Using hardcoded paths reduces flexibility when API paths are updated. Changes in paths would require locating and updating these hardcoded strings manually. Severity marked as **Medium**.

## Location

`func (client *client) GetPolicies(ctx context.Context) ([]dlpPolicyModelDto, error)`

## Code Issue

```go
Path: "providers/PowerPlatform.Governance/v2/policies",
```

## Fix

Put the path into a constant variable and use it instead.

```go
const GetPoliciesPath = "providers/PowerPlatform.Governance/v2/policies"
apiUrl := &url.URL{
Scheme: constants.HTTPS,
Host: client.Api.GetConfig().Urls.BapiUrl,
Path: GetPoliciesPath,
}
```