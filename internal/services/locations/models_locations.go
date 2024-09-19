// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package locations

type locationDto struct {
	Value []locationsArrayDto `json:"value"`
}

type locationsArrayDto struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Name       string             `json:"name"`
	Properties locationProperties `json:"properties"`
}

type locationProperties struct {
	DisplayName                            string   `json:"displayName"`
	Code                                   string   `json:"code"`
	IsDefault                              bool     `json:"isDefault"`
	IsDisabled                             bool     `json:"isDisabled"`
	CanProvisionDatabase                   bool     `json:"canProvisionDatabase"`
	CanProvisionCustomerEngagementDatabase bool     `json:"canProvisionCustomerEngagementDatabase"`
	AzureRegions                           []string `json:"azureRegions"`
}
