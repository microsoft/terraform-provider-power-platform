// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

const TEST_ENV_VARIRONMENT_VARIABLE_NAME = "TEST_ENV_VAR"

// func TestGetConfigString_ConfigValueIsNotNull(t *testing.T) {
// 	expectedValue := "expectedValue"
// 	t.Setenv(TEST_ENV_VARIRONMENT_VARIABLE_NAME, "environmentVariableValue")
// 	ctx := context.Background()
// 	configValue := basetypes.NewStringValue(expectedValue)
// 	defaultValue := "defaultValue"

// 	result := helpers.GetConfigString(ctx, configValue, TEST_ENV_VARIRONMENT_VARIABLE_NAME, defaultValue)

// 	if result != expectedValue {
// 		t.Errorf("Expected '%s', got %s", expectedValue, result)
// 	}
// }

// func TestGetConfigString_EnvironmentVariableIsSet(t *testing.T) {
// 	ctx := context.Background()
// 	configValue := basetypes.NewStringNull()
// 	defaultValue := "defaultValue"
// 	os.Setenv(TEST_ENV_VARIRONMENT_VARIABLE_NAME, "environmentVariableValue")

// 	result := helpers.GetConfigString(ctx, configValue, TEST_ENV_VARIRONMENT_VARIABLE_NAME, defaultValue)

// 	if result != "environmentVariableValue" {
// 		t.Errorf("Expected 'environmentVariableValue', got %s", result)
// 	}
// }

func TestUnitGetConfigString_Matrix(t *testing.T) {
	type testData struct{
		name string
		configValue *string
		environmentValue *string
		defaultValue string
		expectedValue string
	}

	environmentValue := "environmentValue"
	configValue := "configValue"
	for _, testCase := range []testData {
		{
			name : "default only",
			configValue : nil,
			environmentValue : nil,
			defaultValue : "default",
			expectedValue : "default",
		},
		{
			name : "environment set",
			configValue : nil,
			environmentValue : &environmentValue,
			defaultValue : "default",
			expectedValue : environmentValue,
		},
		{
			name : "config and environment set",
			configValue : &configValue,
			environmentValue : &environmentValue,
			defaultValue : "default",
			expectedValue : configValue,
		},
		{
			name : "config set",
			configValue : &configValue,
			environmentValue : nil,
			defaultValue : "default",
			expectedValue : configValue,
		},
	} {
		t.Run(testCase.name, func(t *testing.T){
			var configValue basetypes.StringValue

			if(testCase.environmentValue != nil) {
				t.Setenv(TEST_ENV_VARIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
			}

			if( testCase.configValue != nil) {
				configValue = basetypes.NewStringValue(*testCase.configValue)
			} else {
				configValue = basetypes.NewStringNull()
			}
			
			ctx := context.Background()
		
			result := helpers.GetConfigString(ctx, configValue, TEST_ENV_VARIRONMENT_VARIABLE_NAME, testCase.defaultValue)
		
			if result != testCase.expectedValue {
				t.Errorf("Expected '%s', got '%s'", testCase.expectedValue, result)
			}
		})
	}
}
