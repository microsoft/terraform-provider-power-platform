# Title

Repetition of `Execute` Calls Causes Code Duplication

##

`/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go`

## Problem

The `GetConnectors` function includes multiple calls to the `client.Api.Execute` function that repeat similar logic for executing HTTP requests. For example:
```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
```
These calls are repeated for different purposes but share identical parameters except for the API URL and response object. This leads to code duplication and reduces maintainability.

## Impact

Repetition of similar code increases the likelihood of bugs and makes the function harder to maintain or extend. If changes are required, they need to be repeated across multiple locations, increasing potential risk. This is a **medium severity** issue because while it does not break functionality, it significantly impacts maintainability.

## Location

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch connectors: %w", err)
}

_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &unblockableConnectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch unblockable connectors: %w", err)
}

_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &virtualConnectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch virtual connectors: %w", err)
}
```

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch connectors: %w", err)
}

_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &unblockableConnectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch unblockable connectors: %w", err)
}

_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &virtualConnectorArray)
if err != nil {
    return nil, fmt.Errorf("failed to fetch virtual connectors: %w", err)
}
```

## Fix

Refactor the repetitive code into a helper function that executes HTTP requests and returns the necessary data structure. This eliminates duplication and ensures changes are made in one location. Here's an example:

```go
func executeRequest(ctx context.Context, client *api.Client, urlString string, responseObj any) error {
    _, err := client.Execute(ctx, nil, "GET", urlString, nil, nil, []int{http.StatusOK}, responseObj)
    if err != nil {
        return fmt.Errorf("failed to execute request to %s: %w", urlString, err)
    }
    return nil
}

func (client *client) GetConnectors(ctx context.Context) ([]connectorDto, error) {
    apiUrlBase := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.PowerAppsUrl,
        Path:   "/providers/Microsoft.PowerApps/apis",
    }
    values := url.Values{}
    // Set values before building URL
    values.Add("api-version", "2019-05-01")
    values.Add("showApisWithToS", "true")
    values.Add("hideDlpExemptApis", "true")
    values.Add("showAllDlpEnforceableApis", "true")
    values.Add("$filter", "environment eq '~Default'")

    apiUrlBase.RawQuery = values.Encode()

    connectorArray := connectorArrayDto{}
    if err := executeRequest(ctx, client.Api, apiUrlBase.String(), &connectorArray); err != nil {
        return nil, err
    }

    unblockableUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.BapiUrl,
        Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable",
    }
    unblockableConnectorArray := []unblockableConnectorDto{}
    if err := executeRequest(ctx, client.Api, unblockableUrl.String(), &unblockableConnectorArray); err != nil {
        return nil, err
    }

    virtualUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.BapiUrl,
        Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual",
    }
    virtualConnectorArray := []virtualConnectorDto{}
    if err := executeRequest(ctx, client.Api, virtualUrl.String(), &virtualConnectorArray); err != nil {
        return nil, err
    }

    // Logic for processing response objects remains unchanged
    return connectorArray.Value, nil
}
```

This refactor improves readability, reduces maintenance overhead, and makes future changes easier to implement.