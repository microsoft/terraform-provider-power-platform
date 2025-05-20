// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package environment_application_admin

type applicationUserDto struct {
	ApplicationId string `json:"applicationId"`
}

type applicationUserResponseDto struct {
	Value []applicationUserDataverseDto `json:"value"`
}

// This structure represents the response from Dataverse API for applicationUsers query
type applicationUserDataverseDto struct {
	ApplicationId string `json:"applicationid"`
}
