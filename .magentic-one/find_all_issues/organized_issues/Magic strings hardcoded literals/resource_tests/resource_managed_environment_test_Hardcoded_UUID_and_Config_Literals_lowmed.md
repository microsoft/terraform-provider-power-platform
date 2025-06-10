# Title

Wrong Use of Hardcoded UUID and Config Strings

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment_test.go

## Problem

Multiple hardcoded GUIDs, URLs, values, and so forth appear in multiple tests and config strings. This is somewhat expected in testing, but risks accidental copy-paste errors and does not centralize test fixtures or allow for adjustment in the event the GUID must change globally.

## Impact

Low-Medium. Not centralizing fixtures or magic values can make maintenance harder; bugs can slip in if a required value changes and is not updated everywhere.

## Location

E.g.,

```go
config: "... security_group_id = \"00000000-0000-0000-0000-000000000000\" ..."
httpmock.RegisterResponder("GET", "...environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?"
...
"powerplatform_managed_environment\" \"managed_development\" { environment_id = \"00000000-0000-0000-0000-000000000001\"
```

## Code Issue

```go
security_group_id = "00000000-0000-0000-0000-000000000000"
// ...etc...
```

## Fix

Centralize such test fixtures/constants:

```go
const (
    testSecurityGroupID = "00000000-0000-0000-0000-000000000000"
    testEnvironmentID   = "00000000-0000-0000-0000-000000000001"
    // etc.
)

// Then use in config strings:
config := fmt.Sprintf(`
resource ... {
    ...
    security_group_id = "%s"
}
`, testSecurityGroupID)
```

