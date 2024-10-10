// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

type createEnvironmentGroupRuleSetDto struct {
	Parameters []environmentGroupRuleSetParameterDto `json:"parameters"`
}

type environmentGroupRuleSetDto struct {
	Value []environmentGroupRuleSetValueSetDto `json:"value"`
}

type environmentGroupRuleSetValueSetDto struct {
	Parameters        []environmentGroupRuleSetParameterDto       `json:"parameters"`
	Id                string                                      `json:"id"`
	LastModified      string                                      `json:"lastModified"`
	EnvironmentFilter environmentGroupRuleSetEnvironmentFilterDto `json:"environmentFilter"`
}

type environmentGroupRuleSetEnvironmentFilterDto struct {
	Type  string   `json:"type"`
	Value []string `json:"values"`
}

type environmentGroupRuleSetValueDto struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type environmentGroupRuleSetParameterDto struct {
	Type         string                            `json:"type"`
	ResourceType string                            `json:"resourceType"`
	Value        []environmentGroupRuleSetValueDto `json:"value"`
}
