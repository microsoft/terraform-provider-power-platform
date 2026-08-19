<!-- markdownlint-disable-file -->
# Current State Research: `powerplatform_locations`, `powerplatform_languages`, `powerplatform_currencies`

## Research Topics and Questions

1. Full Terraform attribute inventory for each of the three data sources, including nested list-item attributes, required/optional/computed, Go DTO field, and exact BAPI JSON path.
2. The exact BAPI call each data source makes: method, URL template, api-version, query params, request/response DTO struct names, real fixture JSON.
3. Provider conventions used by these three: `helpers.TypeInfo` embedding, `EnterRequestContext`, verbatim `MarkdownDescription` text, httpmock unit-test structure, fixture folder naming, synthetic `id` computation.
4. Cross-dependency blast radius: does `powerplatform_environment` (resource) or any validator depend on these clients/DTOs at create/validate time?

Companion document (new API target shapes): `.copilot-tracking/research/subagents/2026-08-18/powerplatform-api-environmentmanagement.md`

---

## 1. File Inventory

### `internal/services/locations/`

| File | Purpose |
| --- | --- |
| internal/services/locations/datasource_locations.go | Data source: factory, `DataSource` struct, Metadata, Schema, Configure, Read |
| internal/services/locations/models.go | `DataSourceModel`, `DataModel` (tfsdk tags) |
| internal/services/locations/dto.go | `locationDto`, `locationsArrayDto`, `locationProperties` |
| internal/services/locations/api_locations.go | `newLocationsClient`, `client`, `GetLocations` |
| internal/services/locations/datasource_locations_test.go | `TestAccLocationsDataSource_Validate_Read`, `TestUnitLocationsDataSource_Validate_Read` |
| internal/services/locations/tests/datasource/Validate_Read/get_locations.json | BAPI fixture, 18 locations |

### `internal/services/languages/`

| File | Purpose |
| --- | --- |
| internal/services/languages/datasource_languages.go | Data source: factory, Metadata, Schema, Configure, Read |
| internal/services/languages/models.go | **`DataSource` struct lives here (not in the datasource file)**, `DataSourceModel`, `DataModel` |
| internal/services/languages/dto.go | `languagesArrayDto`, `languageDto`, `languagePropertiesDto` |
| internal/services/languages/api_languages.go | `newLanguagesClient`, `client`, `GetLanguagesByLocation` |
| internal/services/languages/datasource_languages_test.go | `TestAccLanguagesDataSource_Validate_Read`, `TestUnitLanguagesDataSource_Validate_Read` |
| internal/services/languages/tests/datasource/Validate_Read/get_languages.json | BAPI fixture, 45 languages |

### `internal/services/currencies/`

| File | Purpose |
| --- | --- |
| internal/services/currencies/datasource_currencies.go | Data source: factory, `DataSource` struct, Metadata, Schema, Configure, Read |
| internal/services/currencies/models.go | `DataSourceModel`, `DataModel` |
| internal/services/currencies/dto.go | `currenciesDto`, `currenciesArrayDto`, `currenciesPropertiesDto` |
| internal/services/currencies/api_currencies.go | `newCurrenciesClient`, `client`, `GetCurrenciesByLocation` |
| internal/services/currencies/datasource_currencies_test.go | `TestAccCurrenciesDataSource_Validate_Read`, `TestUnitCurrenciesDataSource_Validate_Read` |
| internal/services/currencies/tests/datasource/Validate_Read/get_currencies.json | BAPI fixture, 112 currencies |

Registration:

* internal/provider/provider.go:458 — `func() datasource.DataSource { return locations.NewLocationsDataSource() }`
* internal/provider/provider.go:459 — `languages.NewLanguagesDataSource()`
* internal/provider/provider.go:460 — `currencies.NewCurrenciesDataSource()`
* internal/provider/provider_test.go:68-70 — same three asserted in the expected data source list.

---

## 2. Attribute Inventory — `powerplatform_locations`

Schema function: internal/services/locations/datasource_locations.go:49-104
Read/mapping loop: internal/services/locations/datasource_locations.go:145-157

| Terraform attribute | TF type | R/O/C | Go model field (models.go) | Go DTO field (dto.go) | BAPI JSON path | Schema line |
| --- | --- | --- | --- | --- | --- | --- |
| `timeouts` | Object (`timeouts.Attributes`, `Read` only) | Optional | `DataSourceModel.Timeouts` (models.go:9) | n/a | n/a | datasource_locations.go:56 |
| `locations` | List of Object | Computed | `DataSourceModel.Value` `tfsdk:"locations"` (models.go:10) | `locationDto.Value` (dto.go:7) | `value` | datasource_locations.go:58 |
| `locations[].id` | String | Computed | `DataModel.ID` (models.go:14) | `locationsArrayDto.ID` (dto.go:11) | `value[].id` | datasource_locations.go:63 |
| `locations[].name` | String | Computed | `DataModel.Name` (models.go:15) | `locationsArrayDto.Name` (dto.go:13) | `value[].name` | datasource_locations.go:67 |
| `locations[].display_name` | String | Computed | `DataModel.DisplayName` (models.go:16) | `locationProperties.DisplayName` (dto.go:18) | `value[].properties.displayName` | datasource_locations.go:71 |
| `locations[].code` | String | Computed | `DataModel.Code` (models.go:17) | `locationProperties.Code` (dto.go:19) | `value[].properties.code` | datasource_locations.go:75 |
| `locations[].is_default` | Bool | Computed | `DataModel.IsDefault` (models.go:18) | `locationProperties.IsDefault` (dto.go:20) | `value[].properties.isDefault` | datasource_locations.go:79 |
| `locations[].is_disabled` | Bool | Computed | `DataModel.IsDisabled` (models.go:19) | `locationProperties.IsDisabled` (dto.go:21) | `value[].properties.isDisabled` | datasource_locations.go:83 |
| `locations[].can_provision_database` | Bool | Computed | `DataModel.CanProvisionDatabase` (models.go:20) | `locationProperties.CanProvisionDatabase` (dto.go:22) | `value[].properties.canProvisionDatabase` | datasource_locations.go:87 |
| `locations[].can_provision_customer_engagement_database` | Bool | Computed | `DataModel.CanProvisionCustomerEngagementDatabase` (models.go:21) | `locationProperties.CanProvisionCustomerEngagementDatabase` (dto.go:23) | `value[].properties.canProvisionCustomerEngagementDatabase` | datasource_locations.go:91 |
| `locations[].azure_regions` | List of String | Computed | `DataModel.AzureRegions` (models.go:22) | `locationProperties.AzureRegions` (dto.go:24) | `value[].properties.azureRegions` | datasource_locations.go:95 |

