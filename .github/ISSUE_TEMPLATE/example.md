---
name: "Quickstart Example Request"
about: "Request a new quickstart example to be added to the Terraform provider repository."
labels: [guide,triage]
assignees: ""

---

## Description

<!-- Short description here describing the new quickstart example that you're requesting.  Include a use case for why users need this example. -->

### Example Details

- Proposed Name: (e.g. 101-my-example)
- Supporting documentation: <!-- links to product documentation (if public). -->
- Estimated complexity/effort: <!--  (e.g., easy, moderate, hard) -->
- Related resources/data sources: <!-- what data sources and/or resources will this example use? -->

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

- [ ] Example in the /examples/quickstarts folder
- [ ] Example documentation in README.md.tmpl
- [ ] Updated auto-generated provider docs with `make quickstarts`
- [ ] Confirmation that you have manually tested this

## Contributions

Do you plan to raise a PR to address this issue?

- [ ] Yes
- [ ] No

See the [contributing guide](/CONTRIBUTING.md?) for more information about what's expected for contributions.
