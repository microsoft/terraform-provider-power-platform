# Unused or misleading variable name reuse within for-loops

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

In the `GetPolicies` method, the variable named `policy` is both the range variable in the main loop as well as redeclared as a new variable inside the loop. This shadowing can cause confusion about which variable is being used.

## Impact

This can lead to subtle bugs, make the code significantly harder to read and maintain, and generally reduces code clarity. (Severity: Medium)

## Location

Lines 29–44

## Code Issue

```go
	for _, policy := range policiesArray.Value {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   client.Api.GetConfig().Urls.BapiUrl,
			Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", policy.PolicyDefinition.Name),
		}
		policy := dlpPolicyDto{}
		_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)
		if err != nil {
			return nil, err
		}
		v, err := covertDlpPolicyToPolicyModelDto(policy)
		if err != nil {
			return nil, err
		}
		policies = append(policies, *v)
	}
```

## Fix

Rename the inner `policy` variable to something else such as `policyDto` to avoid shadowing the range variable.

```go
	for _, policy := range policiesArray.Value {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   client.Api.GetConfig().Urls.BapiUrl,
			Path:   fmt.Sprintf("providers/PowerPlatform.Governance/v2/policies/%s", policy.PolicyDefinition.Name),
		}
		policyDto := dlpPolicyDto{}
		_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policyDto)
		if err != nil {
			return nil, err
		}
		v, err := covertDlpPolicyToPolicyModelDto(policyDto)
		if err != nil {
			return nil, err
		}
		policies = append(policies, *v)
	}
```

