// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitParseShareImportId(t *testing.T) {
	const (
		env        = "00000000-0000-0000-0000-000000000000"
		connector  = "shared_commondataserviceforapps"
		connection = "00000000-0000-0000-0000-000000000001"
		principal  = "0000000a-000b-000c-000d-00000000000e"
	)

	for name, tc := range map[string]struct {
		importId string
		want     shareImportId
		wantErr  string
	}{
		"valid": {
			importId: env + "/" + connector + "/" + connection + "/" + principal,
			want:     shareImportId{EnvironmentId: env, ConnectorName: connector, ConnectionId: connection, PrincipalId: principal},
		},
		"uppercase principal is kept as given": {
			importId: env + "/" + connector + "/" + connection + "/0000000A-000B-000C-000D-00000000000E",
			want:     shareImportId{EnvironmentId: env, ConnectorName: connector, ConnectionId: connection, PrincipalId: "0000000A-000B-000C-000D-00000000000E"},
		},
		"legacy single id": {
			importId: principal,
			wantErr:  "expected import ID in format 'environment_id/connector_name/connection_id/principal_entra_object_id'",
		},
		"three parts": {
			importId: env + "/" + connector + "/" + connection,
			wantErr:  "expected import ID in format",
		},
		"five parts": {
			importId: env + "/" + connector + "/" + connection + "/" + principal + "/extra",
			wantErr:  "expected import ID in format",
		},
		"trailing slash": {
			importId: env + "/" + connector + "/" + connection + "/" + principal + "/",
			wantErr:  "expected import ID in format",
		},
		"empty environment": {
			importId: "/" + connector + "/" + connection + "/" + principal,
			wantErr:  "environment_id in import ID",
		},
		"empty connector": {
			importId: env + "//" + connection + "/" + principal,
			wantErr:  "connector_name in import ID",
		},
		"empty connection": {
			importId: env + "/" + connector + "//" + principal,
			wantErr:  "connection_id in import ID",
		},
		"empty principal": {
			importId: env + "/" + connector + "/" + connection + "/",
			wantErr:  "principal_entra_object_id in import ID",
		},
		"whitespace around the id": {
			importId: " " + env + "/" + connector + "/" + connection + "/" + principal + " ",
			wantErr:  "environment_id in import ID",
		},
		"whitespace around a part": {
			importId: env + "/ " + connector + " /" + connection + "/" + principal,
			wantErr:  "connector_name in import ID",
		},
		"environment not a guid": {
			importId: "not-a-guid/" + connector + "/" + connection + "/" + principal,
			wantErr:  "environment_id 'not-a-guid' in import ID is not a guid",
		},
		"principal not a guid": {
			importId: env + "/" + connector + "/" + connection + "/admin@contoso.com",
			wantErr:  "principal_entra_object_id 'admin@contoso.com' in import ID is not a guid",
		},
		"principal braced guid": {
			importId: env + "/" + connector + "/" + connection + "/{" + principal + "}",
			wantErr:  "is not a guid",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := parseShareImportId(tc.importId)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