Notes:

* No filter/input attributes. The data source always returns the full tenant location list.
* `locationsArrayDto.Type` (dto.go:12, JSON `value[].type`) is deserialized but **not** exposed in the schema.
* There is no synthetic ID: `locations[].id` is the raw ARM-style BAPI resource path, and there is no top-level `id` attribute on the data source itself.

---

## 3. Attribute Inventory — `powerplatform_languages`

Schema function: internal/services/languages/datasource_languages.go:44-90
Read/mapping loop: internal/services/languages/datasource_languages.go:139-148

| Terraform attribute | TF type | R/O/C | Go model field (models.go) | Go DTO field (dto.go) | BAPI JSON path | Schema line |
| --- | --- | --- | --- | --- | --- | --- |
| `timeouts` | Object (`Read` only) | Optional | `DataSourceModel.Timeouts` (models.go:18) | n/a | n/a | datasource_languages.go:51 |
| `location` | String | **Required** (input filter) | `DataSourceModel.Location` (models.go:19) | n/a — becomes the `{location}` path segment | request path only | datasource_languages.go:53 |
| `languages` | List of Object | Computed | `DataSourceModel.Value` `tfsdk:"languages"` (models.go:20) | `languagesArrayDto.Value` (dto.go:7) | `value` | datasource_languages.go:57 |
| `languages[].name` | String | Computed | `DataModel.Name` (models.go:24) | `languageDto.Name` (dto.go:10) | `value[].name` | datasource_languages.go:62 |
| `languages[].id` | String | Computed | `DataModel.ID` (models.go:25) | `languageDto.ID` (dto.go:11) | `value[].id` | datasource_languages.go:66 |
| `languages[].display_name` | String | Computed | `DataModel.DisplayName` (models.go:26) | `languagePropertiesDto.DisplayName` (dto.go:19) | `value[].properties.displayName` | datasource_languages.go:70 |
| `languages[].localized_name` | String | Computed | `DataModel.LocalizedName` (models.go:27) | `languagePropertiesDto.LocalizedName` (dto.go:18) | `value[].properties.localizedName` | datasource_languages.go:74 |
| `languages[].locale_id` | Int64 | Computed | `DataModel.LocaleID` (models.go:28) | `languagePropertiesDto.LocaleID` (dto.go:17) | `value[].properties.localeId` | datasource_languages.go:78 |
| `languages[].is_tenant_default` | Bool | Computed | `DataModel.IsTenantDefault` (models.go:29) | `languagePropertiesDto.IsTenantDefault` (dto.go:20) | `value[].properties.isTenantDefault` | datasource_languages.go:82 |

Notes:

* `location` is `Required: true` with **no validators** (no `stringvalidator`, no cross-check against `powerplatform_locations`). An invalid location surfaces as an API error from `Api.Execute`.
* `languageDto.Type` (dto.go:12, JSON `value[].type`) is deserialized but **not** exposed.
* `state.Location = types.StringValue(state.Location.ValueString())` (datasource_languages.go:137) is a no-op round-trip that just preserves the input.

---

## 4. Attribute Inventory — `powerplatform_currencies`

Schema function: internal/services/currencies/datasource_currencies.go:49-95
Read/mapping loop: internal/services/currencies/datasource_currencies.go:135-144

| Terraform attribute | TF type | R/O/C | Go model field (models.go) | Go DTO field (dto.go) | BAPI JSON path | Schema line |
| --- | --- | --- | --- | --- | --- | --- |
| `timeouts` | Object (`Read` only) | Optional | `DataSourceModel.Timeouts` (models.go:12) | n/a | n/a | datasource_currencies.go:56 |
| `location` | String | **Required** (input filter) | `DataSourceModel.Location` (models.go:13) | n/a — becomes the `{location}` path segment | request path only | datasource_currencies.go:58 |
| `currencies` | List of Object | Computed | `DataSourceModel.Value` `tfsdk:"currencies"` (models.go:14) | `currenciesDto.Value` (dto.go:7) | `value` | datasource_currencies.go:62 |
| `currencies[].id` | String | Computed | `DataModel.ID` (models.go:18) | `currenciesArrayDto.ID` (dto.go:12) | `value[].id` | datasource_currencies.go:67 |
| `currencies[].name` | String | Computed | `DataModel.Name` (models.go:19) | `currenciesArrayDto.Name` (dto.go:11) | `value[].name` | datasource_currencies.go:71 |
| `currencies[].type` | String | Computed | `DataModel.Type` (models.go:20) | `currenciesArrayDto.Type` (dto.go:13) | `value[].type` (**top level, not under `properties`**) | datasource_currencies.go:75 |
| `currencies[].code` | String | Computed | `DataModel.Code` (models.go:21) | `currenciesPropertiesDto.Code` (dto.go:18) | `value[].properties.code` | datasource_currencies.go:79 |
| `currencies[].symbol` | String | Computed | `DataModel.Symbol` (models.go:22) | `currenciesPropertiesDto.Symbol` (dto.go:19) | `value[].properties.symbol` | datasource_currencies.go:83 |
| `currencies[].is_tenant_default` | Bool | Computed | `DataModel.IsTenantDefault` (models.go:23) | `currenciesPropertiesDto.IsTenantDefault` (dto.go:20) | `value[].properties.isTenantDefault` | datasource_currencies.go:87 |

