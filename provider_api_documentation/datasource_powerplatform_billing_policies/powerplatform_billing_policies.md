# `powerplatform_billing_policies` (Data Source)

This data source fetches the list of billing policies in a Power Platform tenant. A billing policy is a set of rules that define how a tenant is billed for usage of Power Platform services. A billing policy is associated with a billing instrument, which is a subscription and resource group that is used to pay for usage of Power Platform services.

## API Endpoints

| Terraform Operation | HTTP Method | API Endpoint URL (Public Cloud) |
| ------------------- | ----------- | -------------------------------- |
| Read                | `GET`       | `https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview` |

## Attribute Mapping

| Data Source Attribute | API Response JSON Field |
| --------------------- | ----------------------- |
| `billing_policies` | `value` (array) |
| `billing_policies[].id` | `value[].id` |
| `billing_policies[].name` | `value[].name` |
| `billing_policies[].location` | `value[].location` |
| `billing_policies[].status` | `value[].status` |
| `billing_policies[].billing_instrument.id` | `value[].billingInstrument.id` |
| `billing_policies[].billing_instrument.resource_group` | `value[].billingInstrument.resourceGroup` |
| `billing_policies[].billing_instrument.subscription_id` | `value[].billingInstrument.subscriptionId` |

### Example API Response

```json
{
  "value": [
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
  ]
}
```

## Notes

- This data source returns all billing policies in the tenant
- Each billing policy includes its billing instrument details (Azure subscription and resource group)

## References

- [What is a billing policy](https://learn.microsoft.com/power-platform/admin/pay-as-you-go-overview#what-is-a-billing-policy)
- [Power Platform Billing Policies API](https://learn.microsoft.com/rest/api/power-platform/licensing/billing-policy/list-billing-policies)
