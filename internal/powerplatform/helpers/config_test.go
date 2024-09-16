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


func TestUnitGetConfigBool_Matrix(t *testing.T) {
	type testData struct{
		name string
		configValue *bool
		environmentValue *string
		defaultValue bool
		expectedValue bool
	}

	trueValue1 := "TRUE"
	trueValue2 := "true"
	trueValue3 := "1"
	falseValue1 := "FALSE"
	falseValue2 := "false"
	falseValue3 := "0"

	trueValue := true
	//falseValue := false

	for _, testCase := range []testData {
		{
			name : "default false",
			configValue : nil,
			environmentValue : nil,
			defaultValue : false,
			expectedValue : false,
		},
		{
			name : "default true",
			configValue : nil,
			environmentValue : nil,
			defaultValue : true,
			expectedValue : true,
		},
		{
			name : "environment set to true 1",
			configValue : nil,
			environmentValue : &trueValue1,
			defaultValue : false,
			expectedValue : true,
		},
		{
			name : "environment set to true 2",
			configValue : nil,
			environmentValue : &trueValue2,
			defaultValue : false,
			expectedValue : true,
		},
		{
			name : "environment set to true 3",
			configValue : nil,
			environmentValue : &trueValue3,
			defaultValue : false,
			expectedValue : true,
		},
		{
			name : "environment set to false 1",
			configValue : nil,
			environmentValue : &falseValue1,
			defaultValue : true,
			expectedValue : false,
		},
		{
			name : "environment set to false 2",
			configValue : nil,
			environmentValue : &falseValue2,
			defaultValue : true,
			expectedValue : false,
		},
		{
			name : "environment set to false 3",
			configValue : nil,
			environmentValue : &falseValue3,
			defaultValue : true,
			expectedValue : false,
		},
		{
			name : "config and environment set",
			configValue : &trueValue,
			environmentValue : &falseValue1,
			defaultValue : false,
			expectedValue : true,
		},
		{
			name : "config set",
			configValue : &trueValue,
			environmentValue : nil,
			defaultValue : false,
			expectedValue : true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T){
			var configValue basetypes.BoolValue

			if(testCase.environmentValue != nil) {
				t.Setenv(TEST_ENV_VARIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
			}

			if( testCase.configValue != nil) {
				configValue = basetypes.NewBoolValue(*testCase.configValue)
			} else {
				configValue = basetypes.NewBoolNull()
			}
			
			ctx := context.Background()
		
			result := helpers.GetConfigBool(ctx, configValue, TEST_ENV_VARIRONMENT_VARIABLE_NAME, testCase.defaultValue)
		
			if result != testCase.expectedValue {
				t.Errorf("Expected '%t', got '%t'", testCase.expectedValue, result)
			}
		})
	}
}