Notes:

* `currencies` is the only one of the three that exposes the BAPI ARM `type` field as a Terraform attribute. This is a migration hazard — the new API has no `type`.
* `location` is `Required: true` with no validators.
* The `location` round-trip no-op is at datasource_currencies.go:134.

---

## 5. Exact BAPI Calls

Host is always `client.Api.GetConfig().Urls.BapiUrl`; for the public cloud this resolves to `api.bap.microsoft.com` (internal/constants/constants.go:28 `PUBLIC_BAPI_DOMAIN`).
`api-version` is always `constants.BAP_API_VERSION` = `"2023-06-01"` (internal/constants/constants.go:188).
Query-param name is `constants.API_VERSION_PARAM` = `"api-version"` (internal/constants/constants.go:137).

### 5.1 Locations

Source: internal/services/locations/api_locations.go:26-42

```http
GET https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01
```

| Item | Value |
| --- | --- |
| Method | `GET` (api_locations.go:38) |
| Host | `client.Api.GetConfig().Urls.BapiUrl` (api_locations.go:32) |
| Path | `/providers/Microsoft.BusinessAppPlatform/locations` (api_locations.go:33) |
| Query | `api-version=2023-06-01` (api_locations.go:28) |
| Request body | none |
| Response DTO | `locationDto` → `[]locationsArrayDto` → `locationProperties` (dto.go:6-25) |
| Acceptable status | `[]int{http.StatusOK}` |
| Deserialization | `Api.Execute(..., &locations)` — auto-unmarshal via the `responseObj` param (api_locations.go:38) |
| Error wrap | `fmt.Errorf("failed to get locations: %w", err)` (api_locations.go:40) |

Verbatim fixture (first 3 of 18 entries) — internal/services/locations/tests/datasource/Validate_Read/get_locations.json:1-52:

```json
{
  "value": [
    {
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates",
      "type": "Microsoft.BusinessAppPlatform/locations",
      "name": "unitedstates",
      "properties": {
        "displayName": "United States",
        "code": "NA",
        "isDefault": true,
        "isDisabled": false,
        "canProvisionDatabase": true,
        "canProvisionCustomerEngagementDatabase": true,
        "azureRegions": [
          "eastus",
          "westus"
        ]
      }
    },
    {
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/unitedstatesfirstrelease",
      "type": "Microsoft.BusinessAppPlatform/locations",
      "name": "unitedstatesfirstrelease",
      "properties": {
        "displayName": "Preview (United States)",
        "code": "NA",
        "isDefault": false,
        "isDisabled": false,
        "canProvisionDatabase": true,
        "canProvisionCustomerEngagementDatabase": true,
        "azureRegions": [
          "eastus",
          "westus"
        ]
      }
    },
    {
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/europe",
      "type": "Microsoft.BusinessAppPlatform/locations",
      "name": "europe",
      "properties": {
        "displayName": "Europe",
        "code": "EMEA",
        "isDefault": false,
        "isDisabled": false,
        "canProvisionDatabase": true,
        "canProvisionCustomerEngagementDatabase": true,
        "azureRegions": [
          "westeurope",
          "northeurope"
        ]
      }
    }
  ]
}
```

### 5.2 Languages

Source: internal/services/languages/api_languages.go:28-55

```http
GET https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/{location}/environmentLanguages?api-version=2023-06-01
```

| Item | Value |
| --- | --- |
| Method | `GET` (api_languages.go:40) |
| Host | `client.Api.GetConfig().Urls.BapiUrl` (api_languages.go:31) |
| Path | `fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentLanguages", location)` (api_languages.go:32) |
| Query | `api-version=2023-06-01` (api_languages.go:35) |
| Request body | none |
| Response DTO | `languagesArrayDto` → `[]languageDto` → `languagePropertiesDto` (dto.go:6-21) |
| Acceptable status | `[]int{http.StatusOK}` |
| Deserialization | Manual: `responseObj` is `nil`, then `json.Unmarshal(response.BodyAsBytes, &languages)` (api_languages.go:49) |
| Extra guard | Explicit empty-body check returning `errors.New("empty response body")` (api_languages.go:46-48) — **unique to languages** |

Verbatim fixture (first 3 of 45 entries) — internal/services/languages/tests/datasource/Validate_Read/get_languages.json:1-35:

