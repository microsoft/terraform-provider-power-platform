// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

const TEST_ENVIRONMENT_VARIABLE_NAME = "TEST_ENV_VAR"
const TEST_ENVIRONMENT_VARIABLE_NAME1 = "TEST_ENV_VAR_1"
const TEST_ENVIRONMENT_VARIABLE_NAME2 = "TEST_ENV_VAR_2"

// TestUnitGetConfigString tests the GetConfigString function
// This function should return the value of the configValue if it is not null,
// otherwise it should return the value of the environmentVariableName environment variable if it is set,
// otherwise it should return the defaultValue.
func TestUnitGetConfigString_Matrix(t *testing.T) {
	// Do not run in parallel as we are setting environment variables.

	type testData struct {
		name             string
		configValue      *string
		environmentValue *string
		defaultValue     string
		expectedValue    string
	}

	environmentValue := "environmentValue"
	configValue := "configValue"
	for _, testCase := range []testData{
		{
			name:             "default only",
			configValue:      nil,
			environmentValue: nil,
			defaultValue:     "default",
			expectedValue:    "default",
		},
		{
			name:             "environment set",
			configValue:      nil,
			environmentValue: &environmentValue,
			defaultValue:     "default",
			expectedValue:    environmentValue,
		},
		{
			name:             "config and environment set",
			configValue:      &configValue,
			environmentValue: &environmentValue,
			defaultValue:     "default",
			expectedValue:    configValue,
		},
		{
			name:             "config set",
			configValue:      &configValue,
			environmentValue: nil,
			defaultValue:     "default",
			expectedValue:    configValue,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var configValue basetypes.StringValue

			if testCase.environmentValue != nil {
				t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
			}

			if testCase.configValue != nil {
				configValue = basetypes.NewStringValue(*testCase.configValue)
			} else {
				configValue = basetypes.NewStringNull()
			}

			ctx := context.Background()

			result := helpers.GetConfigString(ctx, configValue, TEST_ENVIRONMENT_VARIABLE_NAME, testCase.defaultValue)

			if result != testCase.expectedValue {
				t.Errorf("Expected '%s', got '%s'", testCase.expectedValue, result)
			}
		})
	}
}

// TestUnitGetConfigBool tests the GetConfigBool function
// This function should return the value of the configValue if it is not null,
// otherwise it should return the value of the environmentVariableName environment variable if it is set,
// (if the environment variable can not be parsed as a bool it should return the defaultValue.)
// otherwise it should return the defaultValue.
func TestUnitGetConfigBool_Matrix(t *testing.T) {
	// Do not run in parallel as we are setting environment variables.

	type testData struct {
		name             string
		configValue      *bool
		environmentValue *string
		defaultValue     bool
		expectedValue    bool
	}

	trueValue1 := "TRUE"
	trueValue2 := "true"
	trueValue3 := "1"
	falseValue1 := "FALSE"
	falseValue2 := "false"
	falseValue3 := "0"
	invalidValue := "invalid"

	trueValue := true
	falseValue := false

	for _, testCase := range []testData{
		{
			name:             "default false",
			configValue:      nil,
			environmentValue: nil,
			defaultValue:     false,
			expectedValue:    false,
		},
		{
			name:             "default true",
			configValue:      nil,
			environmentValue: nil,
			defaultValue:     true,
			expectedValue:    true,
		},
		{
			name:             "environment set to true 1",
			configValue:      nil,
			environmentValue: &trueValue1,
			defaultValue:     false,
			expectedValue:    true,
		},
		{
			name:             "environment set to true 2",
			configValue:      nil,
			environmentValue: &trueValue2,
			defaultValue:     false,
			expectedValue:    true,
		},
		{
			name:             "environment set to true 3",
			configValue:      nil,
			environmentValue: &trueValue3,
			defaultValue:     false,
			expectedValue:    true,
		},
		{
			name:             "environment set to false 1",
			configValue:      nil,
			environmentValue: &falseValue1,
			defaultValue:     true,
			expectedValue:    false,
		},
		{
			name:             "environment set to false 2",
			configValue:      nil,
			environmentValue: &falseValue2,
			defaultValue:     true,
			expectedValue:    false,
		},
		{
			name:             "environment set to false 3",
			configValue:      nil,
			environmentValue: &falseValue3,
			defaultValue:     true,
			expectedValue:    false,
		},
		{
			name:             "environment set to invalid with default true",
			configValue:      nil,
			environmentValue: &invalidValue,
			defaultValue:     true,
			expectedValue:    true,
		},
		{
			name:             "environment set to invalid with default false",
			configValue:      nil,
			environmentValue: &falseValue3,
			defaultValue:     false,
			expectedValue:    false,
		},
		{
			name:             "config and environment set",
			configValue:      &trueValue,
			environmentValue: &falseValue1,
			defaultValue:     false,
			expectedValue:    true,
		},
		{
			name:             "config set to true",
			configValue:      &trueValue,
			environmentValue: nil,
			defaultValue:     false,
			expectedValue:    true,
		},
		{
			name:             "config set to false",
			configValue:      &falseValue,
			environmentValue: nil,
			defaultValue:     true,
			expectedValue:    false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var configValue basetypes.BoolValue

			if testCase.environmentValue != nil {
				t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
			}

			if testCase.configValue != nil {
				configValue = basetypes.NewBoolValue(*testCase.configValue)
			} else {
				configValue = basetypes.NewBoolNull()
			}

			ctx := context.Background()

			result := helpers.GetConfigBool(ctx, configValue, TEST_ENVIRONMENT_VARIABLE_NAME, testCase.defaultValue)

			if result != testCase.expectedValue {
				t.Errorf("Expected '%t', got '%t'", testCase.expectedValue, result)
			}
		})
	}
}

