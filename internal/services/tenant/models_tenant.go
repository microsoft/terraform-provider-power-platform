// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant

type TenantDto struct {
	TenantId                         string `json:"tenantId"`
	State                            string `json:"state"`
	Location                         string `json:"location"`
	AadCountryGeo                    string `json:"aadCountryGeo"`
	DataStorageGeo                   string `json:"dataStorageGeo"`
	DefaultEnvironmentGeo            string `json:"defaultEnvironmentGeo"`
	AadDataBoundary                  string `json:"aadDataBoundary"`
	FedRAMPHighCertificationRequired bool   `json:"fedRAMPHighCertificationRequired"`
}
