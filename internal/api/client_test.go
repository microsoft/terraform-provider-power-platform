// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/config"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func TestUnitApiClient_GetConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := config.ProviderConfig{
		UseCli:       false,
		UseOidc:      false,
		UseMsi:       false,
		TenantId:     uuid.NewString(),
		ClientId:     uuid.NewString(),
		ClientSecret: uuid.NewString(),
		TestMode:     true,
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.Execute(ctx, []string{"test"}, "GET", "/relativeurl", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an error for relativeurl but got nil error")
	}

	switch err.(type) {
	case customerrors.UrlFormatError:
		return
	default:
		t.Errorf("Expected error type %s but got %s", reflect.TypeOf(customerrors.UrlFormatError{}), reflect.TypeOf(err))
	}
}

func TestUnitSleepWithContext_TimeoutError(t *testing.T) {
	t.Parallel()

	a := api.NewApiClientBase(&config.ProviderConfig{}, api.NewAuthBase(&config.ProviderConfig{}))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(1)*time.Second)
	err := a.SleepWithContext(ctx, time.Duration(5)*time.Second)
	if err == nil {
		t.Error("Expected an error but got nil error")
	}

	if err.Error() != "context deadline exceeded" {
		t.Errorf("Expected error message %s but got %s", "context deadline exceeded", err.Error())
	}

	cancel()
}

func TestUnitSleepWithContext_HappyPath(t *testing.T) {
	t.Parallel()

	a := api.NewApiClientBase(&config.ProviderConfig{}, api.NewAuthBase(&config.ProviderConfig{}))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	err := a.SleepWithContext(ctx, time.Duration(1)*time.Second)
	if err != nil {
		t.Error("Expected to complete without error but got an error")
	}

	cancel()
}

func TestUnitApiClient_SystemManagedIdentity_No_Identity(t *testing.T) {
	expectedError := "ManagedIdentityCredential: failed to authenticate a system assigned identity."

	ctx := context.Background()
	cfg := config.ProviderConfig{
		UseMsi: true,
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.Execute(ctx, []string{"test"}, "GET", "https://api.bap.microsoft.com", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an authentication error but got nil error")
	}

	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message '%s' but got '%s'", expectedError, err.Error())
	}
}

func TestUnitApiClient_UserManagedIdentity_No_Identity(t *testing.T) {
	expectedError := "ManagedIdentityCredential authentication failed. the requested identity isn't assigned to this resource"

	ctx := context.Background()
	cfg := config.ProviderConfig{
		UseMsi:   true,
		ClientId: uuid.NewString(),
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.Execute(ctx, []string{"test"}, "GET", "https://api.bap.microsoft.com", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an authentication error but got nil error")
	}

	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message '%s' but got '%s'", expectedError, err.Error())
	}
}

func TestUnitApiClient_AzDOWorkloadIdentity_No_TenantId(t *testing.T) {
	expectedError := "tenant ID must be provided to use Azure DevOps Workload Identity Federation"

	ctx := context.Background()
	cfg := config.ProviderConfig{
		UseOidc:                 true,
		AzDOServiceConnectionID: "test",
		ClientId:                "test",
		OidcRequestToken:        "test",
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.Execute(ctx, []string{"test"}, "GET", "https://api.bap.microsoft.com", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an authentication error but got nil error")
	}

	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message '%s' but got '%s'", expectedError, err.Error())
	}
}

func TestUnitApiClient_AzDOWorkloadIdentity_No_ClientId(t *testing.T) {
	expectedError := "client ID must be provided to use Azure DevOps Workload Identity Federation"

	ctx := context.Background()
	cfg := config.ProviderConfig{
		UseOidc:                 true,
		AzDOServiceConnectionID: "test",
		TenantId:                "test",
		OidcRequestToken:        "test",
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.Execute(ctx, []string{"test"}, "GET", "https://api.bap.microsoft.com", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an authentication error but got nil error")
	}

	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message '%s' but got '%s'", expectedError, err.Error())
	}
}

// This is technically not possible with the current control flow but it's still worth testing for.
func TestUnitApiClient_AzDOWorkloadIdentity_No_AzDOServiceConnection(t *testing.T) {
	expectedError := "request URL is required for OIDC credential"

	ctx := context.Background()
	cfg := config.ProviderConfig{
		UseOidc:          true,
		ClientId:         "test",
		TenantId:         "test",
		OidcRequestToken: "test",
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.Execute(ctx, []string{"test"}, "GET", "https://api.bap.microsoft.com", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an authentication error but got nil error")
	}

	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message '%s' but got '%s'", expectedError, err.Error())
	}
}

// This should fall through to base OIDC.
func TestUnitApiClient_AzDOWorkloadIdentity_No_OIDC_Token(t *testing.T) {
	expectedError := "could not obtain an OIDC request token for Azure DevOps Workload Identity Federation"

	ctx := context.Background()
	cfg := config.ProviderConfig{
		UseOidc:                 true,
		ClientId:                "test",
		TenantId:                "test",
		AzDOServiceConnectionID: "test",
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.Execute(ctx, []string{"test"}, "GET", "https://api.bap.microsoft.com", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an authentication error but got nil error")
	}

	if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("Expected error message '%s' but got '%s'", expectedError, err.Error())
	}
}
