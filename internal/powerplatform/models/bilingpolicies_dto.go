// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_models

type BillingPolicyDto struct {
	Id                string `json:"id,omitempty"`
	Name              string `json:"name"`
	Type              string `json:"type"`
	Status            string `json:"status"`
	Location          string `json:"location"`
	BillingInstrument struct {
		Id             string `json:"id,omitempty"`
		SubscriptionId string `json:"subscriptionId"`
		ResourceGroup  string `json:"resourceGroup"`
	} `json:"billingInstrument"`
	CreatedOn string `json:"createdOn,omitempty"`
	CreatedBy struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"createdBy,omitempty"`
	LastModifiedOn string `json:"lastModifiedOn,omitempty"`
	LastModifiedBy struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"lastModifiedBy,omitempty"`
}

type BillingPolicyDtoArray struct {
	Value []BillingPolicyDto `json:"value"`
}