package powerplatform

type BillingPolicyCreateDto struct {
	Location            string                       `json:"location"`
	Name                string                       `json:"name"`
	TenantType          string                       `json:"type"`
	BillingInstrument   BillingInstrumentDto         `json:"billingInstrument"`
	PowerAppsPolicy     PolicyDto                    `json:"powerAppsPolicy"`
	PowerAutomatePolicy PowerAutomatePolicyCreateDto `json:"powerAutomatePolicy"`
	StoragePolicy       PolicyDto                    `json:"storagePolicy"`
}

type BillingInstrumentDto struct {
	Id               string `json:"id"`
	Location         string `json:"location"`
	ResourceGroup    string `json:"resourceGroup"`
	SubscriptionId   string `json:"subscriptionId"`
	SubscriptionName string `json:"subscriptionName"`
	//Tags             []string `json:"tags,omitempty"`
}

type PolicyDto struct {
	PayAsYouGoState string `json:"payAsYouGoState"`
}

type PowerAutomatePolicyCreateDto struct {
	PayAsYouGoState string `json:"payAsYouGoState"`
}

type BillingPolicyDto struct {
	Id                    string                 `json:"id"`
	Name                  string                 `json:"name"`
	TenantType            string                 `json:"type"`
	Status                string                 `json:"status"`
	Location              string                 `json:"location"`
	PowerAutomatePolicy   PowerAutomatePolicyDto `json:"powerAutomatePolicy"`
	PowerAppsPolicy       PolicyDto              `json:"powerAppsPolicy"`
	StoragePolicy         PolicyDto              `json:"storagePolicy"`
	PowerPlatformRequests PolicyDto              `json:"powerPlatformRequests"`
	PowerPagesPolicy      PolicyDto              `json:"powerPagesPolicy"`
	PowerVirtualAgent     PolicyDto              `json:"powerVirtualAgent"`
	BillingInstrument     BillingInstrumentDto   `json:"billingInstrument"`
	CreatedOn             string                 `json:"createdOn"`
	CreatedBy             PrincipalDto           `json:"createdBy"`
	LastModifiedOn        string                 `json:"lastModifiedOn"`
	LastModifiedBy        PrincipalDto           `json:"lastModifiedBy"`
}

type PowerAutomatePolicyDto struct {
	CloudFlowRunsPayAsYouGoState             string `json:"cloudFlowRunsPayAsYouGoState"`
	DesktopFlowUnattendedRunsPayAsYouGoState string `json:"desktopFlowUnattendedRunsPayAsYouGoState"`
	DesktopFlowAttendedRunsPayAsYouGoState   string `json:"desktopFlowAttendedRunsPayAsYouGoState"`
}

type PrincipalDto struct {
	Id            string `json:"id"`
	PrincipalType string `json:"type"`
}
