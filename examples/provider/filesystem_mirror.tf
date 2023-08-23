provider_installation  {
  filesystem_mirror {
    path = "/usr/share/terraform/providers"
    include = ["registry.terraform.io/microsoft/power-platform"]
  }
}
