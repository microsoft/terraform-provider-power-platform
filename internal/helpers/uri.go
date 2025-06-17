// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

// BuildEnvironmentHostUri builds the host uri for the environmentid
// example input: 00000000-0000-0000-0000-000000000123
// example output: 0000000000000001.23.environment.api.powerplatform.com.
func BuildEnvironmentHostUri(environmentId, powerPlatformUrl string) string {
	envId := strings.ReplaceAll(environmentId, "-", "")
	realm := string(envId[len(envId)-2:])
	envId = envId[:len(envId)-2]

	return fmt.Sprintf("%s.%s.environment.%s", envId, realm, powerPlatformUrl)
}

// BuildTenantHostUri builds the host uri for the tenantId
// example input: 00000000-0000-0000-0000-000000000123
// example output: 0000000000000001.23.tenant.api.powerplatform.com.
func BuildTenantHostUri(tenantId, powerPlatformUrl string) string {
	envId := strings.ReplaceAll(tenantId, "-", "")
	realm := string(envId[len(envId)-2:])
	envId = envId[:len(envId)-2]

	return fmt.Sprintf("%s.%s.tenant.%s", envId, realm, powerPlatformUrl)
}

// BuildApiUrl builds a URL for API endpoints with a given host, path, and query parameters.
func BuildApiUrl(host, path string, query url.Values) string {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   host,
		Path:   path,
	}
	if query != nil {
		apiUrl.RawQuery = query.Encode()
	}
	return apiUrl.String()
}

// BuildDataverseApiUrl builds a URL for Dataverse API endpoints.
func BuildDataverseApiUrl(environmentHost, path string, query url.Values) string {
	return BuildApiUrl(environmentHost, path, query)
}

// BuildBapiUrl builds a URL for BAPI endpoints.
func BuildBapiUrl(bapiHost, path string, query url.Values) string {
	return BuildApiUrl(bapiHost, path, query)
}
