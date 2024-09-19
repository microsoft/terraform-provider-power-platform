// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity

type Dto struct {
	TenantId         string              `json:"tenantId"`
	LicenseModelType string              `json:"licenseModelType"`
	TenantCapacities []TenantCapacityDto `json:"tenantCapacities"`
}

type TenantCapacityDto struct {
	CapacityType  string         `json:"capacityType"`
	CapacityUnits string         `json:"capacityUnits"`
	TotalCapacity float32        `json:"totalCapacity"`
	MaxCapacity   float32        `json:"maxCapacity"`
	Consumption   ConsumptionDto `json:"consumption"`
	Status        string         `json:"status"`
}

type ConsumptionDto struct {
	Actual          float32 `json:"actual"`
	Rated           float32 `json:"rated"`
	ActualUpdatedOn string  `json:"actualUpdatedOn"`
	RatedUpdatedOn  string  `json:"ratedUpdatedOn"`
}
