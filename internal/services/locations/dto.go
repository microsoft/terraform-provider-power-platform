// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package locations

type locationDto struct {
	Value                  []locationsArrayDto `json:"value"`
	TenantProvisioningMode string              `json:"tenantProvisioningMode"`
	MacroRegions           []macroRegionDto    `json:"macroRegions"`
}

type macroRegionDto struct {
	MacroRegionId string `json:"macroRegionId"`
	DisplayName   string `json:"displayName"`
	// Pointer because the live API omits the key for some macro regions (for example `europe-uk`).
	DataResidencyNote *string `json:"dataResidencyNote"`
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
