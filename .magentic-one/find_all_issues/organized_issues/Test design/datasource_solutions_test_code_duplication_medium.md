# Title

Duplication of Test Logic without Table-Driven Test Pattern

##

/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions_test.go

## Problem

The current unit tests for different cases (for example, validating a read or the absence of Dataverse) duplicate a significant amount of setup and mock logic, making the tests harder to maintain, extend, and scale as new scenarios are added. This is considered inefficient, and could be greatly improved by adopting the idiomatic Go table-driven test pattern.

## Impact

The duplicated code increases the risk of inconsistencies, escalates maintenance effort when changes are required, and can reduce clarity. While this does not introduce a bug, it is a medium-severity codebase maintainability issue.

## Location

Spread throughout the different functions:

```go
func TestUnitSolutionsDataSource_Validate_Read(t *testing.T) { ... }
func TestUnitSolutionsDataSource_Validate_No_Dataverse(t *testing.T) { ... }
```

## Code Issue

```go
// Two full sets of mocks and test cases, separated in functions, but nearly identical in structure.
// Any further "Validate_..." scenarios will continue the pattern instead of reusing logic.
```

## Fix

Refactor into a table-driven style to improve reuse and reduce duplication:

```go
func TestUnitSolutionsDataSource_TableDriven(t *testing.T) {
	type testCase struct {
		name    string
		setup   func()
		config  string
		check   resource.TestCheckFunc
		wantErr *regexp.Regexp
	}
	tests := []testCase{
		{
			name: "Read",
			setup: func() {
				// put all relevant registers here
			},
			config: `...`, // test config
			check:  /* resource.ComposeAggregateTestCheckFunc(...) */,
		},
		{
			name: "NoDataverse",
			setup: func() {
				// put all relevant registers here
			},
			config: `...`,
			wantErr: regexp.MustCompile(`No Dataverse exists in environment`),
			check:   resource.ComposeAggregateTestCheckFunc(),
		},
		// add more test cases if desired
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{{
					Config:      tt.config,
					ExpectError: tt.wantErr,
					Check:       tt.check,
				}},
			})
		})
	}
}
```
