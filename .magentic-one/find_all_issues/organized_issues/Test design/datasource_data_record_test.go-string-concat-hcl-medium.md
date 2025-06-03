# Title

Improper Use of String Concatenation for HCL Generation

##

internal/services/data_record/datasource_data_record_test.go

## Problem

The `BootstrapDataRecordTest` function (and potentially others) builds large Terraform HCL configuration strings using raw string literals with interspersed string concatenation and heavy indentation. This construction technique is error-prone and hard to maintain, especially as the configuration grows or changes over time.

## Impact

- Increased risk of missing delimiters, unmatched parentheses, or copy-paste errors.
- Harder to spot logic/data issues (because the HCL is mixed with code).
- Making changes is cumbersome, e.g., if more resources or variables are needed.
- Negatively impacts readability and maintainability.

Severity: Low to Medium (low for small files, medium as config/test complexity grows).

## Location

Function: `BootstrapDataRecordTest(name string) string`
And similar large HCL strings in test steps.

## Code Issue

```go
func BootstrapDataRecordTest(name string) string {
    return `
resource "powerplatform_environment" "data_env" {
    display_name     = "` + name + `"
    location         = "unitedstates"
    ...
}`
}
```

## Fix

Use Go's `text/template` or `fmt.Sprintf` for composable HCL configuration generation. For instance:

```go
import "text/template"
...
tpl := `resource \"powerplatform_environment\" \"data_env\" {
    display_name     = \"{{ .Name }}\"
    ...
}
`
var sb strings.Builder
t := template.Must(template.New("hcl").Parse(tpl))
_ = t.Execute(&sb, struct{ Name string }{name})
return sb.String()
```

Or, if using only a few fields, `fmt.Sprintf` can suffice:

```go
return fmt.Sprintf(`
resource \"powerplatform_environment\" \"data_env\" {
    display_name     = \"%s\"
    ...
}`, name)
```

This approach improves readability, maintainability, and reduces the risk of subtle syntax bugs.

Save as a code structure/maintainability issue.
