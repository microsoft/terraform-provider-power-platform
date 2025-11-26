# `powerplatform_billing_policy`

This resource is used to manage a Power Platform Billing Policy. A billing policy is linked to an Azure subscription and is used to set up pay-as-you-go billing for Power Platform environments.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Create              | `POST`      | `https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview` |
| Read                | `GET`       | `https://api.powerplatform.com/licensing/billingPolicies/{billing_policy_id}?api-version=2022-03-01-preview` |
| Update              | `PUT`       | `https://api.powerplatform.com/licensing/billingPolicies/{billing_policy_id}?api-version=2022-03-01-preview` |
| Delete              | `DELETE`    | `https://api.powerplatform.com/licensing/BillingPolicies/{billing_policy_id}?api-version=2022-03-01-preview` |

## Attribute Mapping

| Resource Attribute | API Request/Response JSON Field |
| ------------------ | ------------------------------- |
| `id`               | `id` |
| `name`             | `name` |
| `location`         | `location` |
| `status`           | `status` |
| `billing_instrument.id` | `billingInstrument.id` |
| `billing_instrument.resource_group` | `billingInstrument.resourceGroup` |
| `billing_instrument.subscription_id` | `billingInstrument.subscriptionId` |

### Example API Response

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "name": "payAsYouGoBillingPolicyExample",
  "location": "europe",
  "status": "Enabled",
  "billingInstrument": {
    "id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg",
    "resourceGroup": "example-rg",
    "subscriptionId": "00000000-0000-0000-0000-000000000000"
  }
}
```

## Notes

- The billing policy status may transition through intermediate states (e.g., "Enabling", "Disabling") before reaching terminal states ("Enabled", "Disabled")
- The provider waits for the billing policy to reach a terminal status during Create and Update operations
- The `location`, `billing_instrument.resource_group`, and `billing_instrument.subscription_id` attributes require resource replacement when changed

## References

- [What is a billing policy](https://learn.microsoft.com/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy)
- [Power Platform Billing Policy API](https://learn.microsoft.com/rest/api/power-platform/licensing/billing-policy/get-billing-policy)
