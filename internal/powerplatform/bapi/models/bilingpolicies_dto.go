package powerplatform_bapi

// {
//     "value": [
//         {
//             "id": "e261f6bb-dddf-4691-98be-e6e2e6257536",
//             "name": "aaaasdasdasd",
//             "type": "TenantOwned",
//             "status": "Enabled",
//             "location": "europe",
//             "billingInstrument": {
//                 "id": "/subscriptions/2bc1f261-7e26-490c-9fd5-b7ca72032ad3/resourceGroups/tmp/providers/Microsoft.PowerPlatform/accounts/aaaasdasdasd",
//                 "subscriptionId": "2bc1f261-7e26-490c-9fd5-b7ca72032ad3",
//                 "resourceGroup": "tmp"
//             },
//             "createdOn": "2023-09-21T14:27:08Z",
//             "createdBy": {
//                 "id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
//                 "type": "User"
//             },
//             "lastModifiedOn": "2023-09-21T14:27:08Z",
//             "lastModifiedBy": {
//                 "id": "f99f844b-ce3b-49ae-86f3-e374ecae789c",
//                 "type": "User"
//             }
//         }
//     ]
// }

type BillingPolicyDto struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Status            string `json:"status"`
	Location          string `json:"location"`
	BillingInstrument struct {
		Id             string `json:"id"`
		SubscriptionId string `json:"subscriptionId"`
		ResourceGroup  string `json:"resourceGroup"`
	} `json:"billingInstrument"`
	CreatedOn string `json:"createdOn"`
	CreatedBy struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"createdBy"`
	LastModifiedOn string `json:"lastModifiedOn"`
	LastModifiedBy struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"lastModifiedBy"`
}

type BillingPoliciesDto struct {
	Value []BillingPolicyDto `json:"value"`
}
