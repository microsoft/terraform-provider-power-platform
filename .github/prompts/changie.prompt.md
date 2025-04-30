# Changelog entries

## 1. Prerequisites

Before running changie, validate that all changes are staged by running `git status`.  
If there are unstaged changes you may offer to stage them using `git add`

After files are staged, you can run:

```bash
git -P diff --merge-base origin/main
```

Use the diff to see what's changed and use that to write a one line description for changie.

## 2. Creating a changelog entry

Create a changelog entry using:

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed (see [copilot-commit-message-instructions.md](../copilot-commit-message-instructions.md))
- `<issue_number>` is the GitHub issue number from the URL provided by the user (just the number, not the full URL)
