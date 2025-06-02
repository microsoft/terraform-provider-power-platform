# Title

Use of `fmt.Sprintf` for URL concatenation in `BuildODataQueryFromModel`

##

`/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go`

## Problem

In the `BuildODataQueryFromModel` function, `fmt.Sprintf` is used for constructing a URL with a query string. This approach does not guarantee correct URL formatting or encoding, and it could lead to malformed URLs if the inputs are not properly sanitized. This is especially critical when dealing with dynamic data.

## Impact

Using `fmt.Sprintf` for URL concatenation can break the functionality of the generated URLs if special characters are present in the input, leading to incorrect query execution or runtime errors. The severity is **high**, as this can affect the interaction with external services and APIs.

## Location

`BuildODataQueryFromModel` function in `/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go`

## Code Issue

```go
if len(resultQuery) > 0 {
    return fmt.Sprintf("%s?%s", model.EntityCollection.ValueString(), resultQuery), headers, nil
}
```

## Fix

Use the `net/url` package's `URL` struct to construct URLs in a robust manner. This ensures proper encoding and formatting.

```go
if len(resultQuery) > 0 {
    baseUrl := url.URL{Path: model.EntityCollection.ValueString()}
    baseUrl.RawQuery = resultQuery
    return baseUrl.String(), headers, nil
}
```
