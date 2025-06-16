# Use of Magic Number for Expand Schema Depth

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

The call to `returnExpandSchema(10)` in the schema definition uses a magic number (the value 10) directly in code as the maximum recursion depth for the expand schema. This is not self-explanatory and could confuse maintainers or lead to future bugs if the value needs to be adjusted, as there is no explanation or symbolic constant.

## Impact

- Makes code less maintainable (not self-documenting).
- Changing the maximum depth requires hunting for the magic value.
- Severity: **low**.

## Location

In `Schema`:
```go
"expand":   returnExpandSchema(10),
```

## Code Issue

```go
"expand":   returnExpandSchema(10),
```

## Fix

Define a package-level constant for expand recursion depth with a descriptive name and short comment, and use the constant in the schema:

```go
const maxExpandRecursionDepth = 10

...

"expand":   returnExpandSchema(maxExpandRecursionDepth),
```
This documents intent, makes future maintenance easier, and aligns with Go best practices.
