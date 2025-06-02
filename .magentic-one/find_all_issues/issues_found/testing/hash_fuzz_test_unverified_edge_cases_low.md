# Unverified Edge Case Coverage in Fuzz Seed Corpus

##

/workspaces/terraform-provider-power-platform/internal/helpers/hash_fuzz_test.go

## Problem

The fuzz test adds a comprehensive set of edge cases to the fuzz corpus, but it does not validate that all are exercised, and some paths (e.g., Windows reserved names) may not exist or be meaningful on all OSes. There is also no validation of what constitutes expected behavior (for example, specific error types for different invalid inputs).

## Impact

Severity: Low

While not a functional bug, this reduces the diagnostic value of the fuzz test and can lead to spurious or missed error reporting in cross-platform scenarios, making the test less useful for reliably surfacing bugs.

## Location

```go
	f.Add("/dev/null")               // Reserved name on Linux
	f.Add("CON")                     // Reserved name on Windows
	// ... various other OS-dependent paths ...
```

## Code Issue

```go
	f.Add("/dev/null")               // Reserved name on Linux
	f.Add("CON")                     // Reserved name on Windows
	f.Add(string(make([]byte, 300))) // Extremely long path
	f.Add("../relative/path")
	f.Add("./current/dir")
	f.Add(" ")                    // Single space
	f.Add("\n")                   // Newline character
	f.Add("Z:/nonexistent/drive") // Nonexistent drive on Windows
	f.Add("//network/share")
	f.Add("\\\\network\\share")
	f.Add("/dev/random")  // Special device file on Linux
	f.Add("/dev/urandom") // Special device file on Linux
```

## Fix

Structure the test to conditionally include OS-specific paths using build constraints (`runtime.GOOS`), and assert specific error types for known edge cases to clarify and improve diagnostic output. Document expected behaviors for platform-dependent seeds.

```go
import (
    //...
    "runtime"
)

// In FuzzCalculateSHA256:
if runtime.GOOS == "windows" {
    f.Add("CON")                     // Reserved name on Windows
    f.Add("Z:/nonexistent/drive")    // Nonexistent drive on Windows
    // etc.
} else {
    f.Add("/dev/null")               // Reserved name on Linux
    f.Add("/dev/random")
    f.Add("/dev/urandom")
    // etc.
}

// In the fuzz function, optionally improve error assertion:
if filePath == "CON" && runtime.GOOS == "windows" {
    if err == nil {
        t.Errorf("Expected error for reserved Windows name 'CON', but got nil")
    }
}
// Similar checks for each special case...
```
