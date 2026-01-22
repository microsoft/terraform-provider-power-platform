// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package mocks_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestUnitGetKnownValue_CheckValueAndString(t *testing.T) {
	state := &mocks.StateValue{}
	known := mocks.GetStateValue(state)

	if err := known.CheckValue("ok"); err != nil {
		t.Fatalf("expected CheckValue to succeed, got error: %v", err)
	}

	if state.Value != "ok" {
		t.Fatalf("expected state value to be updated, got %q", state.Value)
	}

	if got := known.String(); got != "ok" {
		t.Fatalf("expected String to return %q, got %q", "ok", got)
	}
}

func TestUnitGetKnownValue_CheckValueTypeMismatch(t *testing.T) {
	state := &mocks.StateValue{}
	known := mocks.GetStateValue(state)

	if err := known.CheckValue(42); err == nil {
		t.Fatal("expected CheckValue to return an error for non-string input")
	}
}

func TestUnitStateValueMatch(t *testing.T) {
	a := &mocks.StateValue{Value: "same"}
	b := &mocks.StateValue{Value: "same"}
	called := false

	check := func(left, right *mocks.StateValue) error {
		called = true
		if left.Value != right.Value {
			return errors.New("values do not match")
		}
		return nil
	}

	checkFunc := mocks.TestStateValueMatch(a, b, check)
	if err := checkFunc(nil); err != nil {
		t.Fatalf("expected check to succeed, got error: %v", err)
	}
	if !called {
		t.Fatal("expected check function to be called")
	}
}
