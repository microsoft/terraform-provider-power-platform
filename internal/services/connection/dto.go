// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

type connectionArrayDto struct {
	Value []connectionDto `json:"value"`
}

type connectionDto struct {
	Name       string                  `json:"name"`
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties connectionPropertiesDto `json:"properties"`
}

type connectionPropertiesDto struct {
	ApiId                   string          `json:"apiId"`
	DisplayName             string          `json:"displayName"`
	IconUri                 string          `json:"iconUri"`
	Statuses                []statusDto     `json:"statuses"`
	ConnectionParametersSet map[string]any  `json:"connectionParametersSet,omitempty"`
	ConnectionParameters    map[string]any  `json:"connectionParameters,omitempty"`
	KeywordsRemaining       int             `json:"keywordsRemaining"`
	CreatedBy               createdByDto    `json:"createdBy"`
	CreatedTime             string          `json:"createdTime"`
	LastModifiedTime        string          `json:"lastModifiedTime"`
	ExpirationTime          string          `json:"expirationTime"`
	TestLinks               []testLinkDto   `json:"testLinks"`
	Environment             environmentDto  `json:"environment"`
	Permissions             []permissionDto `json:"permissions"`
	AccountName             string          `json:"accountName"`
	AllowSharing            bool            `json:"allowSharing"`
}

type statusDto struct {
	Status string `json:"status"`
}

type createdByDto struct {
	Id                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Email             string `json:"email"`
	Type              string `json:"type"`
	TenantId          string `json:"tenantId"`
	UserPrincipalName string `json:"userPrincipalName"`
}

type testLinkDto struct {
	RequestURI string `json:"requestUri"`
	Method     string `json:"method"`
}

type environmentDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type permissionDto struct {
	Name       string                  `json:"name"`
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties permissionPropertiesDto `json:"properties"`
}

type permissionPropertiesDto struct {
	RoleName                string       `json:"roleName"`
	Principal               principalDto `json:"principal"`
	NotifyShareTargetOption string       `json:"NotifyShareTargetOption"`
	InviteGuestToTenant     bool         `json:"inviteGuestToTenant"`
}

type principalDto struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	TenantId string `json:"tenantId"`
}

type createDto struct {
	Properties createPropertiesDto `json:"properties"`
}

type createPropertiesDto struct {
	DisplayName             string               `json:"displayName"`
	ConnectionParametersSet map[string]any       `json:"connectionParametersSet,omitempty"`
	ConnectionParameters    map[string]any       `json:"connectionParameters,omitempty"`
	Environment             createEnvironmentDto `json:"environment"`
}

type createEnvironmentDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type shareConnectionRequestDto struct {
	Put    []shareConnectionRequestPutDto    `json:"put"`
	Delete []shareConnectionRequestDeleteDto `json:"delete"`
}

type shareConnectionRequestPutDto struct {
	Properties shareConnectionRequestPutPropertiesDto `json:"properties"`
}

type shareConnectionRequestPutPropertiesDto struct {
	RoleName                string                                          `json:"roleName"`
	Capabilities            []any                                           `json:"capabilities"`
	Principal               shareConnectionRequestPutPropertiesPrincipalDto `json:"principal"`
	NotifyShareTargetOption string                                          `json:"NotifyShareTargetOption"`
}

type shareConnectionRequestPutPropertiesPrincipalDto struct {
	Id       string  `json:"id"`
	Type     string  `json:"type"`
	TenantId *string `json:"tenantId"`
}

type shareConnectionRequestDeleteDto struct {
	Id string `json:"id"`
}

type shareConnectionResponseArrayDto struct {
	Value []shareConnectionResponseDto `json:"value"`
}

type shareConnectionResponseDto struct {
	Name       string                               `json:"name"`
	Id         string                               `json:"id"`
	Type       string                               `json:"type"`
	Properties shareConnectionResponsePropertiesDto `json:"properties"`
}

type shareConnectionResponsePropertiesDto struct {
	RoleName                string         `json:"roleName"`
	Principal               map[string]any `json:"principal"`
	NotifyShareTargetOption string         `json:"NotifyShareTargetOption"`
	InviteGuestToTenant     bool           `json:"inviteGuestToTenant"`
}
