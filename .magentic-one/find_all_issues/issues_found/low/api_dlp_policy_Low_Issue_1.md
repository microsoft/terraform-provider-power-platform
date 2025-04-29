# Title

Inefficient String Splitting in `convertPolicyModelToDlpPolicy`

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

The `strings.Split` function is called multiple times to split the connector ID, which is inefficient and unnecessarily verbose.

## Impact

While the efficiency impact is negligible in this specific implementation, repeated string splitting reduces code readability. Severity marked as **Low**.

## Location

`func convertPolicyModelToDlpPolicy(policy dlpPolicyModelDto) dlpPolicyDto`

## Code Issue

```go
nameSplit := strings.Split(connector.Id, "/")
Name: nameSplit[len(nameSplit)-1],
```

## Fix

Utilize a helper function for improved readability and reduce redundancy.

```go
func GetLastSegment(id string) string {
parts := strings.Split(id, "/")
return parts[len(parts)-1]
}

... // Inside the loop
Name: GetLastSegment(connector.Id),
```