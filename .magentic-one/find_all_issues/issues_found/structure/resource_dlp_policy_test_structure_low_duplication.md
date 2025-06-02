# Overly Duplicated Test Data Strings

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

The test file contains many large string constants comprising JSON payloads (policyResponse1–policyResponse6), with significant duplication between them. This bloat makes maintenance difficult and increases the likelihood of inconsistencies or errors when updating test data for new API features or fixes.

## Impact

**Severity: Low**  
While not affecting functionality, this code duplication makes the file difficult to maintain, increases merge conflict risk, and obscures test differences.

## Location

Example of duplicated responses:

```go
policyResponse1 := fmt.Sprintf(`{
    "policyDefinition": {
        "name": "%s",
        "displayName": "Block All Policy",
        ...
    },
    "customConnectorUrlPatternsDefinition": {
        "rules": [ ... ]
    }
}`, policyId)

policyResponse2 := fmt.Sprintf(`{
    "policyDefinition": {
        "name": "%s",
        "displayName": "Block All Policy_1",
        ...
    },
    "customConnectorUrlPatternsDefinition": {
        "rules": [ ... ]
    }
}`, policyId)

// ...Up to policyResponse6, mostly duplicating structure.
```

## Fix

Extract common portions to helper functions or use a `makePolicyResponse` function that takes only the necessary fields as arguments (e.g., name, displayName, groups, etc.), and generates the JSON payload as needed.

```go
func makePolicyResponse(name, displayName string, ...otherArgs) string {
    // Use fmt.Sprintf or a struct+json.Marshal to build JSON
}

// Use: 
policyResponse1 := makePolicyResponse(policyId, "Block All Policy", ...)
policyResponse2 := makePolicyResponse(policyId, "Block All Policy_1", ...)
```

Or, use Go's template engine or even defined Go structs marshaled to JSON for more robust and type-safe test data.
