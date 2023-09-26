---
name: "New Data Source Request"
about: "Request a new data source to be added to the Terraform provider."
labels: [enhancement, data source, triage]
assignees: ""

---

## Description

<!-- Short description here describing the new data source that you're requesting.  Include a use case for why users need this data source. -->

### Resource

- Resource Name: powerplatform_[your data source name]
- API documentation: <!-- links to API documentation (if public).  What APIs are needed for read/list data? -->
- Estimated complexity/effort: <!--  (e.g., easy, moderate, hard) -->
- Related resources/data sources: <!-- are there any existing or potential data sources that are related to this one -->

### Potential Terraform Configuration

```hcl
# Sample Terraform config that describes how the new resource might look.

data "powerplatform_[your data source name]" "example_data_source" {
  name = "example"
  parameter1 = "value1"
  parameter2 = "value2"
}

```

## Definition of Done

- [ ] Data Transfer Objects (dtos)
- [ ] Data Client functions
- [ ] Resource Implementation
- [ ] Resource Added to Provider
- [ ] Unit Tests for Happy Path
- [ ] Unit Tests for error path
- [ ] Acceptance Tests
- [ ] Example in the /examples folder
- [ ] Schema Documentation in code
- [ ] Updated auto-generated provider docs with `make docs`

## Contributions

Do you plan to raise a PR to address this issue?

- [ ] Yes
- [ ] No

See the [contributing guide](/CONTRIBUTING.md?) for more information about what's expected for contributions.
