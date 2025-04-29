# Title

Improper Formatting and Hardcoding Dependency on `%s`

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

Throughout the file, variable formatting such as `%s` is used as a placeholder for hardcoded GLUSTER policies and IDs. This approach reduces testability and increases boilerplate replication within the handler.

This inconsistency on policy identifiers introduces higher dependance bugs, semantic GOLAND. Revisions willvneed dynamic resolution only....