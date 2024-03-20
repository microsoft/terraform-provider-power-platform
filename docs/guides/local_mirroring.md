---
page_title: "Installation: Local mirroring for the Power Platform Provider"
description: |-
  <no value>
---

# Local mirroring for the Power Platform Provider

Local mirroring allows use of a local installation of the provider instead of downloading it from the Terraform Registry.  This allows using private builds or releases from the repo which may not be published in the Terraform Registry.

## Installation

To use the provider you can download the binaries from [Releases](https://github.com/microsoft/terraform-provider-power-platform/releases) to your local file system and configure Terraform to use your local mirror.  See the [Explicit Installation Method Configuration](https://developer.hashicorp.com/terraform/cli/config/config-file#explicit-installation-method-configuration) for more information about using local binaries.

```terraform
provider_installation {
  filesystem_mirror {
    path    = "/usr/share/terraform/providers"
    include = ["registry.terraform.io/microsoft/power-platform"]
  }
}
```

## Example

An example of how to use the provider with local mirroring is available in the [Quickstarts bootstrap](https://github.com/microsoft/power-platform-terraform-quickstarts/tree/main/bootstrap/mirror) directory.
