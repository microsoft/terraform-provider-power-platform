# Title

Missing Error Handling for `DataRecordClient`

##

`/workspaces/terraform-provider-power-platform/internal/services/data_record/models.go`

## Problem

The `DataRecordClient` field is declared in both `DataRecordDataSource` and `DataRecordResource` structs without proper error-handling mechanisms for its operations or initialization. This could lead to runtime panics if the client fails during its setup or encounters errors during execution.

## Impact

The absence of error handling for the `DataRecordClient` could make the service layer fragile, especially during communication with external dependencies or data services. This can result in unpredictable behavior, performance degradation, and a lack of debug visibility. Severity: **High**

## Location

Found in the following structs:
```go
type DataRecordDataSource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type DataRecordResource struct {
	helpers.TypeInfo
	DataRecordClient client
}
```

## Code Issue

```go
type DataRecordDataSource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type DataRecordResource struct {
	helpers.TypeInfo
	DataRecordClient client
}
```

## Fix

Introduce error-handling logic for the initialization and operations of `DataRecordClient`.

### Adjusted Code Example:
```go
type DataRecordDataSource struct {
	helpers.TypeInfo
	DataRecordClient client
}

// Constructor or Initialization with Error Handling
func (d *DataRecordDataSource) InitializeClient() error {
	var err error
	d.DataRecordClient, err = InitializeClientFunction()
	if err != nil {
		return fmt.Errorf("failed to initialize DataRecordClient: %w", err)
	}
	return nil
}

type DataRecordResource struct {
	helpers.TypeInfo
	DataRecordClient client
}

// Safe operation on DataRecordClient with error checks
func (r *DataRecordResource) FetchData() ([]Data, error) {
	data, err := r.DataRecordClient.GetData()
	if err != nil {
		return nil, fmt.Errorf("error fetching data from DataRecordClient: %w", err)
	}
	return data, nil
}
```

This fix ensures robust error handling and avoids runtime errors that may arise when `DataRecordClient` fails.

