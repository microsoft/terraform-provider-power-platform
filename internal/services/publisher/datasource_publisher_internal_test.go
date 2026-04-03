// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package publisher

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUnitPublisherLookupTarget_UsesIDWhenKnown(t *testing.T) {
	got := publisherLookupTarget(DataSourceModel{
		Id:         types.StringValue("11111111-1111-1111-1111-111111111111"),
		UniqueName: types.StringNull(),
	})

	if got != publisherLookupTargetID {
		t.Fatalf("expected id lookup target, got %q", got)
	}
}

func TestUnitPublisherLookupTarget_DefersUnknownID(t *testing.T) {
	got := publisherLookupTarget(DataSourceModel{
		Id:         types.StringUnknown(),
		UniqueName: types.StringNull(),
	})

	if got != publisherLookupTargetDeferred {
		t.Fatalf("expected deferred lookup target for unknown id, got %q", got)
	}
}

func TestUnitPublisherLookupTarget_UsesUniqueNameWhenKnown(t *testing.T) {
	got := publisherLookupTarget(DataSourceModel{
		Id:         types.StringNull(),
		UniqueName: types.StringValue("contoso"),
	})

	if got != publisherLookupTargetUniqueName {
		t.Fatalf("expected uniquename lookup target, got %q", got)
	}
}

func TestUnitPublisherLookupTarget_DefersUnknownUniqueName(t *testing.T) {
	got := publisherLookupTarget(DataSourceModel{
		Id:         types.StringNull(),
		UniqueName: types.StringUnknown(),
	})

	if got != publisherLookupTargetDeferred {
		t.Fatalf("expected deferred lookup target for unknown uniquename, got %q", got)
	}
}
