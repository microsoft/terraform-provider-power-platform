# Title

Confusing/inconsistent type naming and capitalization ("Dataverse" vs "dataverse" vs "environment")

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

Throughout the code, the usage of "Dataverse", "dataverse", "environment", and their related constants/variables is inconsistent in comments, log messages, and documentation. Example: sometimes "Dataverse" is capitalized, sometimes not; "environment user" and "dataverse user" logic, comments, or attributes are intermingled.

## Impact

This reduces readability for future maintainers and can cause misunderstanding for users reading documentation or debugging logs. Cloud provider resources should use consistent naming for clarity and maintainability. Severity: **Low**.

## Location

- Schema attribute descriptions (e.g., "**This attribute applies only when working with dataverse users.**")
- Log and error messages: `tflog.Debug(ctx, fmt.Sprintf("Dataverse exists in environment: %t", hasEnvDataverse))`
- Plan/model variable names and logic

## Code Issue

```go
MarkdownDescription: "Security roles Ids assigned to the Dataverse user" +
    "When working with non Dataverse environments, only 'Environment Admin' and 'Environment Maker' role values are allowed",
```

And elsewhere in code or logs, e.g.:

```go
tflog.Debug(ctx, fmt.Sprintf("Dataverse exists in environment: %t", hasEnvDataverse))
```

## Fix

Audit all user-facing documentation/comments and log messages so "Dataverse" (as a proper noun/platform feature) is consistently capitalized, and "environment" is used for non-Dataverse scenarios.  
For example, standardize as:

- "Dataverse" (capital D, everywhere referring to the Microsoft Dataverse product)
- "environment user" for non-Dataverse users  
- Ensure matching capitalization in schema Markdown descriptions and code comments/logs

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_user.go-naming_consistency-low.md.
