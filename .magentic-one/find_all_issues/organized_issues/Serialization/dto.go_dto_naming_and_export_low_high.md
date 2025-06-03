# Title

Excessive Use of Abbreviations in Type and Field Names Reduces Readability

##

/workspaces/terraform-provider-power-platform/internal/services/connection/dto.go

## Problem

Almost all of the struct types are suffixed with `Dto`, e.g., `connectionDto`, `statusDto`, `createdByDto`. While it's common to distinguish DTOs from domain types, Go idioms recommend full words and capitalized names (`ConnectionDTO`). Also, excessive Hungarian notation (a la DTO) may be unnecessary if these types are only used for JSON unmarshaling.

## Impact

Low. This does not cause code errors, but it reduces readability and may lead to confusion or extra verbosity.

## Location

Every type in the file, e.g.:

```go
type connectionDto struct { ... }
type connectionPropertiesDto struct { ... }
type statusDto struct { ... }
...
```

## Code Issue

```go
type connectionDto struct { ... }
```

## Fix

Use full, capitalized names for exported types if needed. Remove the suffix if it is not essential for disambiguation. For example:

```go
type Connection struct { ... }
type ConnectionProperties struct { ... }
```

If you must keep the `DTO` suffix, use uppercase for clarity:

```go
type ConnectionDTO struct { ... }
```

And only export types (capitalized) if they are used outside the package.

---

# Title

Structs Not Exported Even Though JSON (Un)Marshaling May Require Exported Fields

##

/workspaces/terraform-provider-power-platform/internal/services/connection/dto.go

## Problem

All struct types and their fields are unexported (start with lowercase), but they might need to be exported (start with uppercase) for encoding/json and other packages outside this package to (un)marshal them correctly. In Go, fields must be exported to be marshaled/unmarshaled.

## Impact

High. If these types are intended to be used outside this package, or if JSON (un)marshaling occurs outside this package, unexported fields will be ignored, causing silent bugs.

## Location

Every type and most fields, e.g.:

```go
type connectionDto struct {
	Name       string                  `json:"name"`
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties connectionPropertiesDto `json:"properties"`
}
```

Here, the struct and its fields are unexported.

## Code Issue

```go
type connectionDto struct {
	Name       string                  `json:"name"`
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties connectionPropertiesDto `json:"properties"`
}
```

## Fix

Export all struct types and fields that are (un)marshaled or needed outside the package:

```go
type ConnectionDTO struct {
	Name       string                  `json:"name"`
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties ConnectionPropertiesDTO `json:"properties"`
}
```

Do this for every type/field that needs to be (un)marshaled or exported.
