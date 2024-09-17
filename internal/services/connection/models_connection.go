// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

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
	ApiId                   string                 `json:"apiId"`
	DisplayName             string                 `json:"displayName"`
	IconUri                 string                 `json:"iconUri"`
	Statuses                []StatusDto            `json:"statuses"`
	ConnectionParametersSet map[string]interface{} `json:"connectionParametersSet,omitempty"`
	ConnectionParameters    map[string]interface{} `json:"connectionParameters,omitempty"`
	KeywordsRemaining       int                    `json:"keywordsRemaining"`
	CreatedBy               CreatedByDto           `json:"createdBy"`
	CreatedTime             string                 `json:"createdTime"`
	LastModifiedTime        string                 `json:"lastModifiedTime"`
	ExpirationTime          string                 `json:"expirationTime"`
	TestLinks               []TestLinkDto          `json:"testLinks"`
	Environment             EnvironmentDto         `json:"environment"`
	Permissions             []PermissionDto        `json:"permissions"`
	AccountName             string                 `json:"accountName"`
	AllowSharing            bool                   `json:"allowSharing"`
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
	NotifyShareTargetOption string       `json:"NotifyShareTargetOption"`
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
	DisplayName             string                           `json:"displayName"`
	ConnectionParametersSet map[string]interface{}           `json:"connectionParametersSet,omitempty"`
	ConnectionParameters    map[string]interface{}           `json:"connectionParameters,omitempty"`
	Environment             ConnectionToCreateEnvironmentDto `json:"environment"`
}

type ConnectionToCreateEnvironmentDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ConnectorDefinition struct {
	Name       string `json:"name"`
	Id         string `json:"id"`
	Type       string `json:"type"`
	Properties struct {
		DisplayName             string          `json:"displayName"`
		IconUri                 string          `json:"iconUri"`
		IconBrandColor          string          `json:"iconBrandColor"`
		ApiEnvironment          string          `json:"apiEnvironment"`
		IsCustomApi             bool            `json:"isCustomApi"`
		BlobUrisAreProxied      bool            `json:"blobUrisAreProxied"`
		ConnectionParameters    *map[string]any `json:"connectionParameters,omitempty"`
		ConnectionParameterSets *map[string]any `json:"connectionParameterSets,omitempty"`
		Swagger                 *map[string]any `json:"swagger,omitempty"`
	} `json:"properties"`
}

type ShareConnectionRequestDto struct {
	Put    []ShareConnectionRequestPutDto    `json:"put"`
	Delete []ShareConnectionRequestDeleteDto `json:"delete"`
}

type ShareConnectionRequestPutDto struct {
	Properties ShareConnectionRequestPutPropertiesDto `json:"properties"`
}

type ShareConnectionRequestPutPropertiesDto struct {
	RoleName                string                                          `json:"roleName"`
	Capabilities            []any                                           `json:"capabilities"`
	Principal               ShareConnectionRequestPutPropertiesPrincipalDto `json:"principal"`
	NotifyShareTargetOption string                                          `json:"NotifyShareTargetOption"`
}

type ShareConnectionRequestPutPropertiesPrincipalDto struct {
	Id       string  `json:"id"`
	Type     string  `json:"type"`
	TenantId *string `json:"tenantId"`
}

type ShareConnectionRequestDeleteDto struct {
	Id string `json:"id"`
}

type ShareConnectionResponseArrayDto struct {
	Value []ShareConnectionResponseDto `json:"value"`
}

type ShareConnectionResponseDto struct {
	Name       string                               `json:"name"`
	Id         string                               `json:"id"`
	Type       string                               `json:"type"`
	Properties ShareConnectionResponsePropertiesDto `json:"properties"`
}

type ShareConnectionResponsePropertiesDto struct {
	RoleName string `json:"roleName"`
	//Principal               ShareConnectionResponsePropertiesPrincipalDto `json:"princpal"`
	Principal               map[string]interface{} `json:"principal"`
	NotifyShareTargetOption string                 `json:"NotifyShareTargetOption"`
	InviteGuestToTenant     bool                   `json:"inviteGuestToTenant"`
}

type ShareConnectionResponsePropertiesPrincipalDto struct {
	Id                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Type              string `json:"type"`
	TenantId          string `json:"tenantId"`
	Email             string `json:"email"`
	PreferredLanguage string `json:"preferredLanguage"`
}
