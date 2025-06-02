# Title

Validation Missing for `clientId` Input in CRUD Methods

##

`/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go`

## Problem

Currently, no validation is being performed on the `clientId` parameter in the following methods:
- `client.GetAdminApplication`
- `client.RegisterAdminApplication`
- `client.UnregisterAdminApplication`

If `clientId` is empty or malformed, an invalid API URL will be constructed, which could lead to runtime errors or unhandled exceptions during API calls.

## Impact

- **Severity:** **High**
- Constructing incorrect URLs at runtime could cause failures in making API calls.
- Lack of validation obscures the source of error until it shows up much later during execution, making debugging more difficult.

## Location

- Method `GetAdminApplication`: Line 17
- Method `RegisterAdminApplication`: Line 27
- Method `UnregisterAdminApplication`: Line 41

## Code Issue

```go
func (client *client) GetAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			constants.API_VERSION_PARAM: []string{"2020-10-01"},
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
			"api-version": []string{"2020-10-01"},
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
			"api-version": []string{"2020-10-01"},
		}.Encode(),
	}

	...
}
```

## Fix

Add validation for the `clientId` input in each of these methods to ensure that it is non-empty and meets expected constraints. This can be implemented as follows:

```go
func validateClientId(clientId string) error {
	if clientId == "" {
		return fmt.Errorf("clientId cannot be empty")
	}
	// Add more specific validation logic if applicable
	return nil
}

func (client *client) GetAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	if err := validateClientId(clientId); err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			constants.API_VERSION_PARAM: []string{"2020-10-01"},
		}.Encode(),
	}
	...
}

func (client *client) RegisterAdminApplication(ctx context.Context, clientId string) (*adminManagementApplicationDto, error) {
	if err := validateClientId(clientId); err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			"api-version": []string{"2020-10-01"},
		}.Encode(),
	}
	...
}

func (client *client) UnregisterAdminApplication(ctx context.Context, clientId string) error {
	if err := validateClientId(clientId); err != nil {
		return err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/adminApplications/%s", clientId),
		RawQuery: url.Values{
			"api-version": []string{"2020-10-01"},
		}.Encode(),
	}
	...
}
```