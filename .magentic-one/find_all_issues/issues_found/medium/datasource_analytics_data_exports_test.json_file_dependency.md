# Title

Hardcoded Path to JSON Files in Unit Test

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports_test.go

## Problem

The unit test for `TestUnitAnalyticsDataExportsDataSource_Validate_Read` relies on hardcoded paths to JSON files (`tests/datasource/...`) for mock data. This introduces brittle test design as the location of test data files must remain static, introducing additional maintenance overhead if the directory structure of the project changes.

## Impact

In case of directory structure changes, the test may break, introducing frustrating debugging cycles. Furthermore, if the hardcoded path changes or mock data accidentally goes missing, critical unit testing functionality is lost. Severity: **medium**

## Location

File: `datasource_analytics_data_exports_test.go` 

Function: `TestUnitAnalyticsDataExportsDataSource_Validate_Read`

## Code Issue

```go
// Register mock response for tenant API
httpmock.RegisterResponder(
	"GET",
	"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
	httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant.json").String()))

// Register responder for gateway cluster API with dynamic hostname pattern
httpmock.RegisterResponder(
	"GET",
	`=~^https://.*\.tenant\.api\.powerplatform\.com/gateway/cluster.*`,
	httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_gateway_cluster.json").String()))

// Register responder for analytics data exports API
httpmock.RegisterResponder(
	"GET",
	"https://na.csanalytics.powerplatform.microsoft.com/api/v2/connections",
	httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_analytics_data_exports.json").String()))
```

## Fix

Avoid hardcoded paths by introducing an abstraction for managing test data files. For instance:

```go
func getTestData(fileName string) (string, error) {
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		dataPath = "./tests/datasource/Validate_Read"
	}
	fullPath := filepath.Join(dataPath, fileName)
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read test data file %s: %v", fullPath, err)
	}
	return string(data), nil
}

respondData, err := getTestData("get_tenant.json")
if err != nil {
	t.Fatalf("%v", err)
}
httpmock.RegisterResponder(
	"GET",
	"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
	httpmock.NewStringResponder(http.StatusOK, respondData))
```