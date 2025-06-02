# Title

Hardcoded `api-version` Parameter in URL Construction

##

`/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go`

## Problem

The `api-version` parameter is hardcoded with the value `2020-10-01` in multiple places in the file. As API versions evolve, this hardcoding makes it more difficult to update the version across the codebase in a consistent and controlled manner.

## Impact 

- **Severity:** **Medium**
- Hardcoding increases maintenance costs when upgrading API versions globally.
- Introducing a version management mechanism reduces the likelihood of errors from manually updating individual references to `api-version`.

## Location

- Method `GetAdminApplication`: Line 21
- Method `RegisterAdminApplication`: Line 31
- Method `UnregisterAdminApplication`: Line 45

## Code Issue

```go
apiUrl := &url.URL{
    Scheme: constants.HTTPS,
    Host:   client.Api.GetConfig().Urls.BapiUrl,
    Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
    RawQuery: url.Values{
        "api-version": []string{"2020-10-01"},
    }.Encode(),
}
```

## Fix

Refactor the hardcoded `api-version` parameter by introducing a centralized constant or configuration for managing API versions. This makes it easier to update the API version globally in the future.

```go
const (
    CurrentApiVersion = "2020-10-01" // Define constant for API version management
)

...

func (client *client) GetAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.BapiUrl,
        Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
        RawQuery: url.Values{
            constants.API_VERSION_PARAM: []string{CurrentApiVersion},
        }.Encode(),
    }
    ...
}

func (client *client) RegisterAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.BapiUrl,
        Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
        RawQuery: url.Values{
            "api-version": []string{CurrentApiVersion},
        }.Encode(),
    }
    ...
}

func (client *client) UnregisterAdminApplication(ctx context.Context, clientId string) error {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.BapiUrl,
        Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
        RawQuery: url.Values{
            "api-version": []string{CurrentApiVersion},
        }.Encode(),
    }
    ...
}
```