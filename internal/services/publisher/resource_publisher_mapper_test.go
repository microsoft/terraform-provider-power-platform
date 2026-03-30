// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher

import "testing"

func TestUnitAddressModelFromDto_IgnoresEmptySlotWithOnlyAddressID(t *testing.T) {
	dto := &publisherDto{
		Address2AddressId: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
	}

	model := addressModelFromDto(2, dto)
	if model != nil {
		t.Fatalf("expected address slot 2 to be ignored when only address id remains, got %#v", model)
	}
}
