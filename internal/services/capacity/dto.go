// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity

type capacityDto struct {
	TenantId         string              `json:"tenantId"`
	LicenseModelType string              `json:"licenseModelType"`
	TenantCapacities []tenantCapacityDto `json:"tenantCapacities"`
}

type tenantCapacityDto struct {
	CapacityType  string         `json:"capacityType"`
	CapacityUnits string         `json:"capacityUnits"`
	TotalCapacity float32        `json:"totalCapacity"`
	MaxCapacity   float32        `json:"maxCapacity"`
	Consumption   consumptionDto `json:"consumption"`
	Status        string         `json:"status"`
}

type consumptionDto struct {
	Actual          float32 `json:"actual"`
	Rated           float32 `json:"rated"`
	ActualUpdatedOn string  `json:"actualUpdatedOn"`
	RatedUpdatedOn  string  `json:"ratedUpdatedOn"`
}
