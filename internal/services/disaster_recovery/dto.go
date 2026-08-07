// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package disaster_recovery

import "github.com/microsoft/terraform-provider-power-platform/internal/services/environment"

type disasterRecoveryPatchDto struct {
	Properties disasterRecoveryPatchPropertiesDto `json:"properties"`
}

type disasterRecoveryPatchPropertiesDto struct {
	States disasterRecoveryPatchStatesDto `json:"states"`
}

type disasterRecoveryPatchStatesDto struct {
	DisasterRecovery environment.DisasterRecoveryStateDto `json:"disasterRecovery"`
}
