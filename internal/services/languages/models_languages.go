// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages

type ArrayDto struct {
	Value []Dto `json:"value"`
}
type Dto struct {
	Name       string        `json:"name"`
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Properties PropertiesDto `json:"properties"`
}

type PropertiesDto struct {
	LocaleID        int64  `json:"localeId"`
	LocalizedName   string `json:"localizedName"`
	DisplayName     string `json:"displayName"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}
