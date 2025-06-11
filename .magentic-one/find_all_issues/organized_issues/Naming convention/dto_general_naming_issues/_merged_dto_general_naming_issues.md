# Dto General Naming Issues - Merged Issues

## ISSUE 1

# Inconsistent Field Naming and JSON Tag Usage

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/dto.go

## Problem

The code inconsistently applies omitempty to JSON tags. For example, `displayName` sometimes has omitempty, and sometimes does not, even when it's the same field (e.g., in `dlpPolicyDefinitionDto` and `dlpPolicyModelDto`). Some pointer fields have omitempty, but similar non-pointer fields do not.  

Additionally, some struct field names do not match Go naming conventions for acronyms (like `Id` should be `ID`, `ETag` instead of `ETag` due to Go conventions).

## Impact

Severity: **Medium**

- Unclear serialization behavior for consumers, leading to potential confusion or bugs when empty fields are included/excluded inconsistently.
- Lower readability and maintainability due to non-standard naming.

## Location

Example fields:

```go
type dlpPolicyModelDto struct {
	...
	ETag    string `json:"etag"`
	CreatedBy string `json:"createdBy"`
	...
}
type dlpEnvironmentDto struct {
	Name string `json:"name"`
	Id   string `json:"id"`   // Should be \"ID\"
	Type string `json:"type"` 
}
type dlpActionRuleDto struct {
	ActionId string `json:"actionId"`
	Behavior string `json:"behavior"`
}
```

## Fix

- Apply `omitempty` consistently for optional fields.
- Follow Go standard naming conventions for acronyms: use `ID`, not `Id`; `ETag`, not `ETag`; etc.
- Revise JSON tags for consistent casing and usage.

```go
type dlpActionRuleDto struct {
	ActionID string `json:"actionId"`
	Behavior string `json:"behavior"`
}
type dlpEnvironmentDto struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Type string `json:"type"`
}
```

Ensure all fields have consistent `omitempty` application based on their optionality across the codebase.



---

## ISSUE 2

# Inconsistent Struct Naming Conventions

## /workspaces/terraform-provider-power-platform/internal/services/environment/dto.go

## Problem

Several types in the file use different leading capitalization and struct naming conventions (e.g., `environmentArrayDto`, `environmentCreateDto`, `modifySkuDto`, etc.) versus the majority which use capitalized names (e.g., `EnvironmentDto`). In Go, exported types (usable outside the package) should always use capitalized, CamelCase names.

## Impact

While all types in this file may be internal, this inconsistency confuses both users and maintainers, as some types are exported while others are package-private. Go style recommends using capitalized names for exported types for consistency and clarity. Severity: **low**.

## Location

Throughout the entire file, e.g.
- `environmentArrayDto`
- `environmentCreateDto`
- `modifySkuDto`

## Code Issue

```go
type environmentArrayDto struct {
    Value []EnvironmentDto `json:"value"`
}

// ...

type modifySkuDto struct {
    EnvironmentSku string `json:"environmentSku,omitempty"`
}
```

## Fix

Rename these types to use capitalized CamelCase style, e.g.:

```go
type EnvironmentArrayDto struct {
    Value []EnvironmentDto `json:"value"`
}

type EnvironmentCreateDto struct {
    Location   string                         `json:"location"`
    Properties EnvironmentCreatePropertiesDto `json:"properties"`
}

type ModifySkuDto struct {
    EnvironmentSku string `json:"environmentSku,omitempty"`
}
```


---

## ISSUE 3

# Title

Struct Field Naming Inconsistency: Id vs ID

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/dto.go

## Problem

The fields representing identifiers are named `Id` (e.g., `Id string`) rather than `ID`, which is the Go convention. Acronyms and initialisms should use all capitals (i.e., `ID`), according to Go naming conventions. This also applies to JSON tags unless there is an external need to keep the lowercase `id` (for API compatibility).

## Impact

Reduces readability and breaks standard Go naming conventions. It can also cause confusion about what the field represents. The severity is **low** as it is mainly a style/convention problem.

## Location

All struct definitions where identifier fields are present.

## Code Issue

