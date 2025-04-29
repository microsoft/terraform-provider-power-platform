### 2. Title
Potential Over-Matching in `UrlValidStringRegex`

##  
`/workspaces/terraform-provider-power-platform/internal/helpers/regex.go`  

## Problem  
The regex for `UrlValidStringRegex` is overly permissive. It matches strings that may not actually be valid URLs and could lead to security vulnerabilities if this regex is used for validation in sensitive contexts.  

For instance, the regex does not account for protocols (e.g., `http://` or `https://`), nor does it validate domain structures or TLDs.  

## Impact  
As it stands, the regex could allow invalid or malicious URL input, causing issues such as injection risks or erroneous behavior in systems relying on the pattern for validation. Severity: **high**.  

## Location  
```go
const (  
	UrlValidStringRegex = "(?i)^[A-Za-z0-9-._~%/:/?=]+$"  
)  
```  

## Fix  
Improve the regex to validate more accurate URL patterns, or use a full URL-parser library for increased reliability. For example:  

```go
const (  
	// Improved regex for stricter URL validation.
	UrlValidStringRegex = "^(https?://)?([\\da-z.-]+)\\.([a-z.]{2,6})([/\\w .-]*)*/?$"  
)
```  

### Explanation:  
This improved regex checks for optional protocols (`http://` or `https://`), validates domain names, and ensures a proper structure for URLs.