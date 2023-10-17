package powerplatform_mocks

import (
	"io"
	"net/http"

	"github.com/jarcoal/httpmock"
)

const (
	oAuthWellKnownResponse = `{"token_endpoint":"https://login.microsoftonline.com/_/oauth2/v2.0/token","token_endpoint_auth_methods_supported":["client_secret_post","private_key_jwt","client_secret_basic"],"jwks_uri":"https://login.microsoftonline.com/_/discovery/v2.0/keys","response_modes_supported":["query","fragment","form_post"],"subject_types_supported":["pairwise"],"id_token_signing_alg_values_supported":["RS256"],"response_types_supported":["code","id_token","code id_token","id_token token"],"scopes_supported":["openid","profile","email","offline_access"],"issuer":"https://login.microsoftonline.com/_/v2.0","request_uri_parameter_supported":false,"userinfo_endpoint":"https://graph.microsoft.com/oidc/userinfo","authorization_endpoint":"https://login.microsoftonline.com/_/oauth2/v2.0/authorize","device_authorization_endpoint":"https://login.microsoftonline.com/_/oauth2/v2.0/devicecode","http_logout_supported":true,"frontchannel_logout_supported":true,"end_session_endpoint":"https://login.microsoftonline.com/_/oauth2/v2.0/logout","claims_supported":["sub","iss","cloud_instance_name","cloud_instance_host_name","cloud_graph_host_name","msgraph_host","aud","exp","iat","auth_time","acr","nonce","preferred_username","name","tid","ver","at_hash","c_hash","email"],"kerberos_endpoint":"https://login.microsoftonline.com/_/kerberos","tenant_region_scope":"EU","cloud_instance_name":"microsoftonline.com","cloud_graph_host_name":"graph.windows.net","msgraph_host":"graph.microsoft.com","rbac_url":"https://pas.windows.net"}`
	oAuthTokenResponse     = `{
		"token_type": "Bearer",
		"expires_in": 3599,
		"ext_expires_in": 3599,
		"access_token": "eyJ0eXAiOiJKV1Q"
	}`
)

func ActivateOAuthHttpMocks() {

	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		println("No HttpMock responder for: " + req.Method + " " + req.URL.String())
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			println("Body:" + string(body))
		}
		return httpmock.NewStringResponse(http.StatusTeapot, ""), nil
	})

	httpmock.RegisterResponder("GET", `=~^https://login\.microsoftonline\.com/*./v2.0/\.well-known/openid-configuration`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, oAuthWellKnownResponse), nil
		})

	httpmock.RegisterResponder("POST", `=~^https://login\.microsoftonline\.com/*./oauth2/v2.0/token`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, oAuthTokenResponse), nil
		})
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

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
		})

}
