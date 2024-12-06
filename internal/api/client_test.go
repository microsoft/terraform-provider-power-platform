// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api_test

import (
	"context"
	"net/http"
	"reflect"
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
		UseMi:        false,
		TenantId:     uuid.NewString(),
		ClientId:     uuid.NewString(),
		ClientSecret: uuid.NewString(),
		TestMode:     true,
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.Execute(ctx, []string{"test"}, "GET", "/relativeurl", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an error for relatvieurl but got nil error")
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
