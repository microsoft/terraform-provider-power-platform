// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package git_integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitConvertSourceControlConfigurationDtoToModel_UsesConcreteProjectNameValue(t *testing.T) {
	model := convertSourceControlConfigurationDtoToModel("00000000-0000-0000-0000-000000000001", scopeSolution, sourceControlConfigurationDto{
		ID:               "11111111-1111-1111-1111-111111111111",
		GitProvider:      0,
		OrganizationName: "example-org",
		ProjectName:      "",
		RepositoryName:   "example-repo",
	})

	require.False(t, model.ProjectName.IsNull())
	require.Equal(t, "", model.ProjectName.ValueString())
}
