# Type Safety: Potentially Unhandled Error for JSON Unmarshal

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The code unmarshals `response.BodyAsBytes` to `languages` but does not verify if the payload is indeed valid JSON or is non-empty before trying to unmarshal. While an error is returned if unmarshalling fails, a more defensive check would help with debugging.

## Impact

Poor error diagnosis in cases of empty or malformed responses. Severity: **medium**.

## Location

```go
err = json.Unmarshal(response.BodyAsBytes, &languages)

if err != nil {
	return languages, err
}
```

## Code Issue

```go
err = json.Unmarshal(response.BodyAsBytes, &languages)

if err != nil {
	return languages, err
}
```

## Fix

Optionally, check that `response.BodyAsBytes` is not empty before unmarshalling:

```go
if len(response.BodyAsBytes) == 0 {
    return languages, fmt.Errorf("empty response body")
}
err = json.Unmarshal(response.BodyAsBytes, &languages)
if err != nil {
    return languages, err
}
```
