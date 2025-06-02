# Title
Field names `requestUrl` and related fields should follow Go capitalization conventions

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The `OidcCredential` struct and its options (`requestUrl`, `tokenFilePath`, etc.) use mixed-case capitalization (e.g., "Url") instead of Go's convention ("URL"). Go prefers `URL`, `HTTPRequest`, etc., for acronyms. Similarly, these fields should likely be exported (`RequestURL`, etc.) for consistency with the rest of the package if they are intended for external use (or unexported if private).

## Impact
Low. Reduced code consistency and poor adherence to Go idioms; could cause confusion for maintainers.

## Location
Example in:
```go
requestUrl    string
tokenFilePath string
```
## Code Issue
```go
requestUrl    string
tokenFilePath string
```
## Fix
Use Go's preferred acronym capitalization:
```go
requestURL    string
tokenFilePath string
```
If these fields must be exported:
```go
RequestURL    string
TokenFilePath string
```
