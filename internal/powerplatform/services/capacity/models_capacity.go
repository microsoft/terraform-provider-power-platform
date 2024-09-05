// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

type TenantCapacityDto struct {
	TenantId         string        `json:"tenantId"`
	LicenseModelType string        `json:"licenseModelType"`
	TenantCapacities []CapacityDto `json:"tenantCapacities"`
}

type CapacityDto struct {
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
