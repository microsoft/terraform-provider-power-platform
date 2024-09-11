---
page_title: "Importing existing Power Platform resources"
subcategory: "Guides"
description: |-
  <no value>
---

# Importing existing Power Platform resources

Many of the resources in the Power Platform provider can be imported into your Terraform configuration.  This guide will show you how to import existing resources into your configuration.  Resources that can be imported have an `import` section in their documentation.

## Using `terraform import`

Using `terraform import` will only import resources into the state.  It does not generate configuration and requires that you write the configuration yourself and ensure that it is correct.

```bash
terraform import powerplatform_environment.example_env 00000000-0000-0000-0000-000000000001
```

## Using import blocks

Terraform allows the use of import blocks to import existing resources into your configuration.  This can be useful if you have pre-existing resources that you want to manage with Terraform.  The following is an example of an import block used to import the default Power Platform environment:

```terraform
import {
  to = powerplatform_environment.example_env
  id = "00000000-0000-0000-0000-000000000001"
}

resource "powerplatform_environment" "default" {
  display_name     = "Contoso (default)"
  location         = "unitedstates"
  environment_type = "Default"
  dataverse = {
    currency_code     = "USD"
    language_code     = 1033
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }

  lifecycle {
    prevent_destroy = true
  }
}
```
Terraform processes the import block during the plan stage. Once a plan is approved, Terraform imports the resource into its state during the subsequent apply stage.

```bash
terraform apply
```

## Generating configuration

Terraform has [experimental support for generating configuration from import blocks](https://developer.hashicorp.com/terraform/language/import/generating-configuration).  This feature is not perfect, but it can help you get started with your configuration if you have some pre-existing resources.
