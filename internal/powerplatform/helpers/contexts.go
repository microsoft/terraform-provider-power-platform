// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/constants"
)

type ContextKey string

type ExecutionContextValue struct {
	ProviderVersion string
	OperatingSystem string
	Architecture string
	GoVersion string
}

type RequestContextValue struct{
	ObjectName string
	ObjectType string
	RequestType string
	RequestId string
}

const (
	EXECUTION_CONTEXT_KEY ContextKey = "executionContext"
	REQUEST_CONTEXT_KEY ContextKey = "requestContext"
)


// EnterRequestScope is a helper function that logs the start of a request scope and returns a closure that can be used to defer the exit of the request scope
// This function should be called at the start of a resource or data source request function
// The returned closure should be deferred at the start of the function
// The closure will log the end of the request scope
// The context is updated with the request context so that it can be accessed in lower level functions
func EnterRequestScope(ctx *context.Context, name string, req any) func () {
	reqId := uuid.New().String()
	objType, reqType := getRequestType(req)
	
	tflog.Debug(*ctx, fmt.Sprintf("%s %s START: %s", reqType, objType, name), map[string]any{
		"requestId": reqId,
		"providerVersion": constants.ProviderVersion,
	})

	// Add the request context to the context so that we can access it in lower level functions
	*ctx = context.WithValue(*ctx, REQUEST_CONTEXT_KEY, RequestContextValue { ObjectType: objType, RequestType: reqType, ObjectName: name, RequestId: reqId })

	// This returns a closure that can be used to defer the exit of the request scope
	return func() {
		tflog.Debug(*ctx, fmt.Sprintf("%s %s END: %s", reqType, objType, name))
	}
}

func getRequestType(req any) (string, string) {
	switch req.(type) {
	case resource.CreateRequest:
		return "RESOURCE", "CREATE"
	case resource.ReadRequest:
		return "RESOURCE", "READ"
	case resource.UpdateRequest:
		return "RESOURCE", "UPDATE"
	case resource.DeleteRequest:
		return "RESOURCE", "DELETE"
	case resource.SchemaRequest:
		return "RESOURCE", "SCHEMA"
	case resource.ImportStateRequest:
		return "RESOURCE", "IMPORT"
	case datasource.ReadRequest:
		return "DATA_SOURCE", "READ"
	case datasource.SchemaRequest:
		return "DATA_SOURCE", "SCHEMA"
	case datasource.ConfigureRequest:
		return "DATA_SOURCE", "CONFIGURE"
	case datasource.MetadataRequest:
		return "DATA_SOURCE", "METADATA"
	default:
		return "UNKNOWN", "UNKNOWN"
	}
}
