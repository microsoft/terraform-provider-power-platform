// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package telemetry_export

type ExportLinksDto struct {
	Value []string `json:"value"`
}

type telemetryExportArrayDto struct {
	Value []ExportLinksDto `json:"value"`
}
