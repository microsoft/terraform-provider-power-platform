// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package enterprise_policy

const (
	NETWORK_INJECTION_POLICY_TYPE = "NetworkInjection"
	ENCRYPTION_POLICY_TYPE        = "Encryption"
)

type linkEnterprosePolicyDto struct {
	SystemId string `json:"systemId"`
}