```json
{
  "value": [
    {
      "name": "1033",
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages/1033",
      "type": "Microsoft.BusinessAppPlatform/locations/environmentLanguages",
      "properties": {
        "localeId": 1033,
        "localizedName": "English (United States)",
        "displayName": "English (United States)",
        "isTenantDefault": true
      }
    },
    {
      "name": "1025",
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages/1025",
      "type": "Microsoft.BusinessAppPlatform/locations/environmentLanguages",
      "properties": {
        "localeId": 1025,
        "localizedName": "العربية (المملكة العربية السعودية)",
        "displayName": "العربية (المملكة العربية السعودية)",
        "isTenantDefault": false
      }
    },
    {
      "name": "1069",
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages/1069",
      "type": "Microsoft.BusinessAppPlatform/locations/environmentLanguages",
      "properties": {
        "localeId": 1069,
        "localizedName": "euskara (euskara)",
        "displayName": "euskara (euskara)",
        "isTenantDefault": false
      }
    }
  ]
}
```

### 5.3 Currencies

Source: internal/services/currencies/api_currencies.go:27-52

```http
GET https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/{location}/environmentCurrencies?api-version=2023-06-01
```

| Item | Value |
| --- | --- |
| Method | `GET` (api_currencies.go:39) |
| Host | `client.Api.GetConfig().Urls.BapiUrl` (api_currencies.go:30) |
| Path | `fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location)` (api_currencies.go:31) |
| Query | `api-version=2023-06-01` (api_currencies.go:34) |
| Request body | none |
| Response DTO | `currenciesDto` → `[]currenciesArrayDto` → `currenciesPropertiesDto` (dto.go:6-21) |
| Acceptable status | `[]int{http.StatusOK}` |
| Deserialization | Manual: `responseObj` is `nil`, then `json.Unmarshal(response.BodyAsBytes, &currencies)` (api_currencies.go:47) |
| Error wrap | none — raw error returned |

Verbatim fixture (first 3 of 112 entries) — internal/services/currencies/tests/datasource/Validate_Read/get_currencies.json:1-32:

```json
{
  "value": [
    {
      "name": "DJF",
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies/DJF",
      "type": "Microsoft.BusinessAppPlatform/locations/environmentCurrencies",
      "properties": {
        "code": "DJF",
        "symbol": "Fdj",
        "isTenantDefault": false
      }
    },
    {
      "name": "ZAR",
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies/ZAR",
      "type": "Microsoft.BusinessAppPlatform/locations/environmentCurrencies",
      "properties": {
        "code": "ZAR",
        "symbol": "R",
        "isTenantDefault": false
      }
    },
    {
      "name": "ETB",
      "id": "/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies/ETB",
      "type": "Microsoft.BusinessAppPlatform/locations/environmentCurrencies",
      "properties": {
        "code": "ETB",
        "symbol": "ብር",
        "isTenantDefault": false
      }
    }
  ]
}
```

> Note: the BAPI `environmentCurrencies` fixture does **not** contain a `properties.name` field, contrary to the claim in `.copilot-tracking/research/subagents/2026-08-18/powerplatform-api-environmentmanagement.md:469`. The currency human-readable "name" in this provider comes from `value[].name`, which is the ARM resource name and is identical to `properties.code` in every fixture entry inspected.

---

## 6. Auth / Scope Selection Per Host

Function: internal/api/client.go:253-270 `tryGetScopeFromURL(url string, cloudConfig config.ProviderConfigUrls) (string, error)`

```go
switch {
case strings.LastIndex(url, cloudConfig.BapiUrl) != -1,
    strings.LastIndex(url, cloudConfig.PowerAppsUrl) != -1:
    return cloudConfig.PowerAppsScope, nil
case strings.LastIndex(url, cloudConfig.PowerPlatformUrl) != -1:
    return cloudConfig.PowerPlatformScope, nil
...
}
```

* All three services call `Api.Execute(ctx, nil, ...)` — `scopes` is `nil`, so the scope is inferred from the URL (internal/api/client.go:137-146).
* Today, because the host is `BapiUrl`, the token scope is `PowerAppsScope` = `https://service.powerapps.com/.default` (internal/constants/constants.go:30).
* After migration to `api.powerplatform.com`, `tryGetScopeFromURL` will match `cloudConfig.PowerPlatformUrl` and automatically switch to `PowerPlatformScope` = `https://api.powerplatform.com/.default` (internal/constants/constants.go:31-32). **No code change to the api layer is required** — passing `nil` scopes continues to work.
* Sovereign-cloud equivalents already exist for both hosts and scopes (internal/constants/constants.go:26-122), so the migration is cloud-agnostic as long as `PowerPlatformUrl`/`PowerPlatformScope` are populated for every cloud in `internal/config`.

Relevant constants:

| Constant | Value | Line |
| --- | --- | --- |
| `API_VERSION_PARAM` | `"api-version"` | internal/constants/constants.go:137 |
| `HTTPS` | `"https"` | internal/constants/constants.go:136 |
| `BAP_API_VERSION` | `"2023-06-01"` | internal/constants/constants.go:188 |
| `PUBLIC_BAPI_DOMAIN` | `"api.bap.microsoft.com"` | internal/constants/constants.go:28 |
| `PUBLIC_POWERPLATFORM_API_DOMAIN` | `"api.powerplatform.com"` | internal/constants/constants.go:31 |
| `PUBLIC_POWERPLATFORM_API_SCOPE` | `"https://api.powerplatform.com/.default"` | internal/constants/constants.go:32 |

> There is **no** existing `ENVIRONMENT_MANAGEMENT_API_VERSION` / `2024-10-01` constant in internal/constants/constants.go. The api-version block ends at internal/constants/constants.go:189 (`TENANT_SETTINGS_API_VERSION`). A new constant must be added.

---

## 7. Provider Conventions Used by These Three

