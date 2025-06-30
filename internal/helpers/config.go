// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetConfigString returns the value of the configValue if it is not null, otherwise it returns the value of the
// environmentVariableName environment variable if it is set, otherwise it returns the defaultValue.
func GetConfigString(ctx context.Context, configValue basetypes.StringValue, environmentVariableName string, defaultValue string) string {
	if !configValue.IsNull() {
		return configValue.ValueString()
	} else if value, ok := os.LookupEnv(environmentVariableName); ok && value != "" {
		return value
	}
	return defaultValue
}

// GetConfigMultiString returns the value of the configValue if it is not null, otherwise it returns the value of the
// first environment variable that is set, otherwise it returns the defaultValue.
func GetConfigMultiString(ctx context.Context, configValue basetypes.StringValue, environmentVariableNames []string, defaultValue string) string {
	if !configValue.IsNull() {
		return configValue.ValueString()
	}

	for _, k := range environmentVariableNames {
		if value, ok := os.LookupEnv(k); ok && value != "" {
			return value
		}
	}

	return defaultValue
}

func GetListStringValues(value types.List, environmentVariableNames []string, defaultValue []string) types.List {
	// Populate the list with environment variable values or default values if the list is empty.
	if value.IsUnknown() || value.IsNull() {
		values := []attr.Value{}

		for _, k := range environmentVariableNames {
			if value, ok := os.LookupEnv(k); ok && value != "" {
				values = append(values, types.StringValue(strings.TrimSpace(value)))
			}
		}

		if len(values) == 0 {
			for _, v := range defaultValue {
				values = append(values, types.StringValue(strings.TrimSpace(v)))
			}
		}

		return types.ListValueMust(types.StringType, values)
	}

	return value
}

// GetConfigBool returns the value of the configValue if it is not null, otherwise it returns the value of the default value.
func GetConfigBool(ctx context.Context, configValue basetypes.BoolValue, environmentVariableName string, defaultValue bool) bool {
	if !configValue.IsNull() {
		return configValue.ValueBool()
	} else if value, ok := os.LookupEnv(environmentVariableName); ok && value != "" {
		envValue, err := strconv.ParseBool(value)
		if err != nil {
			tflog.Warn(ctx, "Failed to parse environment variable value as a boolean. Using default value instead.", map[string]any{environmentVariableName: value})
			return defaultValue
		}
		return envValue
	}
	return defaultValue
}

func StringPtr(s string) *string {
	return &s
}

// IsKnown returns true if the value is not null and not unknown.
// This is a helper function to reduce repetitive null/unknown checks.
func IsKnown(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}

// BoolPointer returns a pointer to the bool value if the BoolValue is not null and not unknown, otherwise returns nil.
// This helper reduces boilerplate for null/unknown checks when converting BoolValue to *bool.
func BoolPointer(v basetypes.BoolValue) *bool {
	if !v.IsNull() && !v.IsUnknown() {
		return v.ValueBoolPointer()
	}
	return nil
}
