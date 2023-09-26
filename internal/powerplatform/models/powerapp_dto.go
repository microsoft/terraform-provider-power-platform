package powerplatform_models

type PowerAppBapi struct {
	Name       string                 `json:"name"`
	Properties PowerAppPropertiesBapi `json:"properties"`
}

type PowerAppPropertiesBapi struct {
	DisplayName      string                 `json:"displayName"`
	Owner            PowerAppCreatedByDto   `json:"owner"`
	CreatedBy        PowerAppCreatedByDto   `json:"createdBy"`
	LastModifiedBy   PowerAppCreatedByDto   `json:"lastModifiedBy"`
	LastPublishedBy  PowerAppCreatedByDto   `json:"lastPublishedBy"`
	CreatedTime      string                 `json:"createdTime"`
	LastModifiedTime string                 `json:"lastModifiedTime"`
	LastPublishTime  string                 `json:"lastPublishTime"`
	Environment      PowerAppEnvironmentDto `json:"environment"`
}

type PowerAppEnvironmentDto struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type PowerAppCreatedByDto struct {
	DisplayName       string `json:"displayName"`
	Id                string `json:"id"`
	UserPrincipalName string `json:"userPrincipalName"`
}

type PowerAppDtoArray struct {
	Value []PowerAppBapi `json:"value"`
}
