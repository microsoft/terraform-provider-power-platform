// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package mocks

import (
	"context"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/provider"
)

func TestName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

func TestsEntraLicesingGroupName() string {
	return "pptestusers"
}

var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerplatform": providerserver.NewProtocol6WithError(provider.NewPowerPlatformProvider(context.Background(), true)()),
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerplatform": providerserver.NewProtocol6WithError(provider.NewPowerPlatformProvider(context.Background(), false)()),
}

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

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/locations/(europe|unitedstates)/environmentLanguages\?api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/languages/tests/datasource/Validate_Read/get_languages.json").String()), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/locations/(europe|unitedstates)/environmentCurrencies\?api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/currencies/tests/datasource/Validate_Read/get_currencies.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/locations/tests/datasource/Validate_Read/get_locations.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

}
