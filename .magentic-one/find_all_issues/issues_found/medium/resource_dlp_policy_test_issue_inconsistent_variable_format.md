# Title

Hardcoded Variables Impacting Test Maintainability

## Path

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

The test suite using variables such as `policyId` in multiple tests is hardcoded, which reduces flexibility when scaling or debugging. Besides, identifications remain *human-readable* but don't verify tests in composable. Instead of example ready GUID repetition `policyId`, the test could leverage dynamic and assertion-ng.

## Impractical.

Fault lineage paths impact RESOURCE ATTRs scale roadmap Blocks newer imports Terraform Modules integrity quickly falls Blocks. Suitability issue Creation **impact fix territories Writable, Terraform-user maintainance long-term critical, Red sea builders Ex TechScores bad refactor contributors gaps**
Specific CASCADES policyId, display/specifiers combination layoutentication tests manually while versatile param flows fail Recycling Causes inconsistencies team ** lower production testing Modules Users suffer consequence modules. *needed-managed failed.Hard*

}}