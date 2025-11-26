# `powerplatform_billing_policies_environments` (Data Source)

This data source fetches the list of environments associated with a specific billing policy in a Power Platform tenant.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.powerplatform.com/licensing/billingPolicies/{billing_policy_id}/environments?api-version=2022-03-01-preview` |

## Attribute Mapping

| Data Source Attribute | API Response JSON Field |
| --------------------- | ----------------------- |
| `billing_policy_id`   | URL path parameter |
| `environments`        | `value[].environmentId` (array) |

### Example API Response

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

- This data source requires a billing policy ID as input
- Returns a list of environment IDs that are currently associated with the specified billing policy
- Returns an empty list if no environments are associated with the billing policy

## References

- [What is a billing policy](https://learn.microsoft.com/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy)
- [Power Platform Billing Policy Environments API](https://learn.microsoft.com/rest/api/power-platform/licensing/billing-policy-environment/list-billing-policy-environments)
