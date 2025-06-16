# Title

Absence of DTO Type Definitions in Context

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/models.go

## Problem

The function `convertDtoToModel` depends on the type `AnalyticsDataDto` (and implied subtypes, such as its `Sink` and `Status` contents), but these definitions are not included in the file or imported. As a result, maintainability is hampered because the mapping logic is unclear without reference to the fields and types of the source DTOs. Contributors reading this file in isolation will find it hard to understand the structure or constraints of incoming data.

## Impact

Severity: **Low**

This reduces readability and maintainability, making onboarding and editing more difficult, as anyone attempting to review or change the mapping logic will need to search elsewhere for type definitions.

## Location

All references to `dto *AnalyticsDataDto` and its subfields in `convertDtoToModel`.

```go
func convertDtoToModel(dto *AnalyticsDataDto) *AnalyticsDataModel {
```

## Fix

Either:

- Add a comment with the definition of `AnalyticsDataDto` and its subtypes, or
- Import or link the file/module where the definitions are found, or
- Move critical struct definitions into this file (if they are short and only used here).

At a minimum, add a doc-comment or a link/reference to where the DTO type(s) are defined.

```go
// AnalyticsDataDto is defined in <path/to/definition>. Please update mapping here if the DTO changes.
func convertDtoToModel(dto *AnalyticsDataDto) *AnalyticsDataModel {
```
