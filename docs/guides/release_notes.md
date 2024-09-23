---
page_title: "Release Notes"
subcategory: "Guides"
description: |-
  <no value>
---

# Release Notes

This document lists breaking changes and major features for each release of the Power Platform Terraform provider.  See the release notes for each version for more details.

# v3.0.0

## BREAKING CHANGES

* `powerplatform_solution.solution_name` is no longer needed and has been removed.
* `powerplatform_solution.settings_file_checksum` is now generated using SHA256 instead of MD5. This will cause a change in the checksum value of existing resources.
* `powerplatform_solution.solution_file_checksum` is now generated using SHA256 instead of MD5. This will cause a change in the checksum value of existing resources.
* `powerplatform_rest.expected_http_status` type is changed from []int64 to []int.  Practically, this should not affect any existing configurations.
* `powerplatform_rest_query.expected_http_status` type is changed from []int64 to []int.  Practically, this should not affect any existing configurations.  
* `powerplatform_tenant_settings.id` is now set to the tenant id instead of a random guid.  This will cause a change in the id value of existing resources.

## Features

* Added `powerplatform_tenant` data source to retrieve information about the tenant.
* Added `powerplatform_tenant_capacity` data source to retrieve information about the tenant capacity.
* Added `powerplatform_environment_group` resource to manage environment groups.
