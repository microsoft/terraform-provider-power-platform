# `powerplatform_billing_policy_environment`

This resource is used to manage the association of Power Platform environments with a billing policy. It allows you to add or remove environments from a billing policy for pay-as-you-go billing.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Create/Update (Add) | `POST`      | `https://api.powerplatform.com/licensing/billingPolicies/{billing_policy_id}/environments/add?api-version=2022-03-01-preview` |
| Read                | `GET`       | `https://api.powerplatform.com/licensing/billingPolicies/{billing_policy_id}/environments?api-version=2022-03-01-preview` |
| Delete (Remove)     | `POST`      | `https://api.powerplatform.com/licensing/billingPolicies/{billing_policy_id}/environments/remove?api-version=2022-03-01-preview` |

## Attribute Mapping

| Resource Attribute  | API Request/Response JSON Field |
| ------------------- | ------------------------------- |
| `billing_policy_id` | URL path parameter |
| `environments`      | `environmentIds` (array) |

### Example API Request (Add Environments)

```json
{
  "environmentIds": [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]
}
```

### Example API Response (List Environments)

```json
{
  "value": [
    {
      "environmentId": "00000000-0000-0000-0000-000000000001"
    },
    {
      "environmentId": "00000000-0000-0000-0000-000000000002"
    }
  ]
}
```

## Notes

- This resource manages the many-to-many relationship between billing policies and environments
- The Create and Update operations both use the `/add` endpoint to add environments to the billing policy
- The Delete operation uses the `/remove` endpoint to remove environments from the billing policy
- Multiple environments can be added or removed in a single API call
- The resource uses the billing policy ID as a required input and manages the list of associated environment IDs

## References

- [What is a billing policy](https://learn.microsoft.com/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy)
- [Power Platform Billing Policy Environment API](https://learn.microsoft.com/rest/api/power-platform/licensing/billing-policy-environment/add-billing-policy-environments)
