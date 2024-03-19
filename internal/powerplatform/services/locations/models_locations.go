package powerplatform

type LocationsDto struct {
	Value []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Name       string `json:"name"`
		Properties struct {
			DisplayName                            string   `json:"displayName"`
			Code                                   string   `json:"code"`
			IsDefault                              bool     `json:"isDefault"`
			IsDisabled                             bool     `json:"isDisabled"`
			CanProvisionDatabase                   bool     `json:"canProvisionDatabase"`
			CanProvisionCustomerEngagementDatabase bool     `json:"canProvisionCustomerEngagementDatabase"`
			AzureRegions                           []string `json:"azureRegions"`
		} `json:"properties"`
	} `json:"value"`
}
