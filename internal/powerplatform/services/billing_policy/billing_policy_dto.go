package billing_policy

type BillingPolicyDto struct {
	Id                string               `json:"id"`
	Name              string               `json:"name"`
	Location          string               `json:"location"`
	Status            string               `json:"status"`
	BillingInstrument BillingInstrumentDto `json:"billing_instrument"`
}

type BillingInstrumentDto struct {
	Id             string `json:"id"`
	ResourceGroup  string `json:"resource_group"`
	SubscriptionId string `json:"subscription_id"`
}
