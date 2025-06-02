# Receiver Naming: Non-idiomatic Receiver Name

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The receiver variable is called `client`, which is identical to the type name. By convention, short, lower-case receiver names (often the first letter of the type) should be used to prevent confusion.

## Impact

Potentially confusing code and reduced readability. Severity: **low**.

## Location

```go
func (client *client) GetLanguagesByLocation(ctx context.Context, location string) (languagesArrayDto, error)
```

## Code Issue

```go
func (client *client) GetLanguagesByLocation(ctx context.Context, location string) (languagesArrayDto, error)
```

## Fix

Change the receiver to `c` or another short lowercase identifier:

```go
func (c *client) GetLanguagesByLocation(ctx context.Context, location string) (languagesArrayDto, error)
```
