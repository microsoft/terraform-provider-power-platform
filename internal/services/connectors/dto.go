// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connectors

type connectorDto struct {
	Name       string                 `json:"name"`
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties connectorPropertiesDto `json:"properties"`
}

type connectorPropertiesDto struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Tier        string `json:"tier"`
	Publisher   string `json:"publisher"`
	Unblockable bool
}

type connectorArrayDto struct {
	Value []connectorDto `json:"value"`
}

type unblockableConnectorDto struct {
	Id       string                          `json:"id"`
	Metadata unblockableConnectorMetadataDto `json:"metadata"`
}

type unblockableConnectorMetadataDto struct {
	Unblockable bool `json:"unblockable"`
}

type virtualConnectorDto struct {
	Id       string                      `json:"id"`
	Metadata virtualConnectorMetadataDto `json:"metadata"`
}

type virtualConnectorMetadataDto struct {
	VirtualConnector bool   `json:"virtualConnector"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	DisplayName      string `json:"displayName"`
}
