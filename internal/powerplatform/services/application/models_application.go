package powerplatform

type ApplicationArrayDto struct {
	Value []ApplicationDto `json:"value"`
}

type ApplicationDto struct {
	ApplicationId string `json:"applicationId"`
	Name          string `json:"applicationName"`
	UniqueName    string `json:"uniqueName"`
}

type ApplicationDeleteDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ApplicationCreateDto struct {
	Location string `json:"location"`
}

type ApplicationLifecycleCreatedDto struct {
	Name       string                                   `json:"name"`
	Properties ApplicationLifecycleCreatedPropertiesDto `json:"properties"`
}

type ApplicationLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}

type ApplicationBapi struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ApplicationPropertiesBapi struct {
	TenantID    string `json:"tenantId"`
	DisplayName string `json:"displayName"`
}

type LinkedApplicationMetadataBapi struct {
	Version string `json:"version"`
}

type StatesApplicationBapi struct {
	Management StatesManagementApplicationBapi `json:"management"`
}

type StatesManagementApplicationBapi struct {
	Id string `json:"id"`
}

type ApplicationDtoArray struct {
	Value []ApplicationDto `json:"value"`
}

type ApplicationCreateBapi struct {
	Location string `json:"location"`

	Properties ApplicationPropertiesBapi `json:"properties"`
}

type ApplicationCreatePropertiesBapi struct {
	DisplayName string `json:"displayName"`
}

type ApplicationCreateLinkApplicationMetadataBapi struct {
}
