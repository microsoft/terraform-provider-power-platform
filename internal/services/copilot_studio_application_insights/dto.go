// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights

type copilotStudioAppInsightsDto struct {
	AppInsightsConnectionString string `json:"appInsightsConnectionString"`
	IncludeSensitiveInformation bool   `json:"includeSensitiveInformation"`
	IncludeActivities           bool   `json:"includeActivities"`
	IncludeActions              bool   `json:"includeActions"`
	NetworkIsolation            string `json:"networkIsolation"`
	Errors                      string `json:"errors"`
}

type environmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties environmentIdPropertiesDto `json:"properties"`
}

type environmentIdPropertiesDto struct {
	RuntimeEndpoints runtimeEndpointsDto `json:"runtimeEndpoints"`
}

type runtimeEndpointsDto struct {
	PowerVirtualAgents string `json:"powerVirtualAgents"`
}
