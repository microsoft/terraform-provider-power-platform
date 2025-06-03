# Issue: Unhandled errors when parsing headers in AddDataverseToEnvironment

##

/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go

## Problem

When parsing headers such as the Location and Retry-After headers after environment creation, failures are only logged (`tflog.Error`) and execution usually continues, which can lead to further issues or panics. If, for example, the `locationHeader` is empty or invalid, the subsequent logic depending on it could fail.

## Impact

- Severity: High
- This can cause subtle bugs, panics, or repeated API errors if an invalid URL or duration is used.
- The current approach provides no feedback to the calling function, so error propagation is unclear.

## Location

Within `AddDataverseToEnvironment`:

```go
locationHeader := apiResponse.GetHeader(constants.HEADER_LOCATION)
tflog.Debug(ctx, "Location Header: "+locationHeader)

_, err = url.Parse(locationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
}

retryHeader := apiResponse.GetHeader(constants.HEADER_RETRY_AFTER)
tflog.Debug(ctx, "Retry Header: "+retryHeader)
retryAfter, err := time.ParseDuration(retryHeader)
if err != nil {
    retryAfter = api.DefaultRetryAfter()
} else {
    retryAfter = retryAfter * time.Second
}
```

## Code Issue

```go
_, err = url.Parse(locationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
}

retryHeader := apiResponse.GetHeader(constants.HEADER_RETRY_AFTER)
tflog.Debug(ctx, "Retry Header: "+retryHeader)
retryAfter, err := time.ParseDuration(retryHeader)
if err != nil {
    retryAfter = api.DefaultRetryAfter()
} else {
    retryAfter = retryAfter * time.Second
}
```

## Fix

Error out if critical information cannot be parsed, such as the location header or retry interval.

```go
locationHeader := apiResponse.GetHeader(constants.HEADER_LOCATION)
if locationHeader == "" {
    return nil, errors.New("missing Location header in API response")
}
tflog.Debug(ctx, "Location Header: "+locationHeader)

_, err = url.Parse(locationHeader)
if err != nil {
    return nil, fmt.Errorf("error parsing location header: %w", err)
}

retryHeader := apiResponse.GetHeader(constants.HEADER_RETRY_AFTER)
tflog.Debug(ctx, "Retry Header: "+retryHeader)
retryAfter, err := time.ParseDuration(retryHeader)
if err != nil {
    retryAfter = api.DefaultRetryAfter()
} else {
    retryAfter = retryAfter * time.Second
}
```

---

This issue should be saved under:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/api_environment_add_dataverse_header_parsing_high.md`
