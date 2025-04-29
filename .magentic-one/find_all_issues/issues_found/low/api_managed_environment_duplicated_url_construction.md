# Title

Unnecessary duplication of URL construction logic

##

`/workspaces/terraform-provider-power-platform/internal/services/managed_environment/api_managed_environment.go`

## Problem

Both the `EnableManagedEnvironment` and `DisableManagedEnvironment` functions contain duplicative logic for constructing the API URL. The only difference in these functions’ URL paths lies in the `environmentId` section. This repetition violates the DRY (Don't Repeat Yourself) principle and increases the risk of inconsistencies when updating any URL-related code. Additionally, it makes the functions less modular and harder to maintain.

## Impact

This duplication:
1. Decreases code maintainability.
2. Introduces the risk of inconsistencies if the logic for constructing API URLs needs to be updated.
3. Results in bloated functions, making the code harder to read and test.

Severity: **Low**

## Location

`/workspaces/terraform-provider-power-platform/internal/services/managed_environment/api_managed_environment.go`

## Code Issue

### EnableManagedEnvironment

```go
apiUrl := &url.URL{
    Scheme: constants.HTTPS,
    Host:   client.Api.GetConfig().Urls.BapiUrl,
    Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
}
values := url.Values{}
values.Add("api-version", "2021-04-01")
apiUrl.RawQuery = values.Encode()
```

### DisableManagedEnvironment

```go
apiUrl := &url.URL{
    Scheme: constants.HTTPS,
    Host:   client.Api.GetConfig().Urls.BapiUrl,
    Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
}
values := url.Values{}
values.Add("api-version", "2021-04-01")
apiUrl.RawQuery = values.Encode()
```

## Fix

Introduce a shared function for constructing the API URLs. This will eliminate duplication, making the code easier to maintain and more modular. For example:

### Shared Utility Function

```go
func constructManagedEnvApiUrl(baseHost string, environmentId string, apiVersion string) string {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   baseHost,
        Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
    }
    values := url.Values{}
    values.Add("api-version", apiVersion)
    apiUrl.RawQuery = values.Encode()
    return apiUrl.String()
}
```

### Refactored `EnableManagedEnvironment`

```go
apiUrl := constructManagedEnvApiUrl(client.Api.GetConfig().Urls.BapiUrl, environmentId, "2021-04-01")
// Use apiUrl in your logic
```

### Refactored `DisableManagedEnvironment`

```go
apiUrl := constructManagedEnvApiUrl(client.Api.GetConfig().Urls.BapiUrl, environmentId, "2021-04-01")
// Use apiUrl in your logic
```

This approach ensures that any changes to the API URL format are centralized, improving maintainability and reducing the risk of errors.
