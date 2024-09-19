// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

type billingPolicyCreateDto struct {
	Location          string               `json:"location"`
	Name              string               `json:"name"`
	Status            string               `json:"status"`
	BillingInstrument billingInstrumentDto `json:"billingInstrument"`
}

type billingInstrumentDto struct {
	Id             string `json:"id,omitempty"`
	ResourceGroup  string `json:"resourceGroup"`
	SubscriptionId string `json:"subscriptionId"`
}

type billingPolicyDto struct {
	Id                string               `json:"id"`
	Name              string               `json:"name"`
	TenantType        string               `json:"type"`
	Status            string               `json:"status"`
	Location          string               `json:"location"`
	BillingInstrument billingInstrumentDto `json:"billingInstrument"`
	CreatedOn         string               `json:"createdOn"`
	CreatedBy         principalDto         `json:"createdBy"`
	LastModifiedOn    string               `json:"lastModifiedOn"`
	LastModifiedBy    principalDto         `json:"lastModifiedBy"`
}

type billingPolicyArrayDto struct {
	Value []billingPolicyDto `json:"value"`
}

type billingPolicyUpdateDto struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type principalDto struct {
	Id            string `json:"id"`
	PrincipalType string `json:"type"`
}

type billingPolicyEnvironmentsArrayDto struct {
	EnvironmentIds []string `json:"environmentIds"`
}

type billingPolicyEnvironmentsDto struct {
	BillingPolicyId string `json:"billingPolicyId"`
	EnvironmentId   string `json:"environmentId"`
}

type billingPolicyEnvironmentsArrayResponseDto struct {
	Value []billingPolicyEnvironmentsDto `json:"value"`
}