### 7.1 `helpers.TypeInfo` embedding + `EnterRequestContext`

| Service | `TypeInfo` embedded in `DataSource` | `TypeName` set in factory | `EnterRequestContext` call sites |
| --- | --- | --- | --- |
| locations | internal/services/locations/datasource_locations.go:33 | `"locations"` at datasource_locations.go:26-28 | :41 (Metadata), :50 (Schema), :108 (Configure), :129 (Read) |
| languages | **internal/services/languages/models.go:13** (struct is in models.go, not the datasource file) | `"languages"` at datasource_languages.go:26-28 | :36 (Metadata), :45 (Schema), :94 (Configure), :121 (Read) |
| currencies | internal/services/currencies/datasource_currencies.go:33 | `"currencies"` at datasource_currencies.go:26-28 | :41 (Metadata), :50 (Schema), :99 (Configure), :118 (Read) |

Every call site follows the mandated pattern:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
defer exitContext()
```

`Metadata` additionally sets `d.ProviderTypeName = req.ProviderTypeName` **before** `EnterRequestContext`, then `resp.TypeName = d.FullTypeName()` and a `tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))`.

Inconsistency worth noting: `languages` adds extra `tflog.Debug` "READ DATASOURCE LANGUAGES START/END" lines (datasource_languages.go:129, 146) which the repo conventions discourage (entry/exit tracing is `EnterRequestContext`'s job). `locations` and `currencies` do not.

### 7.2 `Configure` conventions

* All three guard `req.ProviderData == nil` with the comment `// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.`
* All three type-assert to `*api.ProviderClient` and emit `"Unexpected ProviderData Type"` (locations :112, currencies :105) or `"Unexpected ProviderData type"` (languages :100 — lowercase `type`, inconsistent).
* `languages` adds an **extra nil-`Api` guard** (datasource_languages.go:107-113) emitting `"Unexpected nil Api in ProviderClient"`; locations and currencies do not.
* Client assignment: locations :125, languages :117, currencies :114.

### 7.3 Verbatim `MarkdownDescription` strings (for reuse/update)

**locations** (internal/services/locations/datasource_locations.go):

| Line | Text |
| --- | --- |
| 53 | `Fetches the list of available Dynamics 365 locations. For more information see [Power Platform Geos](https://learn.microsoft.com/power-platform/admin/regions-overview)` |
| 59 | `List of available locations` |
| 64 | `Unique identifier of the location` |
| 68 | `Name of the location` |
| 72 | `Display name of the location` |
| 76 | `Code of the location` |
| 80 | `Is the location default` |
| 84 | `Is the location disabled` |
| 88 | `Can the location provision a database` |
| 92 | `Can the location provision a customer engagement database` |
| 96 | `List of Azure regions` |

**languages** (internal/services/languages/datasource_languages.go):

| Line | Text |
| --- | --- |
| 48 | `Fetches the list of Dynamics 365 languages. For more information see [Power Platform Enable Languages](https://learn.microsoft.com/power-platform/admin/enable-languages)` |
| 54 | `Location of the languages` |
| 58 | `List of available languages` |
| 63 | `Name of the location` |
| 67 | `Unique identifier of the location` |
| 71 | `Display name of the location` |
| 75 | `Localized name of the location` |
| 79 | `Locale ID of the location` |
| 83 | `Is the location the default for the tenant` |

> Copy-paste defect: every nested `languages[]` description says "of the location" when it should say "of the language". Same for `is_tenant_default`. A migration PR is a natural place to fix these.

**currencies** (internal/services/currencies/datasource_currencies.go):

| Line | Text |
| --- | --- |
| 53 | `Fetches the list of available Dynamics 365 currencies. For more information see [Power Platform Currencies](https://learn.microsoft.com/power-platform/admin/manage-transactions-with-multiple-currencies)` |
| 59 | `Location of the currencies` |
| 63 | `List of available currencies` |
| 68 | `Unique identifier of the currency` |
| 72 | `Name of the currency` |
| 76 | `Type of the currency` |
| 80 | `Code of the location` |
| 84 | `Symbol of the currency` |
| 88 | `Is the currency the default for the tenant` |

> Copy-paste defect: `code` says "Code of the location" (should be "Code of the currency").

These strings are reproduced verbatim in the generated docs: docs/data-sources/locations.md, docs/data-sources/languages.md, docs/data-sources/currencies.md. Changing any of them requires `make userdocs`.

### 7.4 httpmock unit-test structure and fixture folder naming

Pattern (identical across all three):

```go
func TestUnit<X>DataSource_Validate_Read(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("GET", `<literal full URL incl. api-version>`,
        func(req *http.Request) (*http.Response, error) {
            return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_<x>.json").String()), nil
        })

    resource.Test(t, resource.TestCase{
        IsUnitTest:               true,
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{{ Config: `...`, Check: resource.ComposeAggregateTestCheckFunc(...) }},
    })
}
```

| Service | Test func line | Registered responder URL | Fixture path | Asserted count |
| --- | --- | --- | --- | --- |
| locations | internal/services/locations/datasource_locations_test.go:37 | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01` (:41) | `tests/datasource/Validate_Read/get_locations.json` (:43) | `locations.# == 18` |
| languages | internal/services/languages/datasource_languages_test.go:39 | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages?api-version=2023-06-01` (:43) | `tests/datasource/Validate_Read/get_languages.json` (:45) | `languages.# == 45` |
| currencies | internal/services/currencies/datasource_currencies_test.go:40 | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies?api-version=2023-06-01` (:44) | `tests/datasource/Validate_Read/get_currencies.json` (:46) | `currencies.# == 112` |

