// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package currencies

type currenciesDto struct {
	Value []currenciesArrayDto `json:"value"`
}

type currenciesArrayDto struct {
	Name       string                  `json:"name"`
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties currenciesPropertiesDto `json:"properties"`
}

type currenciesPropertiesDto struct {
	Code            string `json:"code"`
	Symbol          string `json:"symbol"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}
