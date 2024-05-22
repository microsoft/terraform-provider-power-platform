// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

type ConnectionDtoArray struct {
	Value []ConnectionDto `json:"value"`
}

type ConnectionDto struct {
	Name       string        `json:"name"`
	Id         string        `json:"id"`
	Type       string        `json:"type"`
	Properties PropertiesDto `json:"properties"`
}

type PropertiesDto struct {
	ApiId       string      `json:"apiId"`
	DisplayName string      `json:"displayName"`
	IconUri     string      `json:"iconUri"`
	Statuses    []StatusDto `json:"statuses"`
	//ConnectionParametersSet map[string]interface{} `json:"connectionParametersSet,omitempty"`
	KeywordsRemaining int             `json:"keywordsRemaining"`
	CreatedBy         CreatedByDto    `json:"createdBy"`
	CreatedTime       string          `json:"createdTime"`
	LastModifiedTime  string          `json:"lastModifiedTime"`
	ExpirationTime    string          `json:"expirationTime"`
	TestLinks         []TestLinkDto   `json:"testLinks"`
	Environment       EnvironmentDto  `json:"environment"`
	Permissions       []PermissionDto `json:"permissions"`
	AccountName       string          `json:"accountName"`
	AllowSharing      bool            `json:"allowSharing"`
}

type StatusDto struct {
	Status string `json:"status"`
}

type CreatedByDto struct {
	Id                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Email             string `json:"email"`
	Type              string `json:"type"`
	TenantId          string `json:"tenantId"`
	UserPrincipalName string `json:"userPrincipalName"`
}

type TestLinkDto struct {
	RequestURI string `json:"requestUri"`
	Method     string `json:"method"`
}

type EnvironmentDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PermissionDto struct {
	Name       string                  `json:"name"`
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties PermissionPropertiesDto `json:"properties"`
}

type PermissionPropertiesDto struct {
	RoleName                string       `json:"roleName"`
	Principal               PrincipalDto `json:"principal"`
	NotifyShareTargetOption string       `json:"notifyShareTargetOption"`
	InviteGuestToTenant     bool         `json:"inviteGuestToTenant"`
}

type PrincipalDto struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	TenantId string `json:"tenantId"`
}

type ConnectionToCreateDto struct {
	Properties ConnectionToCreatePropertiesDto `json:"properties"`
}

type ConnectionToCreatePropertiesDto struct {
	DisplayName          string                           `json:"displayName"`
	ConnectionParameters map[string]interface{}           `json:"connectionParametersSet,omitempty"`
	Environment          ConnectionToCreateEnvironmentDto `json:"environment"`
}

type ConnectionToCreateEnvironmentDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
