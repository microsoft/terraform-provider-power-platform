// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

// buildHostUri builds the host uri for an environment or tenant id. It returns an error when the id is too short
// to contain the two character realm suffix (for example an empty id).
func buildHostUri(id, kind, powerPlatformUrl string) (string, error) {
	cleaned := strings.ReplaceAll(id, "-", "")
	if len(cleaned) < 3 {
		return "", fmt.Errorf("invalid %s id %q: expected a GUID", kind, id)
	}
	realm := cleaned[len(cleaned)-2:]
	cleaned = cleaned[:len(cleaned)-2]

	return fmt.Sprintf("%s.%s.%s.%s", cleaned, realm, kind, powerPlatformUrl), nil
}

// BuildEnvironmentHostUri builds the host uri for the environmentid, returning an error for an empty or too short id
// example input: 00000000-0000-0000-0000-000000000123
// example output: 0000000000000001.23.environment.api.powerplatform.com.
func BuildEnvironmentHostUri(environmentId, powerPlatformUrl string) (string, error) {
	return buildHostUri(environmentId, "environment", powerPlatformUrl)
}

// BuildTenantHostUri builds the host uri for the tenantId
// example input: 00000000-0000-0000-0000-000000000123
// example output: 0000000000000001.23.tenant.api.powerplatform.com.
func BuildTenantHostUri(tenantId, powerPlatformUrl string) (string, error) {
	return buildHostUri(tenantId, "tenant", powerPlatformUrl)
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