Fixture folder naming convention observed: `tests/datasource/<TestNameWithoutTestUnitPrefix>/<method>_<object>.json`, i.e. `TestUnitLocationsDataSource_Validate_Read` → `tests/datasource/Validate_Read/get_locations.json`.

Acceptance tests: `TestAccLocationsDataSource_Validate_Read` (locations_test.go:16), `TestAccLanguagesDataSource_Validate_Read` (languages_test.go:17), `TestAccCurrenciesDataSource_Validate_Read` (currencies_test.go:17). All use `mocks.TestAccProtoV6ProviderFactories` and `resource.TestMatchResourceAttr` with `helpers.StringRegex`.

> **Defect found**: `TestAccCurrenciesDataSource_Validate_Read` (internal/services/currencies/datasource_currencies_test.go:33-34) asserts `currencies.0.display_name` and `currencies.0.locale_id` — **neither attribute exists in the currencies schema**. Copy-pasted from the languages acceptance test. This test cannot pass as written; it is presumably never run in CI.

### 7.5 Synthetic identifiers

None of the three computes a synthetic ID.

* `locations[].id`, `languages[].id`, `currencies[].id` are passed straight through from the BAPI `value[].id` ARM resource path (datasource_locations.go:147, datasource_languages.go:141, datasource_currencies.go:137).
* There is no top-level `id` attribute on any of the three data sources.
* Consequence for migration: the new `api.powerplatform.com` responses have **no `id` field at all**. If `id` is kept, the provider must synthesize it (for example `unitedstates/1033`), or the attribute must be removed as a breaking change.

### 7.6 Examples

| File | Content |
| --- | --- |
| examples/data-sources/powerplatform_locations/data-source.tf | `data "powerplatform_locations" "all_locations" {}` (no arguments) |
| examples/data-sources/powerplatform_locations/outputs.tf | `output "all_locations" { value = data.powerplatform_locations.all_locations }` |
| examples/data-sources/powerplatform_languages/data-source.tf | `data "powerplatform_locations" "all_locations" {}` then `data "powerplatform_languages" "all_languages_by_location" { location = data.powerplatform_locations.all_locations.locations[0].name }` |
| examples/data-sources/powerplatform_currencies/data-source.tf | Same chained pattern with `powerplatform_currencies` |

All three use `provider "powerplatform" { use_cli = true }`.

> The chained example is important: `locations[0].name` feeds `languages.location` / `currencies.location`. The new API's `Location.name` field still exists, so this chaining pattern survives migration.

---

## 8. Cross-Dependency Check (Blast Radius)

### 8.1 Does `powerplatform_environment` use the locations/languages/currencies **clients or DTOs**?

**No — but it duplicates all three BAPI calls with its own private code.** There is zero Go-level coupling to the three service packages, but there is a hard behavioural coupling to the same three BAPI endpoints.

| Concern | Where | Detail |
| --- | --- | --- |
| Location list | internal/services/environment/api_environment.go:61-77 `(*Client).GetLocations` | Duplicate `GET https://{BapiUrl}/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01`. Uses its **own** DTOs `LocationArrayDto` / `LocationDto` / `LocationPropertiesDto` at internal/services/environment/dto.go:295-314 — structurally identical to `locations.locationDto` but exported and independently declared. |
| Location + Azure region validation | internal/services/environment/api_environment.go:79-100 `(*Client).LocationValidator` | Calls `GetLocations`, then `findLocation` (:38) and `findAzureRegion` (:52). `findAzureRegion` reads `location.Properties.AzureRegions`. |
| Currency validation | internal/services/environment/api_environment.go:119-163 `currencyCodeValidator` | Duplicate `GET .../locations/{location}/environmentCurrencies?api-version=2023-06-01`, own DTOs `currencyCodeValidatorDto` / `currencyCodeValidatorPropertiesDto` / `currencyCodeValidatorArrayDto` (:102-117). Compares against `item.Name` (the ARM resource name), **not** `properties.code`. |
| Language validation | internal/services/environment/api_environment.go:183-227 `languageCodeValidator` | Duplicate `GET .../locations/{location}/environmentLanguages?api-version=2023-06-01`, own DTOs `languageCodeValidatorDto` / `languageCodeValidatorArrayDto` / `languageCodeValidatorPropertiesDto` (:165-181). Compares against `item.Name`, **not** `properties.localeId`. Note `LocaleID` here is `int` (:177) vs `int64` in the languages service. |

Call sites in `Create` — internal/services/environment/resource_environment.go:

| Line | Call |
| --- | --- |
| 429 | `err = r.EnvironmentClient.LocationValidator(ctx, envToCreate.Location, envToCreate.Properties.AzureRegion)` → error `"Location validation failed for %s"` |
| 443 | `err = languageCodeValidator(ctx, r.EnvironmentClient.Api, envToCreate.Location, fmt.Sprintf("%d", envToCreate.Properties.LinkedEnvironmentMetadata.BaseLanguage))` → `"Language code validation failed for %s"` (guarded by `if envToCreate.Properties.LinkedEnvironmentMetadata != nil`, :441) |
| 449 | `err = currencyCodeValidator(ctx, r.EnvironmentClient.Api, envToCreate.Location, envToCreate.Properties.LinkedEnvironmentMetadata.Currency.Code)` → `"Currency code validation failed for %s"` |

**Migration implication:** these three validators are pre-flight `GET`s issued on **every** `powerplatform_environment` create. They break in exactly the same way as the data sources if the underlying endpoints change shape:

