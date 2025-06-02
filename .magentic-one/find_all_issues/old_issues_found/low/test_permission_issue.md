# Title

Incorrect File Permissions When Creating Files/Directories for Testing

##

/workspaces/terraform-provider-power-platform/internal/helpers/hash_test.go

## Problem

The permissions for file creation and directory creation are hard-coded as `0644` for `os.WriteFile` and `os.Mkdir`. While `0644` permissions are suitable for files, `0644` is problematic for directories as it doesn't grant execute permission for traversal.

## Impact

The `os.Mkdir` call sets incorrect permissions for the directory `file5`, preventing traversal into the directory and potentially causing confusion or runtime errors for tests that work with this directory. Severity: Low

## Location

The issue is located in the setup code for the unit tests, where the `TestUnitCalculateSHA256` function is initializing test files and directories.

## Code Issue

```go
	err = os.WriteFile(file1, []byte("same"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	...

	err = os.Mkdir(file5, 0644)
	if err != nil {
		t.Fatal(err)
	}
```

## Fix

Adjust the directory creation to set permissions that include execute (`0755` is more appropriate for directories).

```go
	err = os.WriteFile(file1, []byte("same"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	...

	err = os.Mkdir(file5, 0755)
	if err != nil {
		t.Fatal(err)
	}
```

This ensures proper permissions for the directory while retaining the same file permissions for the files created.