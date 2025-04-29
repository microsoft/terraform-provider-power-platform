# Issue: Redundant Uppercase Conversion for Regions

## Problem
Strings are repeatedly converted to uppercase directly using `strings.ToUpper`. This redundancy could be reduced for a centralized implementation.

## Impact
Minor inefficiencies and potential inconsistency across different functions. **Low Severity** efficiency adjustment.

## Affected Code
```go
analyticDataUrl, exists := urlMap[strings.ToUpper(region)]
```

## Recommended Fix
Utilize a utility function like below to handle such conversions both effectively and uniformly:

### Example Utility Function
```go
func sanitizeRegion(region string) string {
	return strings.ToUpper(region)
}
```

Then call the helper wherever applicable:
```go
func getAnalyticsUrl(region string) (string, error) {
	urlMap := getAnalyticsUrlMap()
	region = sanitizeRegion(region)
	analyticDataUrl, exists := urlMap[region]
	if !exists {
		return "", fmt.Errorf("invalid region: %s", region)
	}
	return analyticDataUrl, nil
}
```