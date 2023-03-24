package powerplatform

type Environment struct {
	EnvironmentName                     string `json:"environmentName"`
	DisplayName                         string `json:"displayName"`
	Url                                 string `json:"url"`
	Domain                              string `json:"domain"`
	Location                            string `json:"location"`
	EnvironmentType                     string `json:"environmentType"`
	CommonDataServiceDatabaseType       string `json:"commonDataServiceDatabaseType"`
	OrganizationId                      string `json:"organizationId"`
	SecurityGroupId                     string `json:"securityGroupId"`
	LanguageName                        int    `json:"LanguageName"`
	CurrencyName                        string `json:"CurrencyName"`
	IsCustomControlsInCanvasAppsEnabled bool   `json:"IsCustomControlsInCanvasAppsEnabled"`
}

type EnvironmentCreate struct {
	DisplayName                         string `json:"DisplayName"`
	Location                            string `json:"Location"`
	EnvironmentType                     string `json:"environmentType"`
	LanguageName                        int    `json:"LanguageName"`
	CurrencyName                        string `json:"CurrencyName"`
	IsCustomControlsInCanvasAppsEnabled bool   `json:"IsCustomControlsInCanvasAppsEnabled"`
}

type Solution struct {
	SolutionName    string `json:"SolutionName"`
	SolutionVersion string `json:"SolutionVersion"`
	EnvironmentName string `json:"EnvironmentName"`
	File            string
	SettingsFile    string
	IsManaged       bool   `json:"IsManaged"`
	DisplayName     string `json:"DisplayName"`
}

type SolutionCreate struct {
	SolutionName    string `json:"SolutionName"`
	SolutionContent []byte
	SettingsContent []byte
}

type PackageDeploy struct {
	PackageName     string `json:"PackageName"`
	PackageContent  []byte
	PackageSettings string `json:"PackageSettings"`
	ImportLogs      string `json:"ImportLogs"`
}

type DlpPolicy struct {
	Name                            string               `json:"Name"`
	DisplayName                     string               `json:"DisplayName"`
	CreatedBy                       string               `json:"CreatedBy"`
	CreatedTime                     string               `json:"CreatedTime"`
	LastModifiedBy                  string               `json:"LastModifiedBy"`
	LastModifiedTime                string               `json:"LastModifiedTime"`
	ETag                            string               `json:"Etag"`
	EnvironmentType                 string               `json:"EnvironmentType"`
	DefaultConnectorsClassification string               `json:"DefaultConnectorsClassification"`
	Environments                    []DlpEnvironment     `json:"Environments"`
	ConnectorGroups                 []DlpConnectorGroups `json:"ConnectorGroups"`
}

type DlpEnvironment struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
	Type string `json:"Type"`
}

type DlpConnectorGroups struct {
	Classification string         `json:"Classification"`
	Connectors     []DlpConnector `json:"Connectors"`
}

type DlpConnector struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
	Type string `json:"Type"`
}

type User struct {
	IsApplicationUser bool     `json:"IsApplicationUser"`
	ApplicationId     string   `json:"ApplicationId"`
	IsDisabled        bool     `json:"IsDisabled"`
	UserId            string   `json:"UserId"`
	DomainName        string   `json:"DomainName"`
	FirstName         string   `json:"FirstName"`
	LastName          string   `json:"LastName"`
	AadObjectId       string   `json:"AadObjectId"`
	SecurityRoles     []string `json:"SecurityRoles"`
}

type App struct {
	Name            string `json:"name"`
	DisplayName     string `json:"displayName"`
	EnvironmentName string `json:"EnvironmentName"`
	CreatedTime     string `json:"createdTime"`
}
