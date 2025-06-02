# Title

Typo in `BillingPoliciesEnvironmetsListDataSourceModel` name

### 

/`internal/services/licensing/datasource_billing_policies_environments.go`

### Problem

The object name `BillingPoliciesEnvironmetsListDataSourceModel` contains a typo in "Environmets"; it should be "Environments".

### Impact

- Leads to potential confusion and might reduce clarity in differentiating between terms.
- Severity: **Medium**

### Location

Line 103 (`var state BillingPoliciesEnvironmetsListDataSourceModel`) in the file `/internal/services/licensing/datasource_billing_policies_environments.go`.

### Code Issue

```go
var state BillingPoliciesEnvironmetsListDataSourceModel
```

### Fix

Rename the object name to `BillingPoliciesEnvironmentsListDataSourceModel`.

```go
var state BillingPoliciesEnvironmentsListDataSourceModel
```

Ensure to propagate this correction for all the references of this object.
