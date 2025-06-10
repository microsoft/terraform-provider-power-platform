# Title

Shadowing the standard library package `net/url` with alias `neturl`

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The import statement `neturl "net/url"` unnecessarily aliases the `net/url` package as `neturl`. This is not idiomatic and may lead to confusion, as all Go code uses `url.Parse()`, not `neturl.Parse()`, unless there is a naming conflict—which is not evident here.

## Impact

Lowers readability, makes it harder for new contributors familiar with Go's standard library. Severity: **low**

## Location

Imports at file top and throughout code wherever `neturl.Parse` is used.

## Code Issue

```go
import (
    //...
    neturl "net/url"
    //...
)

// ...
if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
    // ...
}
```

## Fix

Remove the alias and use the standard package name:

```go
import (
    //...
    "net/url"
    //...
)

// ...
if u, e := url.Parse(url); e != nil || !u.IsAbs() {
    // ...
}
```
