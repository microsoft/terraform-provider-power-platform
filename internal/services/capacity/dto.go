// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity

// tenantCapacityDetailsDto represents the main response model for tenant capacity details.
type tenantCapacityDetailsDto struct {
	CapacitySummary      capacitySummaryDto                     `json:"capacitySummary"`
	LegacyModelCapacity  *legacyCapacityModelDto                `json:"legacyModelCapacity,omitempty"`
	LicenseModelType     string                                 `json:"licenseModelType"`
	TemporaryLicenseInfo temporaryLicenseInfoDto                `json:"temporaryLicenseInfo"`
	TenantCapacities     []tenantCapacityAndConsumptionModelDto `json:"tenantCapacities"`
	TenantId             string                                 `json:"tenantId"`
}

// capacitySummaryDto represents the capacity summary information.
type capacitySummaryDto struct {
	FinOpsStatus            string `json:"finOpsStatus"`
	FinOpsStatusMessage     string `json:"finOpsStatusMessage"`
	FinOpsStatusMessageCode string `json:"finOpsStatusMessageCode"`
	Status                  string `json:"status"`
	StatusMessage           string `json:"statusMessage"`
	StatusMessageCode       string `json:"statusMessageCode"`
}

// legacyCapacityModelDto represents legacy capacity model information.
type legacyCapacityModelDto struct {
	CapacityUnits    string  `json:"capacityUnits"`
	TotalCapacity    float64 `json:"totalCapacity"`
	TotalConsumption float64 `json:"totalConsumption"`
}

// temporaryLicenseInfoDto represents temporary license information.
type temporaryLicenseInfoDto struct {
	HasTemporaryLicense        bool   `json:"hasTemporaryLicense"`
	TemporaryLicenseExpiryDate string `json:"temporaryLicenseExpiryDate,omitempty"`
}

// tenantCapacityAndConsumptionModelDto represents tenant capacity and consumption model.
type tenantCapacityAndConsumptionModelDto struct {
	CapacityEntitlements []tenantCapacityEntitlementModelDto `json:"capacityEntitlements"`
	CapacityType         string                              `json:"capacityType"`
	CapacityUnits        string                              `json:"capacityUnits"`
	Consumption          consumptionModelDto                 `json:"consumption"`
	MaxCapacity          float32                             `json:"maxCapacity"`
	OverflowCapacity     []overflowCapacityModelDto          `json:"overflowCapacity"`
	Status               string                              `json:"status"`
	TotalCapacity        float32                             `json:"totalCapacity"`
}

// tenantCapacityEntitlementModelDto represents tenant capacity entitlement model.
type tenantCapacityEntitlementModelDto struct {
	CapacitySubType      string                   `json:"capacitySubType"`
	CapacityType         string                   `json:"capacityType"`
	Licenses             []licenseDetailsModelDto `json:"licenses"`
	MaxNextLifecycleDate string                   `json:"maxNextLifecycleDate,omitempty"`
	TotalCapacity        float32                  `json:"totalCapacity"`
}

// licenseDetailsModelDto represents license details model.
type licenseDetailsModelDto struct {
	CapabilityStatus           string             `json:"capabilityStatus"`
	DisplayName                string             `json:"displayName"`
	EntitlementCode            string             `json:"entitlementCode"`
	IsTemporaryLicense         bool               `json:"isTemporaryLicense"`
	NextLifecycleDate          string             `json:"nextLifecycleDate,omitempty"`
	Paid                       licenseQuantityDto `json:"paid"`
	ServicePlanId              string             `json:"servicePlanId"`
	SkuId                      string             `json:"skuId"`
	TemporaryLicenseExpiryDate string             `json:"temporaryLicenseExpiryDate,omitempty"`
	TotalCapacity              float32            `json:"totalCapacity"`
	Trial                      licenseQuantityDto `json:"trial"`
}

// licenseQuantityDto represents license quantity information.
type licenseQuantityDto struct {
	Enabled   int32 `json:"enabled"`
	Suspended int32 `json:"suspended"`
	Warning   int32 `json:"warning"`
}

// consumptionModelDto represents consumption model information.
type consumptionModelDto struct {
	Actual          float32 `json:"actual"`
	ActualUpdatedOn string  `json:"actualUpdatedOn,omitempty"`
	Rated           float32 `json:"rated"`
	RatedUpdatedOn  string  `json:"ratedUpdatedOn,omitempty"`
}

// overflowCapacityModelDto represents overflow capacity model.
type overflowCapacityModelDto struct {
	CapacityType string  `json:"capacityType"`
	Value        float32 `json:"value"`
}
