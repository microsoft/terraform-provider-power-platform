// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package mocks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type StateValue struct {
	Value string
}

var _ knownvalue.Check = GetKnownValue{}

type GetKnownValue struct {
	value *StateValue
}

func (v GetKnownValue) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for getKnownValue check, got: %T", other)
	}

	v.value.Value = otherVal

	return nil
}

func (v GetKnownValue) String() string {
	return v.value.Value
}

func GetStateValue(value *StateValue) GetKnownValue {
	return GetKnownValue{
		value: value,
	}
}

type StateCheckFunc func(a, b *StateValue) error

func TestStateValueMatch(a, b *StateValue, checkFunc StateCheckFunc) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		return checkFunc(a, b)
	}
}