// TestUnitGetConfigMultiString tests the GetConfigMultiString function
// This function should return the value of the configValue if it is not null,
// otherwise it should return the value of the first environment variable that is set,
// otherwise it should return the defaultValue.
func TestUnitGetConfigMultiString_Matrix(t *testing.T) {
	// Do not run in parallel as we are setting environment variables.

	type testData struct {
		name              string
		configValue       *string
		environmentValues []string
		defaultValue      string
		expectedValue     string
	}

	environmentValue1 := "environmentValue1"
	environmentValue2 := "environmentValue2"

	t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME1, environmentValue1)
	t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME2, environmentValue2)

	configValue := "configValue"
	for _, testCase := range []testData{
		{
			name:              "default only",
			configValue:       nil,
			environmentValues: nil,
			defaultValue:      "default",
			expectedValue:     "default",
		},
		{
			name:              "environment set 1",
			configValue:       nil,
			environmentValues: []string{TEST_ENVIRONMENT_VARIABLE_NAME1},
			defaultValue:      "default",
			expectedValue:     environmentValue1,
		},
		{
			name:              "environment set 2",
			configValue:       nil,
			environmentValues: []string{TEST_ENVIRONMENT_VARIABLE_NAME1, TEST_ENVIRONMENT_VARIABLE_NAME2},
			defaultValue:      "default",
			expectedValue:     environmentValue1,
		},
		{
			name:              "environments reversed",
			configValue:       nil,
			environmentValues: []string{TEST_ENVIRONMENT_VARIABLE_NAME2, TEST_ENVIRONMENT_VARIABLE_NAME1},
			defaultValue:      "default",
			expectedValue:     environmentValue2,
		},
		{
			name:              "environment 2 set only",
			configValue:       nil,
			environmentValues: []string{TEST_ENVIRONMENT_VARIABLE_NAME2},
			defaultValue:      "default",
			expectedValue:     environmentValue2,
		},
		{
			name:              "empty environments",
			configValue:       nil,
			environmentValues: []string{},
			defaultValue:      "default",
			expectedValue:     "default",
		},
		{
			name:              "config and environment set",
			configValue:       &configValue,
			environmentValues: []string{TEST_ENVIRONMENT_VARIABLE_NAME1},
			defaultValue:      "default",
			expectedValue:     configValue,
		},
		{
			name:              "config set",
			configValue:       &configValue,
			environmentValues: nil,
			defaultValue:      "default",
			expectedValue:     configValue,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var configValue basetypes.StringValue

			for _, value := range testCase.environmentValues {
				t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, value)
			}

			if testCase.configValue != nil {
				configValue = basetypes.NewStringValue(*testCase.configValue)
			} else {
				configValue = basetypes.NewStringNull()
			}

			ctx := context.Background()

			result := helpers.GetConfigMultiString(ctx, configValue, testCase.environmentValues, testCase.defaultValue)

			if result != testCase.expectedValue {
				t.Errorf("Expected '%s', got '%s'", testCase.expectedValue, result)
			}
		})
	}
}
