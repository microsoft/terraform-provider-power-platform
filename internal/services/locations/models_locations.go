// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package locations

type Dto struct {
	Value []ArrayDto `json:"value"`
}

type ArrayDto struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Name       string     `json:"name"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	DisplayName                            string   `json:"displayName"`
	Code                                   string   `json:"code"`
	IsDefault                              bool     `json:"isDefault"`
	IsDisabled                             bool     `json:"isDisabled"`
	CanProvisionDatabase                   bool     `json:"canProvisionDatabase"`
	CanProvisionCustomerEngagementDatabase bool     `json:"canProvisionCustomerEngagementDatabase"`
	AzureRegions                           []string `json:"azureRegions"`
}
