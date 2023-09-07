package powerplatform_bapi

type ApplicationUserDto struct {
	Name       string                       `json:"name"`
	Id         string                       `json:"id"`
	Type       string                       `json:"type"`
	Properties ApplicationUserPropertiesDto `json:"properties"`
}

type ApplicationUserPropertiesDto struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Tier        string `json:"tier"`
	Publisher   string `json:"publisher"`
	Unblockable bool
}

type ApplicationUserDtoArray struct {
	Value []ApplicationUserDto `json:"value"`
}

type UnblockableApplicationUserDto struct {
	Id       string                                `json:"id"`
	Metadata UnblockableApplicationUserMetadataDto `json:"metadata"`
}

type UnblockableApplicationUserMetadataDto struct {
	Unblockable bool `json:"unblockable"`
}

type VirtualApplicationUserDto struct {
	Id       string                            `json:"id"`
	Metadata VirtualApplicationUserMetadataDto `json:"metadata"`
}

type VirtualApplicationUserMetadataDto struct {
	VirtualApplicationUser bool   `json:"virtualApplicationUser"`
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	DisplayName            string `json:"displayName"`
}
