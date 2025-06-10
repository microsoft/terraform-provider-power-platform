# Use of Magic Strings for Test Data and URLs

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages_test.go

## Problem

The test uses hard-coded magic strings for configuration blocks, API URLs, and repeated attribute names. This reduces maintainability and increases the chance of typos or inconsistencies if these values must change.

## Impact

Refactoring and updates may become error-prone or time-consuming. This is a low severity, maintainability issue.

## Location

```go
Config: `
	data "powerplatform_languages" "all_languages_for_unitedstates" {
		location = "unitedstates"
	}`,

Check: resource.ComposeAggregateTestCheckFunc(
	resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.#", regexp.MustCompile(`^[1-9]\d*$`)),
	...
	resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.id", regexp.MustCompile(helpers.StringRegex)),
	...
)
```

## Code Issue

```go
Config: `
	data "powerplatform_languages" "all_languages_for_unitedstates" {
		location = "unitedstates"
	}`,
// and
resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.0.id", regexp.MustCompile(helpers.StringRegex)),
```

## Fix

Define constants at the top of the test file:

```go
const (
	testLocation        = "unitedstates"
	testLanguagesDSName = "all_languages_for_unitedstates"
	testConfig          = `
data "powerplatform_languages" "` + testLanguagesDSName + `" {
	location = "` + testLocation + `"
}`
	testAPIEndpoint = "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/" + testLocation + "/environmentLanguages?api-version=2023-06-01"
)
```

Use these constants in the test body:

```go
Config: testConfig,
resource.TestMatchResourceAttr("data.powerplatform_languages."+testLanguagesDSName, "languages.0.id", regexp.MustCompile(helpers.StringRegex)),
```
