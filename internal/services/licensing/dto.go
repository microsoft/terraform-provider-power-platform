// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

import "encoding/json"

type billingPolicyCreateDto struct {
	Location          string                `json:"location"`
	Name              string                `json:"name"`
	Status            string                `json:"status"`
	BillingInstrument BillingInstrumentDto  `json:"billingInstrument"`
	PayGoEntitlements []PayGoEntitlementDto `json:"payGoEntitlements,omitempty"`
}

type PayGoEntitlementDto struct {
	EntitlementId   string               `json:"entitlementId"`
	ProductCategory PayGoProductCategory `json:"productCategory,omitempty"`
	PayAsYouGoState bool                 `json:"payAsYouGoState"`
}

// PayGoProductCategory tolerates the numeric product category enum used by some API surfaces so that reading a
// billing policy never fails because of it; numeric values are resolved from the entitlement ID instead.
type PayGoProductCategory string

func (c *PayGoProductCategory) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		*c = ""
		return nil
	}
	*c = PayGoProductCategory(value)
	return nil
}

type BillingInstrumentDto struct {
	Id             string `json:"id,omitempty"`
	ResourceGroup  string `json:"resourceGroup"`
	SubscriptionId string `json:"subscriptionId"`
}

type BillingPolicyDto struct {
	Id                string                `json:"id"`
	Name              string                `json:"name"`
	TenantType        string                `json:"type"`
	Status            string                `json:"status"`
	Location          string                `json:"location"`
	BillingInstrument BillingInstrumentDto  `json:"billingInstrument"`
	CreatedOn         string                `json:"createdOn"`
	CreatedBy         PrincipalDto          `json:"createdBy"`
	LastModifiedOn    string                `json:"lastModifiedOn"`
	LastModifiedBy    PrincipalDto          `json:"lastModifiedBy"`
	PayGoEntitlements []PayGoEntitlementDto `json:"payGoEntitlements"`
}

type BillingPolicyArrayDto struct {
	Value []BillingPolicyDto `json:"value"`
}

type BillingPolicyUpdateDto struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	// The API disables every meter when this list is omitted, empty or null, so updates must always send the full list.
	PayGoEntitlements []PayGoEntitlementDto `json:"payGoEntitlements,omitempty"`
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
