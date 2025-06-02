# Title

Missing Error Wrapping for `GetLocations`

##

`/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go`

## Problem

In the `GetLocations` method:

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
return locations, err
```

The error returned from `client.Api.Execute` is directly passed back without wrapping it with additional context. Adding context to the error would improve debugging and make it easy to trace issues back to their source.

## Impact

Without proper error wrapping, debugging becomes difficult when the application logs only surface-level error messages. This can particularly hamper troubleshooting efforts in production environments and impact maintainability.

**Severity: High**

## Location

- File: `api_locations.go`
- Function: `GetLocations`
- Line: `GetLocations method (error handling section)`

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
return locations, err
```

## Fix

Wrap the error with additional contextual information:

```go
import (
	"fmt"
	// ... other imports
)

func (client *client) GetLocations(ctx context.Context) (locationDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/locations",
		RawQuery: url.Values{
			"api-version": []string{"2023-06-01"},
		}.Encode(),
	}

	var locations locationDto
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
	if err != nil {
		return locations, fmt.Errorf("error fetching locations: %w", err) // Adding error context
	}
	return locations, nil
}
```