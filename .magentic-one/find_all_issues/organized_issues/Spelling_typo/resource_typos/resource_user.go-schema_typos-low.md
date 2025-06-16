# Title

Minor code readability and docstring/description typos in schema block

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

Several docstrings and Markdown descriptions in the `Schema` method (especially for `disable_delete` and `security_roles` attributes) contain typos, markdown errors, or unclear phrasing. Example issues include:  
- Spelling: "Delte" instead of "Delete", "propertyto" instead of "property to"  
- Markdown syntax error: `(Disable Delte)[URL]` should be `[Disable Delete](URL)`  
- Unclear documentation logic on relationship of options.

## Impact

These documentation typos and doc issues may confuse users, reduce end-user confidence, and undermine the usability of provider-generated docs in the Terraform Registry. Severity: **Low**.

## Location

Within the `Schema` method, for multiple attribute docstrings, e.g.:

```go
"disable_delete": schema.BoolAttribute{
    MarkdownDescription: "Disable delete. When set to `True` is expects that (Disable Delte)[https://learn.microsoft.com/power-platform/admin/delete-users..." +
        "... If you just want to remove the resource and not delete the user from Dataverse, set this propertyto `False`\n\n" +
        ...
},
```

## Code Issue

```go
MarkdownDescription: "Disable delete. When set to `True` is expects that (Disable Delte)[https://learn.microsoft.com/power-platform/admin/delete-users?WT.mc_id=ppac_inproduct_settings#soft-delete-users-in-power-platform] feature to be enabled." +
    "Removing resource will try to delete the systemuser from Dataverse. This is the default behaviour. If you just want to remove the resource and not delete the user from Dataverse, set this propertyto `False`\n\n" +
    "**This attribute applies only when working with dataverse users.**",
```

## Fix

Correct all spelling and markdown errors for clarity and better user-facing documentation:

```go
MarkdownDescription: "Disable delete. When set to `True`, it expects that [Disable Delete](https://learn.microsoft.com/power-platform/admin/delete-users?WT.mc_id=ppac_inproduct_settings#soft-delete-users-in-power-platform) feature to be enabled." +
    "Removing the resource will try to delete the system user from Dataverse (default behavior). If you just want to remove the resource and not delete the user from Dataverse, set this property to `False`.\n\n" +
    "**This attribute applies only when working with Dataverse users.**",
```

Perform similar corrections for other documentation fields as needed.

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_user.go-schema_typos-low.md.
