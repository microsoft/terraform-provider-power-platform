// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func TestUnitEnterRequestContext(t *testing.T) {
	ctx := context.Background()
	typ := helpers.TypeInfo{
		ProviderTypeName: "powerplatformtest",
		TypeName:         "testresource",
	}

	req := resource.CreateRequest{}

	newCtx, deferFunc := helpers.EnterRequestContext(ctx, typ, req)

	// Assert
	if newCtx == nil {
		t.Error("EnterRequestContext() returned nil context")
	}

	if deferFunc == nil {
		t.Error("EnterRequestContext() returned nil defer function")
	}

	deferFunc()
}
