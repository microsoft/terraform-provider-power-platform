// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages

type languagesArrayDto struct {
	Value []languageDto `json:"value"`
}
type languageDto struct {
	Name       string                `json:"name"`
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Properties languagePropertiesDto `json:"properties"`
}

type languagePropertiesDto struct {
	LocaleID        int64  `json:"localeId"`
	LocalizedName   string `json:"localizedName"`
	DisplayName     string `json:"displayName"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}
