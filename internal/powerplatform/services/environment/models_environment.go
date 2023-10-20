package powerplatform

import "time"

var (
	//https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies?api-version=2023-06-01
	EnvironmentCurrencyCodes = []string{"DJF", "ZAR", "ETB", "AED", "BHD", "DZD", "EGP", "IQD", "JOD", "KWD",
		"LBP", "LYD", "MAD", "OMR", "QAR", "SAR", "SYP", "TND", "YER", "CLP",
		"INR", "AZN", "RUB", "BYN", "BGN", "NGN", "BDT", "CNY", "EUR", "BAM",
		"USD", "CZK", "GBP", "DKK", "CHF", "MVR", "BTN", "XCD", "AUD", "BZD",
		"CAD", "HKD", "IDR", "JMD", "MYR", "NZD", "PHP", "SGD", "TTD", "XDR",
		"ARS", "BOB", "COP", "CRC", "CUP", "DOP", "GTQ", "HNL", "MXN", "NIO",
		"PAB", "PEN", "PYG", "UYU", "VES", "IRR", "XOF", "CDF", "XAF", "HTG",
		"ILS", "HRK", "HUF", "AMD", "ISK", "JPY", "GEL", "KZT", "KHR", "KRW",
		"KGS", "LAK", "MKD", "MNT", "BND", "MMK", "NOK", "NPR", "PKR", "PLN",
		"AFN", "BRL", "MDL", "RON", "RWF", "SEK", "LKR", "SOS", "ALL", "RSD",
		"KES", "TJS", "THB", "ERN", "TMT", "BWP", "TRY", "UAH", "UZS", "VND",
		"MOP", "TWD"}

	//https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01
	EnvironmentLocations = []string{"unitedstates",
		"europe", "asia", "australia", "india", "japan", "canada",
		"unitedkingdom", "southamerica", "france", "unitedarabemirates", "germany",
		"switzerland", "norway", "korea", "southafrica"}

	//https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages?api-version=2023-06-01
	EnvironmentLanguages = []int64{1033, 1025, 1026, 1069, 1027, 3076, 2052, 1028, 1050, 1029, 1030, 1043, 1061,
		1035, 1036, 1110, 1031, 1032, 1037, 1081, 1038, 1040, 1041, 1087, 1042, 1062,
		1063, 1044, 1045, 1046, 2070, 1048, 1049, 2074, 1051, 1060, 3082, 1053, 1054,
		1055, 1058, 1066, 3098, 1086, 1057}

	EnvironmentTypes = []string{"Sandbox", "Production", "Trial", "Developer"}
)

type EnvironmentDto struct {
	Id         string                   `json:"id"`
	Type       string                   `json:"type"`
	Location   string                   `json:"location"`
	Name       string                   `json:"name"`
	Properties EnvironmentPropertiesDto `json:"properties"`
}

type EnvironmentPropertiesDto struct {
	DatabaseType              string                       `json:"databaseType"`
	DisplayName               string                       `json:"displayName"`
	EnvironmentSku            string                       `json:"environmentSku"`
	LinkedAppMetadata         *LinkedAppMetadataDto        `json:"linkedAppMetadata,omitempty"`
	LinkedEnvironmentMetadata LinkedEnvironmentMetadataDto `json:"linkedEnvironmentMetadata"`
	States                    StatesEnvironmentDto         `json:"states"`
	TenantID                  string                       `json:"tenantId"`
}

type LinkedEnvironmentMetadataDto struct {
	DomainName       string                            `json:"domainName,omitempty"`
	InstanceURL      string                            `json:"instanceUrl"`
	BaseLanguage     int                               `json:"baseLanguage"`
	SecurityGroupId  string                            `json:"securityGroupId"`
	ResourceId       string                            `json:"resourceId"`
	Version          string                            `json:"version"`
	Templates        []string                          `json:"template,omitempty"`
	TemplateMetadata EnvironmentCreateTemplateMetadata `json:"templateMetadata,omitempty"`
}

type LinkedAppMetadataDto struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type StatesEnvironmentDto struct {
	Management StatesManagementEnvironmentDto `json:"management"`
}

type StatesManagementEnvironmentDto struct {
	Id string `json:"id"`
}

type EnvironmentDtoArray struct {
	Value []EnvironmentDto `json:"value"`
}

