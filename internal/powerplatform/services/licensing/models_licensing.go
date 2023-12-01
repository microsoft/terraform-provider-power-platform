package powerplatform

type BillingPolicyCreateDto struct {
	Location          string               `json:"location"`
	Name              string               `json:"name"`
	Status            string               `json:"status"`
	BillingInstrument BillingInstrumentDto `json:"billingInstrument"`
}

type BillingInstrumentDto struct {
	Id             string `json:"id,omitempty"`
	ResourceGroup  string `json:"resourceGroup"`
	SubscriptionId string `json:"subscriptionId"`
}

type BillingPolicyDto struct {
	Id                string               `json:"id"`
	Name              string               `json:"name"`
	TenantType        string               `json:"type"`
	Status            string               `json:"status"`
	Location          string               `json:"location"`
	BillingInstrument BillingInstrumentDto `json:"billingInstrument"`
	CreatedOn         string               `json:"createdOn"`
	CreatedBy         PrincipalDto         `json:"createdBy"`
	LastModifiedOn    string               `json:"lastModifiedOn"`
	LastModifiedBy    PrincipalDto         `json:"lastModifiedBy"`
}

type BillingPolicyArrayDto struct {
	Value []BillingPolicyDto `json:"value"`
}

type BillingPolicyUpdateDto struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PrincipalDto struct {
	Id            string `json:"id"`
	PrincipalType string `json:"type"`
}

type BillingPolicyEnvironmentsArrayDto struct {
	EnvironmentIds []string `json:"environmentIds"`
}

type BillingPolicyEnvironmentsDto struct {
	BillingPolicyId string `json:"billingPolicyId"`
	EnvironmentId   string `json:"environmentId"`
}

type BillingPolicyEnvironmentsArrayResponseDto struct {
	Value []BillingPolicyEnvironmentsDto `json:"value"`
}
