# Lack of Contextual Error Information

**Path:** `/workspaces/terraform-provider-power-platform/internal/services/licensing/api_licensing.go`

---

## Description

Several functions (e.g., `CreateBillingPolicy`, `UpdateBillingPolicy`, and others) return errors without including sufficient context about where or why the error occurred. Errors are directly returned or wrapped but lack specific logging or contextual information.

---

## Observed Code

Observed in `CreateBillingPolicy`:

```go
func (client *Client) CreateBillingPolicy(ctx context.Context, policyToCreate billingPolicyCreateDto) (*BillingPolicyDto, error) {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
        Path:   "/licensing/BillingPolicies",
    }

    values := url.Values{}
    values.Add(constants.API_VERSION_PARAM, "2022-03-01-preview")
    apiUrl.RawQuery = values.Encode()

    policy := &BillingPolicyDto{}
    _, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, policyToCreate, []int{http.StatusCreated}, policy)
    if err != nil {
        return nil, err
    }

    if policy.Status != "Enabled" && policy.Status != "Disabled" {
        policy, err = client.DoWaitForFinalStatus(ctx, policy)

        if err != nil {
            return nil, err
        }
    }

    return policy, err
}
```

---

## Impact

Lack of contextual error information makes debugging and root-cause analysis significantly more difficult, especially for downstream consumers of the client API.

**Severity:** Medium

---

## Suggested Fix

Wrap errors returned from API functions with detailed contextual information (e.g., the requested operation, URL, and other metadata where applicable):

```go
func (client *Client) CreateBillingPolicy(ctx context.Context, policyToCreate billingPolicyCreateDto) (*BillingPolicyDto, error) {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
        Path:   "/licensing/BillingPolicies",
    }

    values := url.Values{}
    values.Add(constants.API_VERSION_PARAM, "2022-03-01-preview")
    apiUrl.RawQuery = values.Encode()

    policy := &BillingPolicyDto{}
    _, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, policyToCreate, []int{http.StatusCreated}, policy)
    if err != nil {
        return nil, fmt.Errorf("failed to create billing policy %v, error: %w", policyToCreate, err)
    }

    if policy.Status != "Enabled" && policy.Status != "Disabled" {
        policy, err = client.DoWaitForFinalStatus(ctx, policy)
        if err != nil {
            return nil, fmt.Errorf("failed waiting for terminal billing policy status, error: %w", err)
        }
    }

    return policy, nil
}
```

---

## Recommended Actions

- Update all API functions to return errors wrapped with detailed contextual information using `fmt.Errorf`.
- Log errors using `tflog` or another logging mechanism wherever appropriate before returning them.

---