# Title

Hardcoded Filter Query String in Function `GetConnectors`

##

`/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go`

## Problem

In the `GetConnectors` function, the filter query string for `$filter` is hardcoded. Specifically: 
```go
values.Add("$filter", "environment eq '~Default'")
```
This forces the function to always use the `~Default` environment, reducing its flexibility for use cases that involve querying connectors in other environments.

## Impact

Hardcoding the environment filter limits the reusability of the `GetConnectors` function. It prevents calling code from dynamically setting the target environment and will require modification to support other environments in the future. The severity of this issue is **medium**, as it directly affects function extensibility.

## Location

```go
values.Add("$filter", "environment eq '~Default'")
```

## Code Issue

```go
values.Add("$filter", "environment eq '~Default'")
```

## Fix

The function should accept the environment as a parameter, allowing the caller to specify which environment to use in the filter. Modify the function signature to include an additional parameter, such as `environment`, and replace the hardcoded value with the dynamic parameter.

```go
func (client *client) GetConnectors(ctx context.Context, environment string) ([]connectorDto, error) {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.PowerAppsUrl,
        Path:   "/providers/Microsoft.PowerApps/apis",
    }
    values := url.Values{}
    values.Add("api-version", "2019-05-01")
    values.Add("showApisWithToS", "true")
    values.Add("hideDlpExemptApis", "true")
    values.Add("showAllDlpEnforceableApis", "true")
    values.Add("$filter", fmt.Sprintf("environment eq '%s'", environment)) // Use dynamic environment

    apiUrl.RawQuery = values.Encode()

    connectorArray := connectorArrayDto{}
    _, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
    if err != nil {
        return nil, err
    }

    // Remaining code unchanged...
}
```

This change makes the function flexible and allows it to be used for targeting connectors in different environments. The caller can also easily specify the environment during function calls.