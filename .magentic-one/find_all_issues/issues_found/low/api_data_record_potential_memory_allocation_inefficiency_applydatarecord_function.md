# Title

Potential memory allocation inefficiency in `ApplyDataRecord` function

##

Path to the file `/workspaces/terraform-provider-power-platform/internal/services/data_record/api_data_record.go`

## Problem

In the `ApplyDataRecord` function, nested maps are iterated and deleted from the `columns` map during runtime. While functional, this creates potential memory inefficiencies as the map structure is heavily modified. Additionally, the declaration `relations := make(map[string]any, 0)` is non-optimal; the initial capacity specified as `0` does not provide benefits, and a better approach would be pre-calculating the required capacity.

## Impact

Frequent map modifications during runtime can cause unnecessary memory reallocation, leading to performance bottlenecks in high-load environments. Severity: **Low**

## Location

Occurs in the method `ApplyDataRecord`:

## Code Issue

```go
relations := make(map[string]any, 0)

for key, value := range columns {
    if nestedMap, ok := value.(map[string]any); ok {
        delete(columns, key)
        // Modification logic within iteration
    }
}
```

## Fix

Use a temporary map to collect keys to be deleted before modifying the `columns` map.

```go
relations := make(map[string]any) // Remove unnecessary capacity
keysToDelete := []string{}        // Collect keys to delete

for key, value := range columns {
    if nestedMap, ok := value.(map[string]any); ok {
        keysToDelete = append(keysToDelete, key)
        if len(nestedMap) > 0 {
            tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
            if err != nil {
                return nil, err
            }

            entityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableLogicalName)
            if err != nil {
                return nil, err
            }

            columns[fmt.Sprintf("%s@odata.bind", key)] = fmt.Sprintf("/%s(%s)", entityDefinition.LogicalCollectionName, dataRecordId)
        }
    } else if nestedMapList, ok := value.([]any); ok {
        keysToDelete = append(keysToDelete, key)
        relations[key] = nestedMapList
    }
}

// Perform deletions after iteration
for _, key := range keysToDelete {
    delete(columns, key)
}
```

This approach avoids modifying the map during iteration and improves memory usage.

---

Next, I will proceed to identify additional issues.