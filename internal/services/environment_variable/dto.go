// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environmentvariable

type environmentVariableDefinitionArrayDto struct {
	Value []environmentVariableDefinitionDto `json:"value"`
}

type environmentVariableDefinitionDto struct {
	EnvironmentVariableDefinitionId string `json:"environmentvariabledefinitionid"`
	SchemaName                      string `json:"schemaname"`
	DisplayName                     string `json:"displayname"`
	Description                     string `json:"description"`
	DefaultValue                    string `json:"defaultvalue"`
	Type                            int64  `json:"type"`
	ValueSchema                     string `json:"valueschema"`
	SecretStore                     int64  `json:"secretstore"`
}

type environmentVariableValueArrayDto struct {
	Value []environmentVariableValueDto `json:"value"`
}

type environmentVariableValueDto struct {
	EnvironmentVariableValueId string `json:"environmentvariablevalueid"`
	SchemaName                 string `json:"schemaname"`
	Value                      string `json:"value"`
}

type createEnvironmentVariableValueDto struct {
	SchemaName                            string `json:"schemaname"`
	Value                                 string `json:"value"`
	EnvironmentVariableDefinitionODataRef string `json:"environmentvariabledefinitionid@odata.bind"`
}

type updateEnvironmentVariableValueDto struct {
	Value string `json:"value"`
}

type resolvedEnvironmentVariableDto struct {
	Definition environmentVariableDefinitionDto
	Value      *environmentVariableValueDto
}
