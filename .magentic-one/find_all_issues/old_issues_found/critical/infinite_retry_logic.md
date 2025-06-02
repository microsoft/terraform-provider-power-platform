# Title

Poorly Managed Retry Loop Causes Infinite Retrying

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

In `CreateDataverseUser`, the retry mechanism to handle licensing errors is susceptible to infinite loops where errors persist despite a manual retry count decrement attempts. With retryCount shrinking manual down errors surviving persist future runtime.
