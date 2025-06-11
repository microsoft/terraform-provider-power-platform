# Dto File Specific Issues - Merged Issues

## ISSUE 1

# Incorrect DTO Struct Naming Conventions

## 
/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/dto.go

## Problem

Several struct types use inconsistent or non-idiomatic naming conventions. For example, types such as `environmentGroupRuleSetDto`, `environmentGroupRuleSetValueDto`, and `environmentGroupRuleSetParameterDto` begin with a lowercase letter. In Go, types intended to be used outside their package should use PascalCase (start with uppercase), and exported fields should also start with an uppercase letter. Additionally, type names could be inconsistent regarding abbreviations and suffixes (`Dto`, `ValueDto`, `ParameterDto`, etc.). This makes the code less readable and confusing for maintainers.

## Impact

- **Severity:** Medium
- Types starting with a lowercase letter are unexported, which may restrict usage or testing in other packages.
- Inconsistent naming can lead to misunderstanding, potential misuse, and more difficult codebase navigation and documentation.
- Violates Go idiomatic naming which affects maintainability and collaboration.

## Location

Occurs in the following declarations and references throughout the file:

```go
type environmentGroupRuleSetDto struct {
    Value []EnvironmentGroupRuleSetValueSetDto `json:"value"`
}
type EnvironmentGroupRuleSetValueSetDto struct {
    Parameters        []*environmentGroupRuleSetParameterDto ...
    ...
}
type environmentGroupRuleSetEnvironmentFilterDto struct { ... }
type environmentGroupRuleSetValueTypeDto struct { ... }
type environmentGroupRuleSetValueDto struct { ... }
type environmentGroupRuleSetParameterDto struct { ... }
```
And their usage in function signatures/fields.

## Code Issue

```go
type environmentGroupRuleSetDto struct {  // not exported (lowercase 'e')
    Value []EnvironmentGroupRuleSetValueSetDto `json:"value"`
}

type environmentGroupRuleSetParameterDto struct { // not exported
    ...
}
```

## Fix

Follow Go naming conventions for type names, capitalizing the first letter if export is intended, and ensure naming consistency. For example:

```go
type EnvironmentGroupRuleSetDTO struct { // If export is intended
    Value []EnvironmentGroupRuleSetValueSetDTO `json:"value"`
}

type EnvironmentGroupRuleSetParameterDTO struct {
    ...
}
```

Evaluate if all DTOs need to be exported (used outside package), and refactor accordingly for consistency.

apply for whole code base
---


---

## ISSUE 2

# Field Naming Convention Issue

##

internal/services/copilot_studio_application_insights/dto.go

## Problem

The field names in the struct definitions use mixed naming conventions—most are singular CamelCase (e.g., `EnvironmentId`), but some, notably `Errors` and potentially `NetworkIsolation`, are ambiguous as to whether they represent collections or singular values. In Go, exported fields should follow clear and consistent naming conventions. 

## Impact

Potential confusion or misuse by contributors and API consumers, resulting in reduced maintainability and readability. Severity: **low**.

## Location

Line(s): 7-15

## Code Issue

```go
type CopilotStudioAppInsightsDto struct {
	EnvironmentId               string `json:"environmentId"`
	BotId                       string `json:"botId"`
	AppInsightsConnectionString string `json:"appInsightsConnectionString"`
	IncludeSensitiveInformation bool   `json:"includeSensitiveInformation"`
	IncludeActivities           bool   `json:"includeActivities"`
	IncludeActions              bool   `json:"includeActions"`
	Errors                      string `json:"errors"`
	NetworkIsolation            string `json:"networkIsolation"`
}
```

## Fix

Evaluate if `Errors` and `NetworkIsolation` are single values or should be named in singular or plural form accordingly. Rename for clarity if required. If `Errors` holds a single error message, use `Error`. If it's a list, use `[]string Errors`.

```go
type CopilotStudioAppInsightsDto struct {
	EnvironmentId               string   `json:"environmentId"`
	BotId                       string   `json:"botId"`
	AppInsightsConnectionString string   `json:"appInsightsConnectionString"`
	IncludeSensitiveInformation bool     `json:"includeSensitiveInformation"`
	IncludeActivities           bool     `json:"includeActivities"`
	IncludeActions              bool     `json:"includeActions"`
	Error                      string   `json:"error"`                    // if single error
	// or
	Errors                     []string `json:"errors"`                   // if multiple errors
	NetworkIsolation           string   `json:"networkIsolation"`
}
```


---

## ISSUE 3

# Potential JSON Marshalling/Unmarshalling Inconsistency Due to Unexported Struct Fields

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/dto.go

## Problem

Field `InstanceURL` in `linkedEnvironmentIdMetadataDto` is not exported (it’s not tagged for JSON, but naming still matters for serialization). While the use of JSON tags is not present, unexported struct fields are not accessible during JSON (un)marshalling. If the intent is to serialize/deserialize this field, it will be skipped.

