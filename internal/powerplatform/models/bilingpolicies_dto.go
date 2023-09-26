package powerplatform_models

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

type BillingPolicyDtoArray struct {
	Value []BillingPolicyDto `json:"value"`
}
