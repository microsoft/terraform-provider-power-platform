// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages

type LanguagesDto struct {
	Value []struct {
		Name       string `json:"name"`
		ID         string `json:"id"`
		Type       string `json:"type"`
		Properties struct {
			LocaleID        int64  `json:"localeId"`
			LocalizedName   string `json:"localizedName"`
			DisplayName     string `json:"displayName"`
			IsTenantDefault bool   `json:"isTenantDefault"`
		} `json:"properties"`
	} `json:"value"`
}