## Impact

High severity: breaks serialization/deserialization logic if the field is supposed to be used over API boundaries and results in lost or omitted data.

## Location

- Line 18 (`linkedEnvironmentIdMetadataDto`)

## Code Issue

```go
type linkedEnvironmentIdMetadataDto struct {
    InstanceURL string
}
```

## Fix

Export the field, and consider adding JSON struct tags if they are needed.

```go
type LinkedEnvironmentIdMetadataDto struct {
    InstanceURL string `json:"instanceUrl"`
}
```


---

## ISSUE 4

# Title
Inconsistent Naming Convention for Struct Fields

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/dto.go

## Problem

The naming of some struct fields in this file is inconsistent. Specifically, many fields use `CamelCase` style (e.g., `EnvironmentId`, `OrganizationId`) while Go convention prefers `ID`, `URL`, etc., to be all upper-case (`EnvironmentID`, `OrganizationID`). Similarly, abbreviations should be all-caps as per Go idiomatic style.

## Impact

This has a **low** severity impact because it does not break functionality but affects code readability, maintainability, and consistency with Go standards. It could confuse developers or lead to mistakes and inconsistencies throughout the codebase.

## Location

- `EnvironmentDto` struct
- `SinkDto` struct
- `AnalyticsDataDto` struct

## Code Issue

```go
type EnvironmentDto struct {
	EnvironmentId  string `json:"environmentId"`
	OrganizationId string `json:"organizationId"`
}
type SinkDto struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	SubscriptionId    string `json:"subscriptionId,omitempty"`
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	ResourceName      string `json:"resourceName"`
	Key               string `json:"key"`
}

type AnalyticsDataDto struct {
	ID               string           `json:"id"`
	// ...
	AiType           string           `json:"aiType"`
}
```

## Fix

Update the struct fields to use Go naming conventions, using all-caps for common initialisms:

```go
type EnvironmentDto struct {
	EnvironmentID  string `json:"environmentId"`
	OrganizationID string `json:"organizationId"`
}

type SinkDto struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	SubscriptionID    string `json:"subscriptionId,omitempty"`
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	ResourceName      string `json:"resourceName"`
	Key               string `json:"key"`
}

type AnalyticsDataDto struct {
	ID               string   `json:"id"`
	// ...
	AIType           string   `json:"aiType"`
}
```



---

## ISSUE 5

# Issue 1

Inconsistent Struct Field Naming

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/dto.go

## Problem

Several struct fields do not follow Go's convention for initialisms or consistent naming, particularly `Id` and `Dto`. According to Go style guidelines, initialisms should be capitalized (`ID`, not `Id`), and suffixes like `DTO` should match the typical capitalization (e.g., `ConnectorDTO`).

## Impact

Lack of consistency in naming conventions can reduce code readability and maintainability, especially for teams familiar with Go idioms. This issue is of low severity but impacts the professional quality of the codebase.

## Location

Multiple struct definitions throughout the file.

## Code Issue

```go
type connectorDto struct {
	Name       string                 `json:"name"`
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties connectorPropertiesDto `json:"properties"`
}

type connectorArrayDto struct {
	Value []connectorDto `json:"value"`
}

type unblockableConnectorDto struct {
	Id       string                          `json:"id"`
	Metadata unblockableConnectorMetadataDto `json:"metadata"`
}

type unblockableConnectorMetadataDto struct {
	Unblockable bool `json:"unblockable"`
}

type virtualConnectorDto struct {
	Id       string                      `json:"id"`
	Metadata virtualConnectorMetadataDto `json:"metadata"`
}

type virtualConnectorMetadataDto struct {
	VirtualConnector bool   `json:"virtualConnector"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	DisplayName      string `json:"displayName"`
}
```

## Fix

Update struct and field names to use correct Go conventions for initialisms (`ID`) and consider capitalizing `DTO` in type names to reflect the typical Go style. However, renaming exported types/fields should be synced across the codebase. Here's a suggested fix for one struct as an example:

```go
type ConnectorDTO struct {
	Name       string                  `json:"name"`
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties ConnectorPropertiesDTO  `json:"properties"`
}

type ConnectorArrayDTO struct {
	Value []ConnectorDTO `json:"value"`
}

type UnblockableConnectorDTO struct {
	ID       string                          `json:"id"`
	Metadata UnblockableConnectorMetadataDTO `json:"metadata"`
}

type UnblockableConnectorMetadataDTO struct {
	Unblockable bool `json:"unblockable"`
}

type VirtualConnectorDTO struct {
	ID       string                      `json:"id"`
	Metadata VirtualConnectorMetadataDTO `json:"metadata"`
}

type VirtualConnectorMetadataDTO struct {
	VirtualConnector bool   `json:"virtualConnector"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	DisplayName      string `json:"displayName"`
}
```

Apply for the whole code base
---


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
