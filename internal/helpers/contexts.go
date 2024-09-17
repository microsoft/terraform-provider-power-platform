// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/common"
)

// ContextKey is a custom type for context keys
// A custom type is needed to avoid collisions with other packages that use the same key
type ContextKey string

type ExecutionContextValue struct {
	ProviderVersion string
	OperatingSystem string
	Architecture    string
	GoVersion       string
}

// RequestContextValue is a struct that holds the object type, request type, object name and request id for a given request
// This struct is used to store the request context in the context so that it can be accessed in lower level functions
type RequestContextValue struct {
	ObjectName  string
	RequestType string
	RequestId   string
}

// Context keys for the execution and request context
const (
	EXECUTION_CONTEXT_KEY ContextKey = "executionContext"
	REQUEST_CONTEXT_KEY   ContextKey = "requestContext"
)

// EnterRequestScope is a helper function that logs the start of a request scope and returns a closure that can be used to defer the exit of the request scope
// This function should be called at the start of a resource or data source request function
// The returned closure should be deferred at the start of the function
// The closure will log the end of the request scope
// The context is updated with the request context so that it can be accessed in lower level functions
func EnterRequestContext[T AllowedRequestTypes](ctx context.Context, typ TypeInfo, req T) (context.Context, func()) {
	reqId := strings.ReplaceAll(uuid.New().String(), "-", "")
	reqType := reflect.TypeOf(req).String()
	name := typ.FullTypeName()

	tflog.Debug(ctx, fmt.Sprintf("%s START: %s", reqType, name), map[string]any{
		"requestId":       reqId,
		"providerVersion": common.ProviderVersion,
	})

	// Add the request context to the context so that we can access it in lower level functions
	ctx = context.WithValue(ctx, REQUEST_CONTEXT_KEY, RequestContextValue{RequestType: reqType, ObjectName: name, RequestId: reqId})
	ctx = tflog.SetField(ctx, "request_id", reqId)
	ctx = tflog.SetField(ctx, "request_type", reqType)

	// This returns a closure that can be used to defer the exit of the request scope
	return ctx, func() {
		tflog.Debug(ctx, fmt.Sprintf("%s END: %s", reqType, name))
	}
}

// EnterProviderScope is a helper function that logs the start of a provider scope and returns a closure that can be used to defer the loging of the exit of the provider scope
func EnterProviderContext[T AllowedProviderRequestTypes](ctx context.Context, req T) (context.Context, func()) {
	reqType := reflect.TypeOf(req).String()

	tflog.Debug(ctx, fmt.Sprintf("%s START", reqType), map[string]any{
		"providerVersion": common.ProviderVersion,
	})

	// This returns a closure that can be used to defer the exit of the provider scope
	return ctx, func() {
		tflog.Debug(ctx, fmt.Sprintf("%s END", reqType), map[string]any{
			"providerVersion": common.ProviderVersion,
		})
	}
}

// AllowedRequestTypes is an interface that defines the allowed request types for the getRequestTypeName function
type AllowedRequestTypes interface {
    resource.CreateRequest | 
	resource.MetadataRequest |
	resource.ReadRequest |
	resource.UpdateRequest |
	resource.DeleteRequest |
	resource.SchemaRequest |
	resource.ConfigureRequest |
	resource.ModifyPlanRequest |
	resource.ImportStateRequest |
	resource.UpgradeStateRequest |
	resource.ValidateConfigRequest |
	datasource.ReadRequest |
	datasource.SchemaRequest |
	datasource.ConfigureRequest |
	datasource.MetadataRequest |
	datasource.ValidateConfigRequest
}

// AllowedProviderRequestTypes is an interface that defines the allowed request types for the EnterProviderContext function
type AllowedProviderRequestTypes interface {
	provider.ConfigureRequest |
	provider.MetaSchemaRequest |
	provider.MetadataRequest |
	provider.SchemaRequest |
	provider.ValidateConfigRequest
}
