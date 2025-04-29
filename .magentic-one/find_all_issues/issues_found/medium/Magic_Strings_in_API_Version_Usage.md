# Title

Magic Strings in `api-version` Usage

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/api_environment_group_rule_set.go

## Problem

The `api-version` string is hardcoded multiple times in the file. Such usage is prone to errors during updates and does not conform to clean coding practices. 

## Impact

If the `api-version` value changes, developers would need to manually search and update in multiple locations, introducing the risk of inconsistency. This could lead to runtime errors and non-functioning API calls. Severity: Medium.

## Location

Occurrences in the following methods:
- `GetEnvironmentGroupRuleSet`: Line 34
- `CreateEnvironmentGroupRuleSet`: Line 78
- `UpdateEnvironmentGroupRuleSet`: Line 113
- `DeleteEnvironmentGroupRuleSet`: Line 162

## Code Issue

The pattern of hardcoding the `api-version` string is repeated:
```go
values.Add("api-version", "2021-10-01-preview")
```

## Fix

Create a constant to store the `api-version` string and use it across the methods. For example:

```go
const APIVersion = "2021-10-01-preview"
```

Use it as:
```go
values.Add("api-version", APIVersion)
```

Updated `GetEnvironmentGroupRuleSet` Method:
```go
func (client *Client) GetEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string) (*EnvironmentGroupRuleSetValueSetDto, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/environmentGroups/%s/ruleSets", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", APIVersion) // Use constant value
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := environmentGroupRuleSetDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, &environmentGroupRuleSet)
	if err != nil {
		return nil, err
	}

	if resp.HttpResponse.StatusCode == http.StatusNoContent {
		return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, "rule set '%s' not found")
	}

	if len(environmentGroupRuleSet.Value) == 0 {
		return nil, fmt.Errorf("no environment group ruleset found for environment group id %s", environmentGroupId)
	}

	return &environmentGroupRuleSet.Value[0], nil
}
```

Repeat the usage of the constant for all other instances in the file. This ensures easier updates and minimizes error risks.

---
