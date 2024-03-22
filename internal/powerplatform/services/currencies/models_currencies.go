// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

type CurrenciesDto struct {
	Value []struct {
		Name       string `json:"name"`
		ID         string `json:"id"`
		Type       string `json:"type"`
		Properties struct {
			Code            string `json:"code"`
			Symbol          string `json:"symbol"`
			IsTenantDefault bool   `json:"isTenantDefault"`
		} `json:"properties"`
	} `json:"value"`
}
