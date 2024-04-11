// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_mocks

import (
	"net/http"

	"github.com/jarcoal/httpmock"
)

func ActivateEnvironmentHttpMocks() {
	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/transactioncurrencies\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"isocurrencycode": "PLN"
					}]}`), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://([\d-]+)\.crm4\.dynamics\.com/api/data/v9\.2/organizations\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
				"value": [
					{
						"_basecurrencyid_value": "xyz"
					}]}`), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/locations/tests/datasource/Validate_Read/get_locations.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentLanguages?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/languages/tests/datasource/Validate_Read/get_languages.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/europe/environmentCurrencies?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/currencies/tests/datasource/Validate_Read/get_currencies.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

}
