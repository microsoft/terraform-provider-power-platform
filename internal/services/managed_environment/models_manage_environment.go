// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment

type operationLifecycleDto struct {
	Id                 string                           `json:"id"`
	Links              operationLifecycleLinksDto       `json:"links"`
	State              operationLifecycleStateDto       `json:"state"`
	Type               operationLifecycleStateDto       `json:"type"`
	CreatedDateTime    string                           `json:"createdDateTime"`
	LastActionDateTime string                           `json:"lastActionDateTime"`
	RequestedBy        operationLifecycleRequestedByDto `json:"requestedBy"`
	Stages             []operationLifecycleStageDto     `json:"stages"`
}

type operationLifecycleStageDto struct {
	Id                  string                     `json:"id"`
	Name                string                     `json:"name"`
	State               operationLifecycleStateDto `json:"state"`
	FirstActionDateTime string                     `json:"firstActionDateTime"`
	LastActionDateTime  string                     `json:"lastActionDateTime"`
}

type operationLifecycleLinksDto struct {
	Self        operationLifecycleLinkDto `json:"self"`
	Environment operationLifecycleLinkDto `json:"environment"`
}

type operationLifecycleLinkDto struct {
	Path string `json:"path"`
}

type operationLifecycleStateDto struct {
	Id string `json:"id"`
}

type operationLifecycleRequestedByDto struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}

type operationLifecycleCreatedDto struct {
	Name       string                                 `json:"name"`
	Properties operationLifecycleCreatedPropertiesDto `json:"properties"`
}

type operationLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}
