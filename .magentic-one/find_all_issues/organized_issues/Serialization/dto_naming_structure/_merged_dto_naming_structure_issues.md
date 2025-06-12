# DTO Naming Structure and Export Issues

This document consolidates all issues related to DTO naming structure, export visibility, and consistency across the codebase.

## ISSUE 1

### Inconsistent JSON Struct Tag Naming

**File:** `/workspaces/terraform-provider-power-platform/internal/services/data_record/dto.go`

**Problem:** JSON struct tag names are inconsistently cased. Some are camelCase ("LogicalName"), others are PascalCase ("MetadataId"), and others "@odata.context". Conventionally, JSON properties should follow lowerCamelCase for external APIs.

**Impact:** Medium severity: Inconsistent JSON output can confuse users/consumers of the API, is error-prone, and is counter to established best practices.

**Location:** All struct field tags in this file

**Code Issue:**

```go
type entityDefinitionsDto struct {
 OdataContext          string `json:"@odata.context"`
 PrimaryIDAttribute    string `json:"PrimaryIdAttribute"`
 LogicalCollectionName string `json:"LogicalCollectionName"`
 MetadataID            string `json:"MetadataId"`
}
// ...etc
```

**Fix:** Decide on a consistent naming convention (typically lowerCamelCase for JSON) and update tags accordingly.

```go
type EntityDefinitionsDto struct {
 ODataContext          string `json:"@odata.context"`
 PrimaryIdAttribute    string `json:"primaryIdAttribute"`
 LogicalCollectionName string `json:"logicalCollectionName"`
 MetadataId            string `json:"metadataId"`
}
```

## ISSUE 2

### Excessive Use of Abbreviations in Type and Field Names Reduces Readability

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connection/dto.go`

**Problem:** Almost all of the struct types are suffixed with `Dto`, e.g., `connectionDto`, `statusDto`, `createdByDto`. While it's common to distinguish DTOs from domain types, Go idioms recommend full words and capitalized names (`ConnectionDTO`). Also, excessive Hungarian notation (a la DTO) may be unnecessary if these types are only used for JSON unmarshaling.

**Impact:** Low. This does not cause code errors, but it reduces readability and may lead to confusion or extra verbosity.

**Location:** Every type in the file, e.g.:

```go
type connectionDto struct { ... }
type connectionPropertiesDto struct { ... }
type statusDto struct { ... }
...
```

**Code Issue:**

```go
type connectionDto struct { ... }
```

**Fix:** Use full, capitalized names for exported types if needed. Remove the suffix if it is not essential for disambiguation. For example:

```go
type Connection struct { ... }
type ConnectionProperties struct { ... }
```

If you must keep the `DTO` suffix, use uppercase for clarity:

```go
type ConnectionDTO struct { ... }
```

And only export types (capitalized) if they are used outside the package.

## ISSUE 3

### Structs Not Exported Even Though JSON (Un)Marshaling May Require Exported Fields

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connection/dto.go`

**Problem:** All struct types and their fields are unexported (start with lowercase), but they might need to be exported (start with uppercase) for encoding/json and other packages outside this package to (un)marshal them correctly. In Go, fields must be exported to be marshaled/unmarshaled.

**Impact:** High. If these types are intended to be used outside this package, or if JSON (un)marshaling occurs outside this package, unexported fields will be ignored, causing silent bugs.

**Location:** Every type and most fields, e.g.:

```go
type connectionDto struct {
 Name       string                  `json:"name"`
 Id         string                  `json:"id"`
 Type       string                  `json:"type"`
 Properties connectionPropertiesDto `json:"properties"`
}
```

Here, the struct and its fields are unexported.

**Code Issue:**

```go
type connectionDto struct {
 Name       string                  `json:"name"`
 Id         string                  `json:"id"`
 Type       string                  `json:"type"`
 Properties connectionPropertiesDto `json:"properties"`
}
```

**Fix:** Export all struct types and fields that are (un)marshaled or needed outside the package:

```go
type ConnectionDTO struct {
 Name       string                  `json:"name"`
 Id         string                  `json:"id"`
 Type       string                  `json:"type"`
 Properties ConnectionPropertiesDTO `json:"properties"`
}
```

Do this for every type/field that needs to be (un)marshaled or exported.

---

Apply this fix to the whole codebase

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