```go
type BillingInstrumentDto struct {
	Id             string `json:"id,omitempty"`
	// ...
}
type BillingPolicyDto struct {
	Id                string               `json:"id"`
	// ...
}
type PrincipalDto struct {
	Id            string `json:"id"`
	// ...
}
type BillingPolicyEnvironmentsDto struct {
	BillingPolicyId string `json:"billingPolicyId"`
	EnvironmentId   string `json:"environmentId"`
}
```

## Fix

Rename the fields to `ID` in Go, and update the JSON tag if a different casing is acceptable for your API.

```go
type BillingInstrumentDto struct {
	ID             string `json:"id,omitempty"`
	// ...
}
type BillingPolicyDto struct {
	ID                string               `json:"id"`
	// ...
}
type PrincipalDto struct {
	ID            string `json:"id"`
	// ...
}
type BillingPolicyEnvironmentsDto struct {
	BillingPolicyID string `json:"billingPolicyId"`
	EnvironmentID   string `json:"environmentId"`
}
```

Keep the JSON tags unchanged if you must adhere to an external API, but update the Go field names for consistency. Apply for the whole code base


---

## ISSUE 4

# Inconsistent Struct Naming

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/dto.go

## Problem

There is inconsistency in naming patterns for the DTO structs. Some structs end with `Dto` (e.g., `dlpEnvironmentDto`), some with `ModelDto` (e.g., `dlpPolicyModelDto`), and some use `ArrayDto` or similar suffixes. This may confuse maintainers and users about the intended usage or relationship of these types.

For example:

- `dlpPolicyModelDto`, `dlpPolicyDto`, `dlpPolicyDefinitionDto`, `dlpPolicyLastActionDto`, etc.
- `dlpConnectorGroupsModelDto` vs. `dlpConnectorGroupsDto`

## Impact

Severity: **Medium**

Inconsistent naming can reduce maintainability, make refactoring more difficult, and may lead to mistakes in usage, especially as the codebase grows or when tools rely on predictable naming patterns.

## Location

Throughout the file, e.g.:

```go
type dlpPolicyModelDto struct {
    ...
}
type dlpPolicyDto struct {
    ...
}
type dlpConnectorGroupsModelDto struct {
    ...
}
type dlpConnectorGroupsDto struct {
    ...
}
```

## Fix

Adopt a consistent naming convention for DTO struct names. Prefer using a single suffix (e.g., always use `Dto` for Data Transfer Objects). Remove redundant distinctions between `ModelDto`, `Dto`, and similar postfixes unless there is a very clear semantic distinction.

For example:

```go
// Instead of
type dlpPolicyModelDto struct { ... }
type dlpConnectorGroupsModelDto struct { ... }

// Use
type DlpPolicyDto struct { ... }
type DlpConnectorGroupsDto struct { ... }
```

Capitalize struct names if they are exported, and use consistent suffixes for DTOs. Apply the chosen style consistently throughout the file.



---

## ISSUE 5

# Unexported Struct Types Used as DTOs

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

The structs `userDto`, `securityRoleDto`, `securityRoleArrayDto`, etc., are named with an initial lowercase letter, making them unexported. Since these are Data Transfer Objects (DTOs), they are likely to be used across multiple packages, particularly for marshaling/unmarshaling JSON. Keeping them unexported restricts their usage, reduces potential reusability, and goes against Go naming conventions for types meant for sharing.

## Impact

**Severity: Medium**

DTOs that are unexported but are expected to be used outside of the current package cannot be accessed, resulting in the need for unnecessary wrapper types or copy-pasted structures in other packages. It also reduces testability from external packages and contradicts idiomatic Go naming guidelines for types expected for broader-use.

## Location

Multiple locations:
- Definition of `userDto`
- Definition of `securityRoleDto`
- Definition of `securityRoleArrayDto`
- ... (others follow this pattern)

## Code Issue

```go
type userDto struct {
  ...
}

type securityRoleDto struct {
  ...
}
```

## Fix

Capitalize the struct type names that are intended to be used outside this package. This will make them exported and accessible in other packages.

```go
type UserDto struct {
  ...
}

type SecurityRoleDto struct {
  ...
}
```

This applies to all relevant DTO structs in this file. If some types are intentionally kept private, please comment on the intention for clarity.


---

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