* `languageCodeValidator` matches on `item.Name` (`"1033"`). The new API returns only `collection[].localeId` (integer) — the matcher must become `strconv.FormatInt(int64(item.LocaleId), 10)` or an int comparison.
* `currencyCodeValidator` matches on `item.Name` (`"USD"`). The new API returns only `collection[].code` — matcher must switch to `item.Code`.
* `findAzureRegion` depends on `properties.azureRegions`, which **does not exist** in the new `Location` model. Migrating locations without a replacement for `azureRegions` would remove Azure-region validation for `powerplatform_environment` entirely.

### 8.2 Do any shared test mocks depend on these fixtures?

**Yes.** internal/mocks/mocks.go `ActivateEnvironmentHttpMocks()` cross-references all three service fixture files by relative path:

| Line | Responder | Fixture |
| --- | --- | --- |
| internal/mocks/mocks.go:69-72 | Regexp `^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/locations/(europe\|unitedstates)/environmentLanguages\?api-version=2023-06-01$` | `../../services/languages/tests/datasource/Validate_Read/get_languages.json` |
| internal/mocks/mocks.go:74-77 | Regexp `.../(europe\|unitedstates)/environmentCurrencies\?api-version=2023-06-01$` | `../../services/currencies/tests/datasource/Validate_Read/get_currencies.json` |
| internal/mocks/mocks.go:79-82 | Literal `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01` | `../../services/locations/tests/datasource/Validate_Read/get_locations.json` |

Also relevant: internal/mocks/mocks.go:47-49 registers a `RegisterNoResponder` that hard-fails on any unmocked request. So **any** URL change in the three services or in the environment validators will immediately break every `powerplatform_environment` unit test that calls `ActivateEnvironmentHttpMocks()`.

**This is the single largest blast-radius item**: changing the fixture *file contents* to the new API shape would break the environment tests unless the environment validators are migrated in the same change, and vice versa.

### 8.3 Other consumers

| Consumer | Coupling | File |
| --- | --- | --- |
| `powerplatform_environment_templates` | Uses the **same BAPI location-scoped path family** (`/providers/Microsoft.BusinessAppPlatform/locations/{location}/templates`) but does **not** import the three packages | internal/services/environment_templates/api_environment_templates.go:31 |
| Provider registration | Import + factory call only | internal/provider/provider.go:35, :47, :49, :458-460 |
| Provider test | Import + expected-list assertion only | internal/provider/provider_test.go:26, :38, :40, :68-70 |

No validators under `internal/validators/` reference these clients or DTOs. No other service imports `internal/services/locations`, `internal/services/languages`, or `internal/services/currencies`.

---

## 9. Old vs New API Field Diff (migration impact)

Cross-referenced against `.copilot-tracking/research/subagents/2026-08-18/powerplatform-api-environmentmanagement.md:355-500`.

### Locations

| Terraform attribute | BAPI source | New API source (`/environmentmanagement/provisioning/locations`) | Impact |
| --- | --- | --- | --- |
| `locations[].id` | `value[].id` | **missing** | Breaking — remove or synthesize |
| `locations[].name` | `value[].name` | `collection[].name` | OK |
| `locations[].display_name` | `value[].properties.displayName` | `collection[].displayName` | OK |
| `locations[].code` | `value[].properties.code` | `collection[].code` | Value change: BAPI `"NA"`/`"EMEA"`/`"APAC"` vs new API `"NAM"`/`"EUR"` per the reconstructed sample — **must confirm live** |
| `locations[].is_default` | `value[].properties.isDefault` | `collection[].isDefault` | OK |
| `locations[].is_disabled` | `value[].properties.isDisabled` | `collection[].isDisabled` | OK |
| `locations[].can_provision_database` | `value[].properties.canProvisionDatabase` | `collection[].canProvisionDatabase` | OK |
| `locations[].can_provision_customer_engagement_database` | `value[].properties.canProvisionCustomerEngagementDatabase` | **missing** | Breaking |
| `locations[].azure_regions` | `value[].properties.azureRegions` | **missing** | Breaking — also removes `findAzureRegion` validation for `powerplatform_environment` |
| — | — | `collection[].hasFirstReleaseIslandAvailableForProvisioning` (new) | Additive opportunity |
| — | — | `locationSelectionMode`, `macroRegions[]` (new, top-level) | Additive opportunity |
| Envelope | `value` | `collection` | Envelope rename |

### Languages

| Terraform attribute | BAPI source | New API source (`.../locations/{location}/languages`) | Impact |
| --- | --- | --- | --- |
| `languages[].name` | `value[].name` (`"1033"`) | **missing** | Breaking — could synthesize from `localeId` |
| `languages[].id` | `value[].id` | **missing** | Breaking |
| `languages[].display_name` | `value[].properties.displayName` | **missing** | Breaking — only `localizedName` exists |
| `languages[].localized_name` | `value[].properties.localizedName` | `collection[].localizedName` | OK |
| `languages[].locale_id` | `value[].properties.localeId` | `collection[].localeId` (int32) | OK; note int32 vs current Go `int64` |
| `languages[].is_tenant_default` | `value[].properties.isTenantDefault` | `collection[].isTenantDefault` | OK |
| Envelope | `value` | `collection` | Envelope rename |

### Currencies

