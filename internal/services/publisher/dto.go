// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher

type publisherDto struct {
	Id                             string `json:"publisherid"`
	FriendlyName                   string `json:"friendlyname"`
	UniqueName                     string `json:"uniquename"`
	CustomizationPrefix            string `json:"customizationprefix"`
	CustomizationOptionValuePrefix int64  `json:"customizationoptionvalueprefix"`
	Description                    string `json:"description"`
	EmailAddress                   string `json:"emailaddress"`
	SupportingWebsiteURL           string `json:"supportingwebsiteurl"`
	IsReadOnly                     bool   `json:"isreadonly"`

	Address1AddressId          string   `json:"address1_addressid"`
	Address1AddressTypeCode    *int64   `json:"address1_addresstypecode"`
	Address1City               string   `json:"address1_city"`
	Address1Country            string   `json:"address1_country"`
	Address1County             string   `json:"address1_county"`
	Address1Fax                string   `json:"address1_fax"`
	Address1Latitude           *float64 `json:"address1_latitude"`
	Address1Line1              string   `json:"address1_line1"`
	Address1Line2              string   `json:"address1_line2"`
	Address1Line3              string   `json:"address1_line3"`
	Address1Longitude          *float64 `json:"address1_longitude"`
	Address1Name               string   `json:"address1_name"`
	Address1PostalCode         string   `json:"address1_postalcode"`
	Address1PostOfficeBox      string   `json:"address1_postofficebox"`
	Address1ShippingMethodCode *int64   `json:"address1_shippingmethodcode"`
	Address1StateOrProvince    string   `json:"address1_stateorprovince"`
	Address1Telephone1         string   `json:"address1_telephone1"`
	Address1Telephone2         string   `json:"address1_telephone2"`
	Address1Telephone3         string   `json:"address1_telephone3"`
	Address1UpsZone            string   `json:"address1_upszone"`
	Address1UtcOffset          *int64   `json:"address1_utcoffset"`

	Address2AddressId          string   `json:"address2_addressid"`
	Address2AddressTypeCode    *int64   `json:"address2_addresstypecode"`
	Address2City               string   `json:"address2_city"`
	Address2Country            string   `json:"address2_country"`
	Address2County             string   `json:"address2_county"`
	Address2Fax                string   `json:"address2_fax"`
	Address2Latitude           *float64 `json:"address2_latitude"`
	Address2Line1              string   `json:"address2_line1"`
	Address2Line2              string   `json:"address2_line2"`
	Address2Line3              string   `json:"address2_line3"`
	Address2Longitude          *float64 `json:"address2_longitude"`
	Address2Name               string   `json:"address2_name"`
	Address2PostalCode         string   `json:"address2_postalcode"`
	Address2PostOfficeBox      string   `json:"address2_postofficebox"`
	Address2ShippingMethodCode *int64   `json:"address2_shippingmethodcode"`
	Address2StateOrProvince    string   `json:"address2_stateorprovince"`
	Address2Telephone1         string   `json:"address2_telephone1"`
	Address2Telephone2         string   `json:"address2_telephone2"`
	Address2Telephone3         string   `json:"address2_telephone3"`
	Address2UpsZone            string   `json:"address2_upszone"`
	Address2UtcOffset          *int64   `json:"address2_utcoffset"`
}

type publishersDto struct {
	Value []publisherDto `json:"value"`
}
