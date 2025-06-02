# Title

Hardcoded Depth in `returnExpandSchema`

##

`/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go`

## Problem

The function `returnExpandSchema` uses a hardcoded depth of 10 when invoked from the `Schema` function. This lacks flexibility and can lead to potential issues if the schema requirements change or exceed this hardcoded limit.

## Impact

Hardcoded values reduce code flexibility and adaptability. If the schema needs to handle deeper levels of nesting, this hardcoded value will become a bottleneck, requiring code alteration. Severity: **medium**

## Location

Located in the `Schema` function of the `DataRecordDataSource` struct.

## Code Issue

```go
"expand": returnExpandSchema(10),
```

## Fix

Replace the hardcoded value with a configurable parameter or a named constant that can be updated based on specific requirements.

```go
const DefaultExpandDepth = 10

"expand": returnExpandSchema(DefaultExpandDepth),
```

You can define the constant outside the function, making it easier to maintain and update when needed.