type EnvironmentCreateDto struct {
	Location   string                         `json:"location"`
	Properties EnvironmentCreatePropertiesDto `json:"properties"`
}

type EnvironmentCreatePropertiesDto struct {
	BillingPolicy             string                                      `json:"billingPolicy,omitempty"`
	DataBaseType              string                                      `json:"databaseType,omitempty"`
	DisplayName               string                                      `json:"displayName"`
	EnvironmentSku            string                                      `json:"environmentSku"`
	LinkedEnvironmentMetadata EnvironmentCreateLinkEnvironmentMetadataDto `json:"linkedEnvironmentMetadata"`
}

type EnvironmentCreateLinkEnvironmentMetadataDto struct {
	BaseLanguage     int                               `json:"baseLanguage"`
	DomainName       string                            `json:"domainName,omitempty"`
	Currency         EnvironmentCreateCurrency         `json:"currency"`
	SecurityGroupId  string                            `json:"securityGroupId,omitempty"`
	Templates        []string                          `json:"templates,omitempty"`
	TemplateMetadata EnvironmentCreateTemplateMetadata `json:"templateMetadata,omitempty"`
}
type EnvironmentCreateCurrency struct {
	Code string `json:"code"`
}

type EnvironmentCreateTemplateMetadata struct {
	PostProvisioningPackages []EnvironmentCreatePostProvisioningPackages `json:"PostProvisioningPackages,omitempty"`
}

type EnvironmentCreatePostProvisioningPackages struct {
	ApplicationUniqueName string `json:"applicationUniqueName,omitempty"`
	Parameters            string `json:"parameters,omitempty"`
}

type EnvironmentCreateLinkedAppMetadataDto struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type EnvironmentDeleteDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type EnvironmentLifecycleDto struct {
	Id                 string                             `json:"id"`
	Links              EnvironmentLifecycleLinksDto       `json:"links"`
	State              EnvironmentLifecycleStateDto       `json:"state"`
	Type               EnvironmentLifecycleStateDto       `json:"type"`
	CreatedDateTime    string                             `json:"createdDateTime"`
	LastActionDateTime string                             `json:"lastActionDateTime"`
	RequestedBy        EnvironmentLifecycleRequestedByDto `json:"requestedBy"`
	Stages             []EnvironmentLifecycleStageDto     `json:"stages"`
}

type EnvironmentLifecycleStageDto struct {
	Id                  string                       `json:"id"`
	Name                string                       `json:"name"`
	State               EnvironmentLifecycleStateDto `json:"state"`
	FirstActionDateTime string                       `json:"firstActionDateTime"`
	LastActionDateTime  string                       `json:"lastActionDateTime"`
}

type EnvironmentLifecycleLinksDto struct {
	Self        EnvironmentLifecycleLinkDto `json:"self"`
	Environment EnvironmentLifecycleLinkDto `json:"environment"`
}

type EnvironmentLifecycleLinkDto struct {
	Path string `json:"path"`
}

type EnvironmentLifecycleStateDto struct {
	Id string `json:"id"`
}

type EnvironmentLifecycleRequestedByDto struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	Type        string `json:"type"`
}

type EnvironmentLifecycleCreatedDto struct {
	Name       string                                   `json:"name"`
	Properties EnvironmentLifecycleCreatedPropertiesDto `json:"properties"`
}

type EnvironmentLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}

type OrganizationSettingsArrayDto struct {
	Value []OrganizationSettingsDto `json:"value"`
}

type OrganizationSettingsDto struct {
	ODataEtag      string    `json:"@odata.etag"`
	CreatedOn      time.Time `json:"createdon"`
	BaseCurrencyId string    `json:"_basecurrencyid_value"`
}

type TransactionCurrencyDto struct {
	OrganizationValue     string  `json:"_organizationid_value"`
	CurrencyName          string  `json:"currencyname"`
	CurrencySymbol        string  `json:"currencysymbol"`
	IsoCurrencyCode       string  `json:"isocurrencycode"`
	CreatedOn             string  `json:"createdon"`
	CurrencyPrecision     int     `json:"currencyprecision"`
	ExchangeRate          float32 `json:"exchangerate"`
	TransactionCurrencyId string  `json:"transactioncurrencyid"`
}

type TransactionCurrencyArrayDto struct {
	Value []TransactionCurrencyDto `json:"value"`
}
