// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment

type OperationLifecycleDto struct {
	Id                 string                           `json:"id"`
	Links              OperationLifecycleLinksDto       `json:"links"`
	State              OperationLifecycleStateDto       `json:"state"`
	Type               OperationLifecycleStateDto       `json:"type"`
	CreatedDateTime    string                           `json:"createdDateTime"`
	LastActionDateTime string                           `json:"lastActionDateTime"`
	RequestedBy        OperationLifecycleRequestedByDto `json:"requestedBy"`
	Stages             []OperationLifecycleStageDto     `json:"stages"`
}

type OperationLifecycleStageDto struct {
	Id                  string                     `json:"id"`
	Name                string                     `json:"name"`
	State               OperationLifecycleStateDto `json:"state"`
	FirstActionDateTime string                     `json:"firstActionDateTime"`
	LastActionDateTime  string                     `json:"lastActionDateTime"`
}

type OperationLifecycleLinksDto struct {
	Self        OperationLifecycleLinkDto `json:"self"`
	Environment OperationLifecycleLinkDto `json:"environment"`
}

type OperationLifecycleLinkDto struct {
	Path string `json:"path"`
}

type OperationLifecycleStateDto struct {
	Id string `json:"id"`
}

type OperationLifecycleRequestedByDto struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}

type OperationLifecycleCreatedDto struct {
	Name       string                                 `json:"name"`
	Properties OperationLifecycleCreatedPropertiesDto `json:"properties"`
}

type OperationLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}
