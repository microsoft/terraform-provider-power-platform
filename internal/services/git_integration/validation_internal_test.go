// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitNormalizeSolutionID_RejectsInvalidGUIDSegments(t *testing.T) {
	_, err := normalizeSolutionID("00000000-0000-0000-0000-000000000001", "not-a-guid_33333333-3333-3333-3333-333333333333")
	require.ErrorContains(t, err, "valid environment GUID prefix")

	_, err = normalizeSolutionID("00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000001_not-a-guid")
	require.ErrorContains(t, err, "valid Dataverse solution GUID suffix")
}

func TestUnitNormalizeSolutionID_AllowsCaseInsensitiveEnvironmentMatch(t *testing.T) {
	solutionID, err := normalizeSolutionID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA_33333333-3333-3333-3333-333333333333")
	require.NoError(t, err)
	require.Equal(t, "33333333-3333-3333-3333-333333333333", solutionID)

	solutionID, err = normalizeSolutionID("AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa_33333333-3333-3333-3333-333333333333")
	require.NoError(t, err)
	require.Equal(t, "33333333-3333-3333-3333-333333333333", solutionID)
}
