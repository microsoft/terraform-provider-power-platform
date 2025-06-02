# Title

Reflection-based filtering lacks unit tests

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/api_tenant_settings.go`

## Problem

The core functionality of the `filterDto` method relies heavily on reflection (`reflect.TypeOf`, `reflect.VisibleFields`, etc.) to filter properties dynamically based on the `configuredSettings` and `backendSettings`. Reflection code is prone to subtle bugs and regressions, especially when handling invalid or unexpected input. Testing such code is essential to ensure reliability. However, there is no indication or reference to any test suite verifying the correctness of this function or handling edge cases (e.g., mismatched types, nil values, invalid reflection operations).

## Impact

This issue impacts the codebase in the following ways:
- Increases the risk of unnoticed bugs and regressions, especially in code paths using reflection.
- Reduces the reliability of the `filterDto` function, which is pivotal for filtering data returned by the backend.
- Additional manual debugging effort may be required to resolve issues if the function fails unexpectedly in production.

Severity: **High**

## Location

`filterDto` function, likely used in conjunction with other business logic (e.g., `applyCorrections`).

## Code Issue

```go
func filterDto(ctx context.Context, configuredSettings any, backendSettings any) any {
	configuredType := reflect.TypeOf(configuredSettings)
	backendType := reflect.TypeOf(backendSettings)
	if configuredType != backendType {
		return nil
	}

	output := reflect.New(configuredType).Interface()

	visibleFields := reflect.VisibleFields(configuredType)

	configuredValue := reflect.ValueOf(configuredSettings)
	backendValue := reflect.ValueOf(backendSettings)

	for fieldIndex, fieldInfo := range visibleFields {
		tflog.Debug(ctx, fmt.Sprintf("Field: %s", fieldInfo.Name))

		configuredFieldValue := configuredValue.Field(fieldIndex)
		backendFieldValue := backendValue.Field(fieldIndex)
		outputField := reflect.ValueOf(output).Elem().Field(fieldIndex)

		if !configuredFieldValue.IsNil() && !backendFieldValue.IsNil() && backendFieldValue.IsValid() && outputField.CanSet() {
			if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Struct {
				outputStruct := filterDto(ctx, configuredFieldValue.Elem().Interface(), backendFieldValue.Elem().Interface())
				outputField.Set(reflect.ValueOf(outputStruct))
			} else if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Bool {
				boolValue := backendFieldValue.Elem().Bool()
				newBool := bool(boolValue)
				outputField.Set(reflect.ValueOf(&newBool))
			} else if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.String {
				stringValue := backendFieldValue.Elem().String()
				newString := string(stringValue)
				outputField.Set(reflect.ValueOf(&newString))
			} else if configuredFieldValue.Kind() == reflect.Pointer && configuredFieldValue.Elem().Kind() == reflect.Int64 {
				int64Value := backendFieldValue.Elem().Int()
				newInt64 := int64(int64Value)
				outputField.Set(reflect.ValueOf(&newInt64))
			} else {
				tflog.Debug(ctx, fmt.Sprintf("Skipping unknown field type %s", configuredFieldValue.Kind()))
			}
		}
	}

	return output
}
```

## Fix

Create a dedicated suite of unit tests to verify the functionality of `filterDto`. Test cases should cover the following scenarios:
- Basic functionality with standard configurations.
- Mismatched data types between `configuredSettings` and `backendSettings`.
- Handling `nil` values gracefully.
- Handling invalid `reflect` operations.
- Edge cases such as empty configurations, invalid values, etc.

Example Test Suite:

```go
func TestFilterDto(t *testing.T) {
	// Create mock data for configuredSettings and backendSettings
	configuredSettings := MockConfiguredSettings{Field1: "value1", Field2: nil}
	backendSettings := MockBackendSettings{Field1: "value1", Field2: "value2"}

	// Call the filterDto function and verify output
	result := filterDto(context.TODO(), configuredSettings, backendSettings)

	// Assert that only opted-in fields are present in the output
	assert.NotNil(t, result)
	assert.Equal(t, "value1", result.Field1)
	assert.Nil(t, result.Field2)
}

func TestFilterDtoMismatchedTypes(t *testing.T) {
	// Create mock data with mismatched types
	configuredSettings := MockConfiguredSettings{}
	backendSettings := AnotherType{}

	// Call the filterDto function and verify it returns nil
	result := filterDto(context.TODO(), configuredSettings, backendSettings)

	assert.Nil(t, result)
}

// Other tests for reflection errors, edge cases, etc.
```

Explanation:
- `MockConfiguredSettings` and `MockBackendSettings` are structs simulating real configurations.
- `TestFilterDto` verifies expected functionality.
- `TestFilterDtoMismatchedTypes` handles cases with mismatched data types.
