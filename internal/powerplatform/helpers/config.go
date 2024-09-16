// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package helpers

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetConfigString returns the value of the configValue if it is not null, otherwise it returns the value of the
// environmentVariableName environment variable if it is set, otherwise it returns the defaultValue.
func GetConfigString(ctx context.Context, configValue basetypes.StringValue, environmentVariableName string, defaultValue string) string {
	if !configValue.IsNull() {
		return configValue.ValueString()
	} else if value, ok := os.LookupEnv(environmentVariableName); ok {
		return value
	} else {
		return defaultValue
	}
}

// GetConfigBool returns the value of the configValue if it is not null, otherwise it returns the value of the default value.
func GetConfigBool(ctx context.Context, configValue basetypes.BoolValue, environmentVariableName string, defaultValue bool) bool {
	if !configValue.IsNull() {
		return configValue.ValueBool()
	} else if value, ok := os.LookupEnv(environmentVariableName); ok {
		envValue, err := strconv.ParseBool(value)
		if err == nil {
			return envValue
		} else {
			tflog.Warn(ctx, "Failed to parse environment variable value as a boolean. Using default value instead.", map[string]interface{}{environmentVariableName:value})
			return defaultValue
		}
	} else {
		return defaultValue
	}
}
