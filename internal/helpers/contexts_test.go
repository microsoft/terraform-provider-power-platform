// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
// 	"github.com/hashicorp/terraform-plugin-go/tftypes"

// 	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
// )

// type timeoutModel struct {
// 	Timeouts timeouts.Value `tfsdk:"timeouts"`
// }

// type timelessModel struct {
// 	Name string `tfsdk:"name"`
// }

// func TestUnitEnterRequestContext_With_Timeout(t *testing.T) {
// 	ctx := context.Background()
// 	typ := helpers.TypeInfo{
// 		ProviderTypeName: "powerplatformtest",
// 		TypeName:         "testresource",
// 	}

// 	// rawValue := tftypes.NewValue(tftypes.Object{
// 	// 	AttributeTypes: map[string]tftypes.Type{
// 	// 		"timeouts": tftypes.Object{
// 	// 			AttributeTypes: map[string]tftypes.Type{
// 	// 				"create": tftypes.String,
// 	// 				"update": tftypes.String,
// 	// 				"delete": tftypes.String,
// 	// 			},
// 	// 		},
// 	// 	},
// 	// }, map[string]tftypes.Value{
// 	// 	"timeouts": tftypes.NewValue(tftypes.Object{
// 	// 		AttributeTypes: map[string]tftypes.Type{
// 	// 			"create": tftypes.String,
// 	// 			"update": tftypes.String,
// 	// 			"delete": tftypes.String,
// 	// 		},
// 	// 	}, map[string]tftypes.Value{
// 	// 		"create": tftypes.NewValue(tftypes.String, "10m"),
// 	// 		"update": tftypes.NewValue(tftypes.String, "10m"),
// 	// 		"delete": tftypes.NewValue(tftypes.String, "10m"),
// 	// 	}),
// 	// })

// 	req := resource.CreateRequest{
// 		Plan: tfsdk.Plan{
// 			Raw: tftypes.NewValue(timeouts.AttributesWithOpts(), nil),
// 			Schema: schema.Schema{
// 				MarkdownDescription: "This resource manages a PowerPlatform environment",
// 				Description:         "This resource manages a PowerPlatform environment",
		
// 				Attributes: map[string]schema.Attribute{
// 					"timeouts": timeouts.Attributes(ctx),
// 				},
// 			},
// 		},
// 	}
// 	req.Plan.Set(ctx, &timeoutModel{
// 		Timeouts: timeouts.Value{},
// 	})

// 	newCtx, deferFunc := helpers.EnterRequestContext(ctx, typ, req)

// 	// Assert
// 	if newCtx == nil {
// 		t.Error("EnterRequestContext() returned nil context")
// 	}

// 	if deferFunc == nil {
// 		t.Error("EnterRequestContext() returned nil defer function")
// 	}

// 	deferFunc()
// }

// func TestUnitEnterRequestContext_Without_Timeout(t *testing.T) {
// 	ctx := context.Background()
// 	typ := helpers.TypeInfo{
// 		ProviderTypeName: "powerplatformtest",
// 		TypeName:         "testresource",
// 	}
// 	req := resource.CreateRequest{
// 		Plan: tfsdk.Plan{
// 			Raw: tftypes.NewValue(tftypes.Object{}, nil),
// 			Schema: schema.Schema{
// 				MarkdownDescription: "This resource manages a PowerPlatform environment",
// 				Description:         "This resource manages a PowerPlatform environment",
// 			},
// 		},
// 	}
// 	req.Plan.Set(ctx, &timelessModel{
// 	})

// 	newCtx, deferFunc := helpers.EnterRequestContext(ctx, typ, req)

// 	// Assert
// 	if newCtx == nil {
// 		t.Error("EnterRequestContext() returned nil context")
// 	}

// 	if deferFunc == nil {
// 		t.Error("EnterRequestContext() returned nil defer function")
// 	}

// 	deferFunc()
// }
