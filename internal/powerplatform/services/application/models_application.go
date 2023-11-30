package powerplatform

type ApplicationArrayDto struct {
	Value []ApplicationDto `json:"value"`
}

type ApplicationDto struct {
	ApplicationId         string `json:"applicationId"`
	Name                  string `json:"applicationName"`
	UniqueName            string `json:"uniqueName"`
	Version               string `json:"version"`
	Description           string `json:"localizedDescription"`
	PublisherId           string `json:"publisherId"`
	PublisherName         string `json:"publisherName"`
	LearnMoreUrl          string `json:"learnMoreUrl"`
	State                 string `json:"state"`
	ApplicationVisibility string `json:"applicationVisibility"`
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

type ApplicationLifecycleDto struct {
	Id                 string                             `json:"id"`
	Links              ApplicationLifecycleLinksDto       `json:"links"`
	State              ApplicationLifecycleStateDto       `json:"state"`
	Type               ApplicationLifecycleStateDto       `json:"type"`
	CreatedDateTime    string                             `json:"createdDateTime"`
	LastActionDateTime string                             `json:"lastActionDateTime"`
	RequestedBy        ApplicationLifecycleRequestedByDto `json:"requestedBy"`
	Stages             []ApplicationLifecycleStageDto     `json:"stages"`
}

type ApplicationLifecycleStageDto struct {
	Id                  string                       `json:"id"`
	Name                string                       `json:"name"`
	State               ApplicationLifecycleStateDto `json:"state"`
	FirstActionDateTime string                       `json:"firstActionDateTime"`
	LastActionDateTime  string                       `json:"lastActionDateTime"`
}

type ApplicationLifecycleLinksDto struct {
	Self        ApplicationLifecycleLinkDto `json:"self"`
	Environment ApplicationLifecycleLinkDto `json:"environment"`
}

type ApplicationLifecycleLinkDto struct {
	Path string `json:"path"`
}

type ApplicationLifecycleStateDto struct {
	Id string `json:"id"`
}

type ApplicationLifecycleRequestedByDto struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}
