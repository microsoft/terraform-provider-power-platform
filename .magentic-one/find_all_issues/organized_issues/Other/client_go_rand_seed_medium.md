# Title

Use of `rand.Intn` without explicitly setting a seed

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The function `DefaultRetryAfter` uses `rand.Intn` to generate a random wait duration but does not explicitly set a seed for the pseudo-random number generator. This can lead to repeatable, predictable sequences of values across runs if `math/rand.Seed()` is never called at program startup.

## Impact

May reduce randomness and reliability in production, especially for retry/jitter logic. Severity: **medium**

## Location

Function `DefaultRetryAfter`

## Code Issue

```go
func DefaultRetryAfter() time.Duration {
	return time.Duration((rand.Intn(10) + 10)) * time.Second
}
```

## Fix

Ensure that the random number generator is seeded, ideally once at program/package startup. For example, in an `init()` function or main, add:

```go
import "math/rand"
import "time"

func init() {
	rand.Seed(time.Now().UnixNano())
}
```

No change is needed for the function itself apart from ensuring this is present at startup.
