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
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func TestApiClient_GetConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg := config.ProviderConfig{
		UseCli:       false,
		UseOidc:      false,
		TenantId:     uuid.NewString(),
		ClientId:     uuid.NewString(),
		ClientSecret: uuid.NewString(),
		TestMode:     true,
	}

	x := api.NewApiClientBase(&cfg, api.NewAuthBase(&cfg))
	_, err := x.ExecuteForGivenScope(ctx, "test", http.MethodGet, "/relativeurl", http.Header{}, nil, []int{http.StatusOK}, nil)
	if err == nil {
		t.Error("Expected an error for relatvieurl but got nil error")
	}

	if e, ok := err.(helpers.ProviderError); !ok {
		t.Errorf("Expected error type %s but got %s", reflect.TypeFor[helpers.ProviderError](), reflect.TypeOf(e.ErrorCode))
	} else if e.ErrorCode != helpers.ERROR_INCORRECT_URL_FORMAT {
		t.Errorf("Expected error code %s but got %s", helpers.ERROR_INCORRECT_URL_FORMAT, e.ErrorCode)
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
