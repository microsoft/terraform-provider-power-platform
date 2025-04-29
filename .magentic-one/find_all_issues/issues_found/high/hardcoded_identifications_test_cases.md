# Issue: Hardcoded Identifications in Test Cases

### Path
`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go`

### Problem
Test cases frequently use hardcoded identifiers like `00000000-0000-0000-0000-000000000001`, which makes the tests brittle and hard to adapt to new scenarios or configurations.

### Severity
High

### Suggested Fix
Generate dynamic values for test identifiers during test setup. This approach improves test adaptability and overall reliability.

### Proposed Code Change
```go
environmentID := GenerateTestGUID()
aadID := GenerateTestGUID()

Config: fmt.Sprintf(`
resource "powerplatform_user" "new_user" {
    environment_id = "%s"
    aad_id = "%s"
}`, environmentID, aadID)
```
