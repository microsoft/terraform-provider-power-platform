// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package currencies

type Dto struct {
	Value []ArrayDto `json:"value"`
}

type ArrayDto struct {
	Name       string        `json:"name"`
	ID         string        `json:"id"`
	Type       string        `json:"type"`
	Properties PropertiesDto `json:"properties"`
}

type PropertiesDto struct {
	Code            string `json:"code"`
	Symbol          string `json:"symbol"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}
