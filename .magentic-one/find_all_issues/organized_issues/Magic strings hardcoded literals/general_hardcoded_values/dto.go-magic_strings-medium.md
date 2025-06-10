# Magic Strings and Constants

## 
/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/dto.go

## Problem

There are numerous raw, repeated string literals, especially for keys (e.g., `"ai_generative_settings"`, `"backup_retention"`, `"sharing_controls"`, etc.) and rule identifiers scattered throughout the code. These are "magic strings" and present maintainability and typo risks. Using consts/enums would make refactoring and verification easier, and avoids mistakes due to typos or changes.

## Impact

- **Severity:** Medium
- Risk of silent bugs from typos.
- Harder code refactoring, auditing, and documentation.
- Poor discoverability and easier to miss updates across the codebase.

## Location

Widespread, for example in:

```go
aiGenerativeSettingsObj := attrs["ai_generative_settings"]
...
backupRetentionObj := attrs["backup_retention"]
...
solutionCheckerObj := attrs["solution_checker_enforcement"]
...
makerWelcomeContentObj := attrs["maker_welcome_content"]
...
rule := environmentGroupRuleSetParameterDto{
    ...
    Type:             AI_GENERATED_DESC,
    ...
}
```

## Code Issue

```go
aiGenerativeSettingsObj := attrs["ai_generative_settings"]
```

## Fix

Define constants at the top of the file (or a shared package) for all such keys and identifiers:

```go
const (
    AttrAiGenerativeSettings = "ai_generative_settings"
    AttrBackupRetention      = "backup_retention"
    AttrSharingControls      = "sharing_controls"
    // ... etc.

    TypeAiGeneratedDesc      = "AI_GENERATED_DESC"
    TypeBackupRetention      = "BACKUP_RETENTION"
    // ... etc.
)
...
aiGenerativeSettingsObj := attrs[AttrAiGenerativeSettings]
...
```

This improves IDE support and reduces the chance of copy-paste/typo bugs as well as facilitates any renaming.

---
