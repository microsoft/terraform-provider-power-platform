// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/stretchr/testify/require"
)

func TestUnitConnectionsWithDelegatedPermissionsHint(t *testing.T) {
	httpErr := func(status int, body string) error {
		return fmt.Errorf("failed to create connection: %w", customerrors.UnexpectedHttpStatusCodeError{
			ExpectedStatusCodes: []int{http.StatusCreated},
			StatusCode:          status,
			StatusText:          http.StatusText(status),
			Body:                []byte(body),
		})
	}
	plainErr := errors.New("connection refused")

	for name, tc := range map[string]struct {
		err      error
		wantHint bool
	}{
		"nil error":                        {err: nil},
		"non http error":                   {err: plainErr},
		"403 other code":                   {err: httpErr(http.StatusForbidden, `{"error":{"code":"ConnectionAuthorizationFailed"}}`)},
		"404 with delegated permissions":   {err: httpErr(http.StatusNotFound, `InsufficientDelegatedPermissions`)},
		"403 insufficient delegated perms": {err: httpErr(http.StatusForbidden, `{"innererror":{"code":"InsufficientDelegatedPermissions"}}`), wantHint: true},
	} {
		t.Run(name, func(t *testing.T) {
			got := withDelegatedPermissionsHint(tc.err)
			if !tc.wantHint {
				require.Equal(t, tc.err, got)
				return
			}
			require.ErrorIs(t, got, tc.err)
			require.Contains(t, got.Error(), constants.CONNECTIVITY_INSUFFICIENT_DELEGATED_PERMISSIONS_ERROR_MSG)
			var target customerrors.UnexpectedHttpStatusCodeError
			require.ErrorAs(t, got, &target)
		})
	}
}
