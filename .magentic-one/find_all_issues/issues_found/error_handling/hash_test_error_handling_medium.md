# Use of t.Fatal in Parallel Subtests

##

/workspaces/terraform-provider-power-platform/internal/helpers/hash_test.go

## Problem

The test file uses `t.Fatal` within parallel subtests (i.e., subtests started with `t.Parallel()`). While Go allows this, it can make debugging difficult as the goroutine that runs the parallel subtest immediately stops, and resources may not get cleaned up properly. It is also generally encouraged to use `t.Errorf` or `t.FailNow` only when necessary and for assertions, while using table-driven tests and helper functions for repetitive patterns.

## Impact

Medium severity. Using `t.Fatal` inside parallel subtests may cause abrupt terminations, which could hinder debugging and cleanup. It makes test failures less predictable and harder to diagnose when running tests concurrently.

## Location

Multiple subtests inside `TestUnitCalculateSHA256`:

```go
t.Run("TestUnitCalculateSHA256_SameFile", func(t *testing.T) {
    t.Parallel()
    ...
    if err != nil {
        t.Fatal(err)
    }
    ...
})
```

## Code Issue

```go
if err != nil {
    t.Fatal(err)
}
```

## Fix

Replace `t.Fatal` with `t.Errorf` (to continue running as much of the test as possible), or if you want to stop the subtest immediately, use `t.FailNow`. However, if the error is truly fatal to the subtest, `t.Fatal` is acceptable; just be aware of the limitations.

For more robust testing, handle errors with explicit error messages for better diagnostics. Optionally, factor repetitive error handling into a helper.

```go
if err != nil {
    t.Fatalf("failed to calculate SHA256: %v", err)
}
```

Or use a helper if repeating:

```go
func mustSHA256(t *testing.T, filename string) string {
    t.Helper()
    val, err := helpers.CalculateSHA256(filename)
    if err != nil {
        t.Fatalf("CalculateSHA256(%s) error: %v", filename, err)
    }
    return val
}
```

Then in the test:

```go
f1 := mustSHA256(t, file1)
```
