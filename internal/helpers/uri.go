// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"fmt"
	"strings"
)

// BuildHostUri builds the host uri for the environmentid
// example input: 00000000-0000-0000-0000-000000000123
// example output: 0000000000000001.23.api.powerplatform.com.
func BuildHostUri(environmentId, powerPlatformUrl string) string {
	envId := strings.ReplaceAll(environmentId, "-", "")
	realm := string(envId[len(envId)-2:])
	envId = envId[:len(envId)-2]

	return fmt.Sprintf("%s.%s.environment.%s", envId, realm, powerPlatformUrl)
}
