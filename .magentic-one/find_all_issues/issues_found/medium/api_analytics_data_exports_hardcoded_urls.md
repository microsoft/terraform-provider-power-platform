# Issue: Hardcoded URLs Inside Function

## Problem
Hardcoding URLs inside the `getAnalyticsUrlMap` function can lead to maintenance difficulties and hard-to-manage configurations across multiple environments. The URLs should be configurable rather than hardcoded in the code itself.

## Impact
Environment-specific changes or additions can lead to redeployment when URLs need to be changed. This could influence operational efficiency. The impact is **medium**, as it can be resolved but requires better coding practice.

## Affected Code
```go
func getAnalyticsUrlMap() map[string]string {
	return map[string]string{
		"US":   "https://na.csanalytics.powerplatform.microsoft.com/",
		"CAN":  "https://can.csanalytics.powerplatform.microsoft.com/",
		"EMEA": "https://emea.csanalytics.powerplatform.microsoft.com/",
	}
}
```

## Recommended Fix
Move these configurations into a separate configuration file or use environment variables:

### Example Solution (Environment Variables):
```go
func getAnalyticsUrlMap() map[string]string {
	return map[string]string{
		"US":   os.Getenv("US_ANALYTICS_URL"),
		"CAN":  os.Getenv("CAN_ANALYTICS_URL"),
		"EMEA": os.Getenv("EMEA_ANALYTICS_URL"),
	}
}
```
And set these environment variables appropriately for different environments.

### Example Solution (Configuration File):
Use a configuration manager to load these dynamically at runtime.