| Terraform attribute | BAPI source | New API source (`.../locations/{location}/currencies`) | Impact |
| --- | --- | --- | --- |
| `currencies[].id` | `value[].id` | **missing** | Breaking |
| `currencies[].name` | `value[].name` (`"DJF"`) | **missing** | Breaking — identical to `code` in every fixture entry, so can be aliased to `code` |
| `currencies[].type` | `value[].type` | **missing** | Breaking |
| `currencies[].code` | `value[].properties.code` | `collection[].code` | OK |
| `currencies[].symbol` | `value[].properties.symbol` | `collection[].symbol` | OK |
| `currencies[].is_tenant_default` | `value[].properties.isTenantDefault` | `collection[].isTenantDefault` | OK |
| Envelope | `value` | `collection` | Envelope rename |

---

## 10. Key Discoveries

1. **Zero Go-level coupling, high behavioural coupling.** No other package imports the three services, but `internal/services/environment` re-implements all three BAPI calls privately (internal/services/environment/api_environment.go:61, :119, :183) and `internal/mocks/mocks.go:69-82` reuses the three fixture files. Migrating the data sources without migrating the environment validators leaves the provider calling both APIs.
2. **Auth requires no api-layer change.** `tryGetScopeFromURL` (internal/api/client.go:253) already routes `api.powerplatform.com` to `PowerPlatformScope` for every cloud. Because all three clients pass `nil` scopes, switching the host is sufficient.
3. **No `2024-10-01` api-version constant exists.** internal/constants/constants.go:183-190 has no `environmentmanagement` entry.
4. **Envelope rename `value` → `collection`** applies to all three.
5. **`id`/`name`/`type` disappear across the board.** Every one of the three data sources exposes at least one ARM-derived identifier that the new API does not return. This is the core breaking-change decision for the migration.
6. **`azureRegions` loss is the highest-risk gap** because it silently removes `powerplatform_environment` Azure-region validation (internal/services/environment/api_environment.go:52-59, :93-97).
7. **Two deserialization styles coexist.** `locations` uses `Api.Execute(..., &dto)`; `languages` and `currencies` pass `nil` and hand-roll `json.Unmarshal`. Only `languages` guards against an empty body.
8. **Pre-existing test defect** in `TestAccCurrenciesDataSource_Validate_Read` (internal/services/currencies/datasource_currencies_test.go:33-34): asserts `display_name` and `locale_id` on a schema that has neither.
9. **Pre-existing doc defects**: all `languages[]` nested descriptions say "of the location"; `currencies[].code` says "Code of the location".
10. **Currency `properties.name` does not exist in BAPI.** The companion research doc's gap note at line 469 is inaccurate on this point — verified against internal/services/currencies/tests/datasource/Validate_Read/get_currencies.json.

---

## 11. Clarifying Questions

1. Is the migration allowed to be a **breaking change** (removing `id`, `type`, `name`, `azure_regions`, `can_provision_customer_engagement_database`, `display_name` on languages), or must every current attribute keep a value (synthesized or deprecated-but-populated)?
2. Should the three data sources migrate **together with** `powerplatform_environment`'s private validators (internal/services/environment/api_environment.go), or in separate PRs? The shared fixtures in internal/mocks/mocks.go:69-82 make a split awkward.
3. Should `powerplatform_environment_templates` (internal/services/environment_templates/api_environment_templates.go:31) migrate in the same change? The new API has a matching `/environmentmanagement/provisioning/locations/{location}/templates`.
4. Is a real captured response (mitmproxy per devdocs/adr/mitmproxy.md) available for the three new endpoints? The companion doc's samples are reconstructed from schema, not captured, and the `code` value change (`"NA"` → `"NAM"`) needs confirmation.
5. Should the new `locationSelectionMode` / `macroRegions` / `hasFirstReleaseIslandAvailableForProvisioning` fields be surfaced as new attributes in this change or deferred?

---

## 12. References

Workspace files (plain paths, no links — consumed by agents):

* internal/services/locations/datasource_locations.go
* internal/services/locations/models.go
* internal/services/locations/dto.go
* internal/services/locations/api_locations.go
* internal/services/locations/datasource_locations_test.go
* internal/services/locations/tests/datasource/Validate_Read/get_locations.json
* internal/services/languages/datasource_languages.go
* internal/services/languages/models.go
* internal/services/languages/dto.go
* internal/services/languages/api_languages.go
* internal/services/languages/datasource_languages_test.go
* internal/services/languages/tests/datasource/Validate_Read/get_languages.json
* internal/services/currencies/datasource_currencies.go
* internal/services/currencies/models.go
* internal/services/currencies/dto.go
* internal/services/currencies/api_currencies.go
* internal/services/currencies/datasource_currencies_test.go
* internal/services/currencies/tests/datasource/Validate_Read/get_currencies.json
* internal/services/environment/api_environment.go
* internal/services/environment/dto.go
* internal/services/environment/resource_environment.go
* internal/services/environment_templates/api_environment_templates.go
* internal/mocks/mocks.go
* internal/api/client.go
* internal/constants/constants.go
* internal/provider/provider.go
* internal/provider/provider_test.go
* docs/data-sources/locations.md
* docs/data-sources/languages.md
* docs/data-sources/currencies.md
* examples/data-sources/powerplatform_locations/data-source.tf
* examples/data-sources/powerplatform_locations/outputs.tf
* examples/data-sources/powerplatform_languages/data-source.tf
* examples/data-sources/powerplatform_currencies/data-source.tf
* .copilot-tracking/research/subagents/2026-08-18/powerplatform-api-environmentmanagement.md
* devdocs/adr/httpmocks.md
* devdocs/adr/mitmproxy.md

External:

* [Get Supported Locations](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-supported-locations)
* [Get Provisioning Currencies](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-currencies)
* [Get Provisioning Languages](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-languages)
