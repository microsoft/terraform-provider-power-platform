<!-- markdownlint-disable-file -->
# Research: `powerplatform_environment` attribute VALUE-level dependencies (BAPI -> Power Platform API)

## Research Topics / Questions

1. Hardcoded value slices in `dto.go` (`EnvironmentTypes`, `CadenceTypes`, `ReleaseCycleTypes`, others) — exact strings.
2. `models.go` value derivations: `administration_mode_enabled`, `background_operation_enabled`, `release_cycle`, `enterprise_policies[].type`, `location`/`azure_region` normalization.
3. `api_environment.go` validators: `LocationValidator`, `currencyCodeValidator`, `languageCodeValidator` — matched field + error text.
4. `resource_environment.go` schema (lines 63-377): Required/Optional/Computed, `stringvalidator.OneOf` lists, plan modifiers.
5. Literal values in the richest BAPI environment JSON fixtures.
6. Literal values in locations / languages / currencies fixtures (`code` = "NA" or "NAM"?).

## Status

Complete. All six questions answered with file:line evidence.

---

## 1. Hardcoded value literals (`internal/services/environment/dto.go`)

All constants live in one `const` block at dto.go lines 11-26 and one `var` block at lines 28-35.

| Go identifier | Exact literal string(s) | Line |
| --- | --- | --- |
| `EnvironmentTypesDeveloper` | `"Developer"` | dto.go:12 |
| `EnvironmentTypesSandbox` | `"Sandbox"` | dto.go:13 |
| `EnvironmentTypesProduction` | `"Production"` | dto.go:14 |
| `EnvironmentTypesTrial` | `"Trial"` | dto.go:15 |
| `EnvironmentTypesDefault` | `"Default"` | dto.go:16 |
| `CadenceTypesFrequent` | `"Frequent"` | dto.go:18 |
| `CadenceTypesModerate` | `"Moderate"` | dto.go:19 |
| `ReleaseCycleTypesStandard` | `"Standard"` | dto.go:21 |
| `ReleaseCycleTypesEarly` | `"Early"` | dto.go:22 |
| `ReleaseCycleFirstReleasePublicDto` | `"FirstRelease"` | dto.go:24 |
| `ReleaseCycleFirstReleaseGovDto` | `"GovFR"` | dto.go:25 |

Derived slices / regexes:

* `EnvironmentTypes = []string{"Developer", "Sandbox", "Production", "Trial", "Default"}` — dto.go:29
* `EnvironmentTypesDeveloperOnlyRegex = "^(Developer)$"` — dto.go:30
* `EnvironmentTypesExceptDeveloperRegex = "^(Sandbox|Production|Trial|Default)$"` — dto.go:31
* `CadenceTypes = []string{"Frequent", "Moderate"}` — dto.go:32
* `ReleaseCycleTypes = []string{"Standard", "Early"}` — dto.go:33
* `ReleaseCycleFirstReleaseOnlyRegex = "^(FirstRelease|GovFR)$"` — dto.go:34 (declared; no consumer found in the environment package)

Other hardcoded literals sent to the API (not in the slices above):

* `"CommonDataService"` written to `properties.databaseType` on create — models.go:165
* `"1"` written to `properties.usedBy.type` on create (developer environments) — models.go:156
  * NOTE value asymmetry: the API *returns* `"usedBy": { "type": "User" }` in fixtures, but the provider *sends* `"1"`.

Per-cloud `FirstRelease` cluster name (internal/config/config.go:79-104):

| Cloud | `FirstReleaseClusterName` value | Line |
| --- | --- | --- |
| `Public` | `"FirstRelease"` | config.go:82 |
| `Gcc` | `"GovFR"` | config.go:86 |
| `GccHigh` | `nil` | config.go:89 |
| `Dod` | `nil` | config.go:92 |
| `China` | `nil` | config.go:95 |
| `Ex` | `nil` | config.go:98 |
| `Rx` | `nil` | config.go:101 |

Key name constant: `FirstReleaseClusterName CloudTypeConfigurationKey = "release_cycle"` — config.go:28.

---

## 2. Value derivations in `internal/services/environment/models.go`

### 2.1 `dataverse.administration_mode_enabled`

Read path (models.go:276-280):

```go
if environmentDto.Properties.States != nil && environmentDto.Properties.States.Runtime != nil && environmentDto.Properties.States.Runtime.Id == "AdminMode" {
    attrValuesProductProperties["administration_mode_enabled"] = types.BoolValue(true)
} else {
    attrValuesProductProperties["administration_mode_enabled"] = types.BoolValue(false)
}
```

* Compared string: exactly `"AdminMode"` (case-sensitive `==`), read from `properties.states.runtime.id`.
* Anything else (including `"Enabled"`) yields `false`. There is no error for an unknown value — it silently becomes `false`.

Write path (resource_environment.go:769-782): sets `properties.states.runtime.id` to `"AdminMode"` when `true` and `"Enabled"` when `false`.

### 2.2 `dataverse.background_operation_enabled`

Read path (models.go:281-285):

```go
if environmentDto.Properties.LinkedEnvironmentMetadata.BackgroundOperationsState == "Enabled" {
```

* Compared string: exactly `"Enabled"` on `properties.linkedEnvironmentMetadata.backgroundOperationsState` (dto.go:149). Anything else -> `false`.

Write path (resource_environment.go:784-790): `"Enabled"` / `"Disabled"`.

### 2.3 `release_cycle`

`convertReleaseCycleModelFromDto` (models.go:361-368):

```go
value := providerConfig.GetCurrentCloudConfiguration(config.FirstReleaseClusterName)
if environmentDto.Properties.Cluster != nil && value != nil && environmentDto.Properties.Cluster.Catergory == *value {
    model.ReleaseCycle = types.StringValue(ReleaseCycleTypesEarly)   // "Early"
} else {
    model.ReleaseCycle = types.StringValue(ReleaseCycleTypesStandard) // "Standard"
}
```

* Source field: `properties.cluster.category` (dto.go:85, field misspelled `Catergory` in Go).
* Compared value: the per-cloud constant — `"FirstRelease"` (Public) or `"GovFR"` (GCC). Any other category (fixtures use `"Prod"`) maps to `"Standard"`.

Write path (models.go:123-130): if plan `release_cycle == "Early"` and the cloud constant is non-nil, send `properties.cluster.category = <constant>`. `"Standard"` sends no cluster at all.

### 2.4 `enterprise_policies[].type`

`convertEnterprisePolicyModelFromDto` (models.go:384-458). The BAPI key -> synthesized Terraform `type` string mapping is:

| BAPI JSON key (dto.go) | Terraform `type` output | models.go line |
| --- | --- | --- |
| `identity` (dto.go:97) | `"Identity"` | models.go:406 |
| `vnets` (dto.go:95) | `"NetworkInjection"` | models.go:426 |
| `customerManagedKeys` (dto.go:96) | `"Encryption"` | models.go:446 |

Other sub-fields are passed through verbatim: `id` <- `id`, `location` <- `location`, `system_id` <- `systemId`, `status` <- `linkStatus` (dto.go:100-106).

Structural bug worth carrying into the migration: each branch calls `types.SetValueMust(...)` and **overwrites** `model.EnterprisePolicies` rather than appending, so an environment with more than one policy only reports the last one evaluated (order: identity, then vnets, then customerManagedKeys). When `properties.enterprisePolicies` is absent the set is empty (models.go:456).

Casing note: the DTO tag is `vnets` (dto.go:95) but the recorded BAPI payload uses `vNets` (enterprise_policy fixture line 119). This only works because Go `encoding/json` matches field names case-insensitively.

### 2.5 `location` / `azure_region` normalization

* No normalization. `Location` and `AzureRegion` are passed straight through: models.go:219-220 (`types.StringValue(environmentDto.Location)` and `types.StringValue(environmentDto.Properties.AzureRegion)`); create path models.go:94 and models.go:116-118.
* The only `strings.ToLower` in the whole environment package is resource_environment.go:928, which lowercases an *environment id* for a licensing check — not location.
* Consequence: value casing from the API must match the practitioner's configured `location` string byte-for-byte, or every plan shows a diff (and `location` is `RequiresReplace`).

### 2.6 Other value-sensitive derivations in models.go

| Attribute | Rule | Line |
| --- | --- | --- |
| `dataverse.security_group_id` | forced to `null` and `owner_id` populated when `properties.environmentSku == "Developer"` | models.go:272-275 |
| `billing_policy_id` | `""` or missing `billingPolicy` -> `"00000000-0000-0000-0000-000000000000"` (`constants.ZERO_UUID`) | models.go:353-359 |
| `environment_group_id` | missing `parentEnvironmentGroup` -> `constants.ZERO_UUID` | models.go:345-351 |
| `cadence` | direct pass-through of `properties.updateCadence.id`; **dereferenced without a nil check** | models.go:225 |
| `allow_moving_data_across_regions` | `properties.copilotPolicies.crossGeoCopilotDataMovementEnabled` (nil -> `false`) | models.go:380 |
| `allow_flex_routing` | `properties.copilotPolicies.crossBoundaryCopilotDataMovementEnabled` (nil -> `false`) | models.go:381 |
| `dataverse.currency_code` | NOT from the environment payload; from Dataverse `transactioncurrencies.isocurrencycode` | api_environment.go:780-797, resource_environment.go:573 |
| `dataverse.templates` | read from `properties.linkedEnvironmentMetadata.template` (singular), but written on create as `templates` (plural) | dto.go:156 vs dto.go:230 |

---

## 3. Validators in `internal/services/environment/api_environment.go`

| Validator | Matched field (exact) | Source | Error message (exact format string) |
| --- | --- | --- | --- |
| `findLocation` (used by `LocationValidator`) | `loc.Name` — i.e. `value[].name`, NOT `properties.code` or `properties.displayName` | `GET https://{bapi}/providers/Microsoft.BusinessAppPlatform/locations` (api_environment.go:61-76) | `"location '%s' is not valid. valid locations are: %s"` (api_environment.go:49) — the list is built from `loc.Name` |
| `findAzureRegion` | membership in `location.Properties.AzureRegions` (`properties.azureRegions[]`) | same locations payload | `"region '%s' is not valid for location %s. valid regions are: %s"` (api_environment.go:58) |
| `LocationValidator` | orchestrates both; empty `azureRegion` short-circuits to success | — | api_environment.go:79-99 (early return at :90) |
| `currencyCodeValidator` | `item.Name` — i.e. `value[].name`, NOT `properties.code` | `GET /providers/Microsoft.BusinessAppPlatform/locations/{location}/environmentCurrencies` (api_environment.go:119-127) | `"currency Code %s is not valid. valid currency codes are: %s"` (api_environment.go:159) |
| `languageCodeValidator` | `item.Name` — i.e. `value[].name` (a **string** like `"1033"`), NOT `properties.localeId` (int) | `GET /providers/Microsoft.BusinessAppPlatform/locations/{location}/environmentLanguages` (api_environment.go:183-191) | `"language Code %s is not valid. valid language codes are: %s"` (api_environment.go:223) |

Call sites (resource_environment.go): `LocationValidator` at :431, `languageCodeValidator` at :445 (passes `fmt.Sprintf("%d", BaseLanguage)` — int formatted to string, so `value[].name` must be the decimal LCID as a string), `currencyCodeValidator` at :451.

Other exact-string state machines in api_environment.go (not "validators" but equally value-coupled):

| Compared literal | Field | Line |
| --- | --- | --- |
| `"Failed"` | `lifecycleResponse.State.Id` | :349, :455, :659, :730 |
| `"Succeeded"` | `lifecycleResponse.State.Id` | :518 |
| `"Succeeded"` / `"LinkedDatabaseProvisioning"` | `properties.provisioningState` | :419, :421, :533 |
| `"Ready"` / `"Running"` | `properties.states.management.id` (else -> hard error `"environment update failed. unexpected management state: ..."`) | :751-756 |
| `"OperationNotStartable"` (substring of the 409 body) | conflict-retry gate | :678 |

---

## 4. `resource_environment.go` schema (lines 63-377)

Top-level attributes:

| Attribute | R/O/C | `OneOf` list | Plan modifiers | Default | Lines |
| --- | --- | --- | --- | --- | --- |
| `timeouts` | — | — | — | — | :106-111 |
| `id` | Computed | — | `UseStateForUnknown` | — | :111-117 |
| `environment_group_id` | Optional + Computed | — | `UseStateForUnknown` | (commented out) | :118-132 |
| `description` | Optional + Computed | — | `UseStateForUnknown` | `""` | :133-141 |
| `cadence` | Optional + Computed | `CadenceTypes` = `Frequent`, `Moderate` | `UseStateForUnknown` | `"Moderate"` | :142-153 |
| `release_cycle` | Optional + Computed | `ReleaseCycleTypes` = `Standard`, `Early` | `RequiresReplace`, `UseStateForUnknown` | — | :154-165 |
| `location` | **Required** | — | `RequiresReplace` | — | :166-172 |
| `azure_region` | Optional + Computed (`Required: false`) | — | `RequiresReplace`, `UseStateForUnknown` | — | :173-182 |
| `environment_type` | **Required** | `EnvironmentTypes` = `Developer`, `Sandbox`, `Production`, `Trial`, `Default` | — | — | :183-190 |
| `owner_id` | Optional | — | `UseStateForUnknown`, `RequiresReplace` | — | :191-203 |
| `allow_bing_search` | Optional + Computed | — | — | — | :204-208 |
| `allow_microsoft_365_services` | Optional + Computed | — | — | — | :209-213 |
| `allow_moving_data_across_regions` | Optional + Computed | — | — | — | :214-218 |
| `allow_flex_routing` | Optional + Computed | — | — | — | :219-223 |
| `display_name` | **Required** | — | — | — | :224-230 |
| `billing_policy_id` | Optional + Computed | — | `UseStateForUnknown` | `"00000000-0000-0000-0000-000000000000"` | :231-243 |
| `enterprise_policies` | Computed | — | `setplanmodifier.UseStateForUnknown` | — | :244-253 |
| `dataverse` | Optional + Computed | — | `modifiers.RequireReplaceObjectToEmptyModifier()` | — | :254-261 |

`enterprise_policies` nested object (`policyAttributeSchema`, :69-89) — `type`, `id`, `location`, `system_id`, `status` are all Computed-only, no validators.

`dataverse` nested attributes:

| Attribute | R/O/C | Plan modifiers | Default | Lines |
| --- | --- | --- | --- | --- |
| `unique_name` | Computed | `UseStateForUnknownKeepNonNullStateModifier` | — | :262-268 |
| `administration_mode_enabled` | Optional + Computed | `boolplanmodifier.UseStateForUnknown` | `false` | :269-277 |
| `background_operation_enabled` | Optional + Computed | `boolplanmodifier.UseStateForUnknown` | `true` | :278-286 |
| `currency_code` | **Required** | `RequireReplaceStringFromNonEmptyPlanModifier` | — | :287-293 |
| `url` | Computed | `UseStateForUnknownKeepNonNullStateModifier`, `SetStringAttributeUnknownOnlyIfSecondAttributeChange(dataverse.domain)` | — | :294-301 |
| `domain` | Optional + Computed | `UseStateForUnknownKeepNonNullStateModifier` | — | :302-313 |
| `organization_id` | Computed | `UseStateForUnknownKeepNonNullStateModifier` | — | :314-320 |
| `security_group_id` | Optional + Computed | — | — | :321-330 |
| `language_code` | **Required** (Int64) | `RequireReplaceIntAttributePlanModifier` | — | :330-341 |
| `version` | Computed | `UseStateForUnknownKeepNonNullStateModifier` | — | :337-344 |
| `templates` | Optional + Computed (List of String) | — | — | :344-349 |
| `template_metadata` | Optional only | — | — | :350-353 |
| `linked_app_type` | Computed | `UseStateForUnknownKeepNonNullStateModifier` | — | :354-360 |
| `linked_app_id` | Computed | `UseStateForUnknownKeepNonNullStateModifier` | — | :361-367 |
| `linked_app_url` | Computed | `UseStateForUnknownKeepNonNullStateModifier` | — | :368-374 |

`stringvalidator.OneOf` appears exactly twice: `cadence` (:150) and `environment_type` (:187). `release_cycle` uses `OneOf` at :159. No other enum validation exists — every other API enum flows through unchecked.

Cross-field validators worth noting:

* `environment_group_id`: `LengthAtLeast(1)`, `RegexMatches(helpers.GuidRegex)`, `AlsoRequires(dataverse)` (:125-129).
* `environment_type`: `OtherFieldRequiredWhenValueOf(owner_id, nil, ^(Developer)$, ...)` (:188).
* `owner_id`: `ConflictsWith(dataverse.security_group_id)`, `AlsoRequires(dataverse)`, `OtherFieldRequiredWhenValueOf(environment_type, ^(Developer)$, ...)` (:197-201).
* `dataverse.security_group_id`: `RegexMatches(helpers.GuidRegex)`, `MakeFieldRequiredWhenOtherFieldDoesNotHaveValue(environment_type, ^(Sandbox|Production|Trial|Default)$, ...)` (:325-328).
* `dataverse.domain`: `LengthAtLeast(1)`, `RegexMatches(helpers.DomainNameRegex)` (:308-311).
* `ConfigValidators`: `RequiredTogether(dataverse.administration_mode_enabled, dataverse.background_operation_enabled)` (:382-388).

---

## 5. Literal values in the richest BAPI environment fixtures

Primary reference fixture (Sandbox + Dataverse + billing policy + env group): `internal/services/environment/tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json`

| JSON path | Literal value | File:line |
| --- | --- | --- |
| `location` | `"europe"` | Validate_Read/get_environment_...001.json:4 |
| `location` (alt) | `"unitedstates"` | Validate_Read/get_environment_...002.json:4 |
| `location` (alt) | `"switzerland"` | resource/Validate_Create_Environment_And_Dataverse/get_lifecycle_new_dataverse.json:4 |
| `properties.azureRegion` | `"westeurope"` | Validate_Read/get_environment_...001.json:8 |
| `properties.azureRegion` (alt) | `"northeurope"` | Validate_Read/get_environments.json:10 |
| `properties.environmentSku` | `"Sandbox"` | Validate_Read/get_environment_...001.json:68 |
| `properties.environmentSku` (alt) | `"Developer"` | resource/Validate_Create_Dev_Env/get_environment_...001.json:30; Validate_Read/get_environments.json:74 |
| `properties.environmentSku` (alt) | `"Production"` | resource/Validate_Update_Environment_Type/get_environment_8.json:25 |
| `properties.cluster.category` | `"Prod"` | Validate_Read/get_environment_...001.json:170 |
| `properties.cluster.category` (early release) | `"FirstRelease"` | resource/Validate_Create_Early_Release_Cycle/get_environment_...001.json:125 |
| `properties.cluster` siblings | `"number": "107"`, `"uriSuffix": "eu-il107.gateway.prod.island"`, `"geoShortName": "EU"`, `"environment": "Prod"` | Validate_Read/get_environment_...001.json:171-174 |
| `properties.states.runtime.id` | `"Enabled"` | Validate_Read/get_environment_...001.json:154 |
| `properties.states.management.id` | `"Ready"` | Validate_Read/get_environment_...001.json:145 |
| `properties.linkedEnvironmentMetadata.backgroundOperationsState` | `"Enabled"` | Validate_Read/get_environment_...001.json:131 |
| `properties.linkedEnvironmentMetadata.baseLanguage` | `1033` (JSON number, not string) | Validate_Read/get_environment_...001.json:128 |
| `properties.linkedEnvironmentMetadata.domainName` | `"00000000-0000-0000-0000-000000000001"` | Validate_Read/get_environment_...001.json:122 |
| `properties.linkedEnvironmentMetadata.uniqueName` | `"00000000-0000-0000-0000-000000000001"` | Validate_Read/get_environment_...001.json:121 |
| `properties.linkedEnvironmentMetadata.instanceUrl` | `"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/"` (trailing slash) | Validate_Read/get_environment_...001.json:124 |
| `properties.linkedEnvironmentMetadata.resourceId` | `"orgid"` | Validate_Read/get_environment_...001.json:119 |
| `properties.linkedEnvironmentMetadata.version` | `"9.2.23092.00206"` | Validate_Read/get_environment_...001.json:123 |
| `properties.linkedEnvironmentMetadata.template` | **ABSENT from every environment fixture in the repo** (grep for `"template"` across `internal/services/environment/tests/**` returns zero hits, including `Validate_Create_With_D365_Template`) | — |
| `properties.updateCadence.id` | `"Moderate"` | Validate_Read/get_environment_...001.json:158 |
| `properties.updateCadence.id` (alt) | `"Frequent"` | resource/Validate_Create_Dev_Env/get_environment_...001.json:117 |
| `properties.databaseType` | `"CommonDataService"` | Validate_Read/get_environment_...001.json:117 |
| `properties.provisioningState` | `"Succeeded"` | Validate_Read/get_environment_...001.json:66 |
| `properties.governanceConfiguration.protectionLevel` | `"Basic"` | Validate_Read/get_environment_...001.json:200 |
| `properties.usedBy` | `{ "id": "f99f844b-ce3b-49ae-86f3-e374ecae789c", "type": "User", "tenantId": "00000000-0000-0000-0000-000000000002", "userPrincipalName": "admin" }` | Validate_Read/get_environments.json:22-27 |
| `properties.usedBy` (dev env) | `{ "id": "00000000-0000-0000-0000-000000000001", "tenantId": "123", "type": "User" }` | resource/Validate_Create_Dev_Env/get_environment_...001.json:19-22 |
| `properties.billingPolicy.id` (unset case) | `""` (empty string, not absent) | resource/Validate_Create_Dev_Env/get_environment_...001.json:24 |
| `properties.parentEnvironmentGroup.id` | `"00000000-0000-0000-0000-000000000001"` | Validate_Read/get_environment_...001.json:13 |

### `properties.enterprisePolicies.*`

No environment-service fixture contains `enterprisePolicies`. The only recordings are in `internal/services/enterprise_policy/tests/Validate_Create/`:

| File | JSON key | Sub-fields |
| --- | --- | --- |
| `get_environment_00000000-0000-0000-0000-000000000001_identity.json:118-125` | `"identity"` | `policyId: "00000000-0000-0000-0000-000000000001"`, `location: "europe"`, `id: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.PowerPlatform/enterprisePolicies/bar"`, `systemId: "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000001"`, `linkStatus: "Linked"` |
| `get_environment_00000000-0000-0000-0000-000000000001_network_injection.json:118-125` | `"vNets"` (capital N in the payload; DTO tag is `vnets`) | `policyId: "00000000-0000-0000-0000-000000000002"`, `location: "europe"`, `id: ".../enterprisePolicies/bar"`, `systemId: "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"`, `linkStatus: "Linked"` |
| `get_environment_00000000-0000-0000-0000-000000000001_encryption.json:118-126` | `"customerManagedKeys"` | `policyId: "00000000-0000-0000-0000-000000000003"`, `location: "europe"`, `id: ".../enterprisePolicies/bar"`, `systemId: "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000003"`, `linkStatus: "Linked"` |

Only `"Linked"` is ever recorded for `linkStatus`. The provider passes it through verbatim as `enterprise_policies[].status`, so it never fails on an unexpected value — but the value appears in state.

---

## 6. Locations / languages / currencies fixtures

`internal/services/locations/tests/datasource/Validate_Read/get_locations.json`

| `name` | `properties.displayName` | `properties.code` | `properties.azureRegions` | Line |
| --- | --- | --- | --- | --- |
| `"unitedstates"` | `"United States"` | `"NA"` | `["eastus", "westus"]` | :6-17 |
| `"unitedstatesfirstrelease"` | `"Preview (United States)"` | `"NA"` | `["eastus", "westus"]` | :23-34 |
| `"europe"` | `"Europe"` | `"EMEA"` | `["westeurope", "northeurope"]` | :40-50 |
| `"asia"` | `"Asia"` | `"APAC"` | `["eastasia", "southeastasia"]` | :57-66 |
| `"australia"` | `"Australia"` | `"OCE"` (`isDisabled: true`) | — | :72-80 |

**Answer to the "NA vs NAM" question: the recorded BAPI value is `"NA"`, not `"NAM"`.** Note that `properties.code` is NOT what `findLocation` matches on — it matches `name` (api_environment.go:39). `code` is only surfaced by the `powerplatform_locations` data source.

`internal/services/languages/tests/datasource/Validate_Read/get_languages.json`

* `value[].name` is a **decimal LCID as a string**: `"1033"` (:4), `"1025"` (:15), `"1069"` (:26), `"1026"` (:37), `"1027"` (:47), `"3076"` (:58).
* `properties.localeId` is the same value as a JSON **number** (`1033`, :9).
* `properties.localizedName` / `properties.displayName` are localized human strings, e.g. `"English (United States)"` (:10-11).
* `properties.isTenantDefault` — `true` for `1033` (:12).

`internal/services/currencies/tests/datasource/Validate_Read/get_currencies.json`

* `value[].name` is the ISO 4217 code: `"DJF"` (:4), `"ZAR"` (:14), `"ETB"` (:24), `"AED"` (:34), `"BHD"` (:44).
* `properties.code` duplicates `name` (`"DJF"` at :8, etc.).
* `properties.symbol`: `"Fdj"`, `"R"`, `"ብር"`, `"د.إ."`, `"د.ب."`.

---

## MAIN DELIVERABLE: attribute -> value dependency matrix

| Terraform attribute | Current BAPI literal value(s) in fixtures | Code depending on the exact value (file:line) | What must be true of the new Power Platform API value |
| --- | --- | --- | --- |
| `id` | `properties`-sibling `name` = `"00000000-0000-0000-0000-000000000001"` | models.go:218 (`environmentDto.Name`, NOT `.Id`) | New API must expose a bare GUID in the field mapped to `Name`. If the new API only returns a full ARM-style `id` path, `Id` becomes the path and every state/import breaks. |
| `location` | `"europe"`, `"unitedstates"`, `"switzerland"` | models.go:94 (send), models.go:221 (read), api_environment.go:39 (`loc.Name` match) | Must be the exact same lowercase slug in both the environment payload and the locations list `name` field. No case folding exists; `RequiresReplace` means any casing drift forces destroy/recreate. |
| `azure_region` | `"westeurope"`, `"northeurope"`, `"eastus"`, `"westus"` | models.go:116-118, models.go:222, api_environment.go:53 (membership in `properties.azureRegions`) | Values in the environment payload must be drawn from the same vocabulary as the locations list `azureRegions[]`, byte-identical. |
| `environment_type` | `"Sandbox"`, `"Developer"`, `"Production"` | dto.go:29 (`OneOf`), models.go:224, models.go:272 (`== "Developer"` gate) | `environmentSku`-equivalent must return exactly `Developer`/`Sandbox`/`Production`/`Trial`/`Default`. Any new SKU string (e.g. `"Teams"`, `"Subscription"`) makes `OneOf` reject configs and silently mis-handles the developer `security_group_id`/`owner_id` special case. |
| `cadence` | `"Moderate"`, `"Frequent"` | dto.go:32 (`OneOf`), models.go:225 | The field must be present on EVERY environment: models.go:225 dereferences `Properties.UpdateCadence.Id` **without a nil check**. A missing `updateCadence` panics the provider. Values must remain exactly `Frequent`/`Moderate`. |
| `release_cycle` | `properties.cluster.category` = `"Prod"` (standard) / `"FirstRelease"` (early); GCC = `"GovFR"` | models.go:361-368, config.go:82/86, models.go:123-130 | The category field must still carry `"FirstRelease"` (public) / `"GovFR"` (GCC) for early-release environments, and something *other than* those for standard. If the new API replaces cluster category with a first-class `releaseCycle`/`updateChannel` enum, the comparison must be rewritten, not just re-pathed. |
| `dataverse.administration_mode_enabled` | `properties.states.runtime.id` = `"Enabled"` | models.go:276 (`== "AdminMode"`), resource_environment.go:772/778 (writes `"AdminMode"`/`"Enabled"`) | The runtime state must remain a **string enum** whose admin-mode member is exactly `"AdminMode"` and whose normal member is exactly `"Enabled"`. A boolean `adminMode` field, or renamed members (`"AdministrationMode"`, `"Disabled"`), silently reports `false` on read and silently no-ops on write. |
| `dataverse.background_operation_enabled` | `"Enabled"` | models.go:281, resource_environment.go:786/788 | Must remain the string pair `"Enabled"`/`"Disabled"`, not a boolean. A boolean unmarshals into the `string` DTO field as a JSON type error or empty string -> attribute always `false`. |
| `dataverse.language_code` | `1033` (number) | models.go:264 (`int64(BaseLanguage)`), resource_environment.go:445 (formats `%d` and matches against languages `value[].name`) | Environment payload must return the LCID as a JSON **number**; the languages list must return the same LCID as a **decimal string** in `name`. If the new languages endpoint returns `"en-US"` in `name`, every create fails validation. |
| `dataverse.currency_code` | `"USD"`-style ISO code from Dataverse `isocurrencycode`; validated against currencies `value[].name` = `"DJF"`, `"ZAR"`, `"AED"` | api_environment.go:780-797, resource_environment.go:573, api_environment.go:145-160 | Currencies list `name` must remain the ISO 4217 code (not a display name, not a GUID). Read path is unaffected by the migration (it queries Dataverse OData directly), so a mismatch between the new list endpoint and Dataverse's `isocurrencycode` would cause perpetual diffs. |
| `dataverse.domain` | `"00000000-0000-0000-0000-000000000001"` | models.go:262, resource_environment.go:308-311 (regex `helpers.DomainNameRegex`) | Returned domain must satisfy "starts with lowercase letter or digit, only lowercase letters/digits/`-`". Any uppercase or FQDN form (`foo.crm4.dynamics.com`) breaks the validator on subsequent applies. |
| `dataverse.url` | `"https://<guid>.crm4.dynamics.com/"` (trailing slash) | models.go:261; api_environment.go:237 `strings.TrimSuffix(..., "/")` | Trailing-slash handling only exists in `GetEnvironmentHostById`. The state value keeps the trailing slash verbatim — a new API returning it without the slash is a state-diff for existing users. |
| `dataverse.unique_name` | `"00000000-0000-0000-0000-000000000001"` | models.go:265 | Free-form pass-through; only needs to be stable across reads (it is `Computed` with `UseStateForUnknownKeepNonNullStateModifier`, so a value change shows a diff). |
| `dataverse.organization_id` | `"orgid"` (from `resourceId`) | models.go:263 | Pass-through of `linkedEnvironmentMetadata.resourceId`. Must remain the org GUID, not an ARM resource id. |
| `dataverse.version` | `"9.2.23092.00206"` | models.go:264 | Pass-through, no parsing. |
| `dataverse.security_group_id` | absent in fixtures (empty) | models.go:260, models.go:273 (nulled for Developer), api_environment.go:278 (`== ""` gate), resource_environment.go:326 (GUID regex) | Must be a GUID string or empty. Empty security group is represented as `"00000000-0000-0000-0000-000000000000"` per the schema doc; the API must accept and echo that literal. |
| `dataverse.templates` | **never present in any fixture** | dto.go:156 read tag `template` (singular) vs dto.go:230 create tag `templates` (plural); models.go:288-306 | This is the single most fragile mapping: read and write use different JSON key names. A live capture is required to determine which key the new API uses on read. |
| `dataverse.template_metadata` | never present in fixtures | dto.go:157, models.go:308-321 (JSON round-trip through `createTemplateMetadataDto.PostProvisioningPackages`) | Must round-trip with a `PostProvisioningPackages` array (capital P) — models.go:190-196 unmarshals the practitioner string with that exact casing. |
| `dataverse.linked_app_type` / `_id` / `_url` | absent in environment fixtures | models.go:251-258 (falls back to `""` not null when `linkedAppMetadata` is nil) | Nil -> empty string is the existing behaviour; a new API returning `null` sub-objects keeps this working, but returning an empty object `{}` would set all three to `""` too (same result). |
| `enterprise_policies[].type` | keys `"identity"`, `"vNets"`, `"customerManagedKeys"` | models.go:406 -> `"Identity"`, models.go:426 -> `"NetworkInjection"`, models.go:446 -> `"Encryption"` | The new API must still discriminate policies by these three object keys (case-insensitive match is available). If it returns a flat array with a `kind`/`type` discriminator, the synthesis logic must be rewritten and the output strings `Identity`/`NetworkInjection`/`Encryption` preserved to avoid a breaking state change. |
| `enterprise_policies[].status` | `"Linked"` | models.go:410/430/450 (pass-through of `linkStatus`) | Pass-through; value change is a visible state diff, not an error. |
| `enterprise_policies[].id` / `system_id` / `location` | ARM path strings + `"europe"` | models.go:407-409 etc. | Pass-through. |
| `billing_policy_id` | `"00000000-0000-0000-0000-000000000001"`; empty case is `""` | models.go:353-359, resource_environment.go:242 (default `"00000000-0000-0000-0000-000000000000"`) | The API must represent "no billing policy" as either an absent `billingPolicy` object or `id: ""`. A literal all-zero GUID from the API also works. Any other sentinel (e.g. `"None"`) becomes a permanent diff against the `ZERO_UUID` default. |
| `environment_group_id` | `"00000000-0000-0000-0000-000000000001"` | models.go:345-351, api_environment.go:282 (`== ""` gate) | Absent `parentEnvironmentGroup` must mean "no group"; the provider maps it to `ZERO_UUID`. Sending `ZERO_UUID` on create is deliberately suppressed (models.go:131-137) — the new API must preserve that "empty guid = remove from group" update semantic. |
| `owner_id` | `properties.usedBy.id` = `"f99f844b-..."`; `usedBy.type` returned as `"User"` | models.go:371-377 (read), models.go:150-159 (write sends `Type: "1"`) | Read must expose the owner object id. Write asymmetry (`"1"` sent vs `"User"` returned) must be re-verified: the new API may reject `"1"` or require the string `"User"`. |
| `allow_bing_search` | `bingChatEnabled` (absent in most fixtures) | models.go:222, dto.go:66 (`bool` with `omitempty`) | Must be a JSON boolean. Because the DTO uses `omitempty` on a non-pointer `bool`, `false` is indistinguishable from absent on write — the AI-features write path (resource_environment.go:717-742) works around this with `*bool`. |
| `allow_microsoft_365_services` | `m365Enabled` | models.go:223, dto.go:67 | Same as above. |
| `allow_moving_data_across_regions` | `copilotPolicies.crossGeoCopilotDataMovementEnabled` (requires `$expand=properties/copilotPolicies` on list — api_environment.go:768) | models.go:380 | Must be a nullable boolean; the new API must return it on a plain GET, or an equivalent expand parameter must be identified. Absent -> `false` silently. |
| `allow_flex_routing` | `copilotPolicies.crossBoundaryCopilotDataMovementEnabled` | models.go:381, resource_environment.go:731-740 | Same. Note the write is suppressed outside the EU data boundary because the service rejects the field entirely — the new API's rejection error text/behaviour needs re-verification. |
| `description` | `"aaa"` | models.go:217, dto.go:59 (tag has **no** `omitempty`) | Empty description must round-trip as `""`; the schema default is `""` (resource_environment.go:140). |
| `display_name` | `"displayname"`, `"Admin AdminOnMicrosoft's Environment"` | models.go:220 | Pass-through. |
| (internal) provisioning gate | `"Succeeded"`, `"LinkedDatabaseProvisioning"` | api_environment.go:419-421, :533 | The new API's long-running-operation terminal/intermediate states must use these exact strings, or every create/update hangs or errors. |
| (internal) management state gate | `"Ready"`, `"Running"` | api_environment.go:751-756 | Any third value produces a hard error `"environment update failed. unexpected management state: <value>"`. |
| (internal) lifecycle state gate | `"Failed"`, `"Succeeded"` | api_environment.go:349, :455, :518, :659, :730 | Same coupling for the lifecycle/operation polling endpoint. |
| (internal) conflict retry | body substring `"OperationNotStartable"` | api_environment.go:678 | The new API's 409 body must carry the same error code, or delete/update retries stop working. |
| (internal) database type | `"CommonDataService"` | models.go:165 (write), fixture :117 (read) | The new API must accept the same literal to provision Dataverse, or expose an equivalent flag. |

---

## Key discoveries / risks

1. **`cadence` is a nil-deref waiting to happen.** models.go:225 does `environmentDto.Properties.UpdateCadence.Id` with no nil guard. Every existing fixture includes `updateCadence`. If the new API omits it for some SKUs, the provider panics.
2. **`administration_mode_enabled` and `background_operation_enabled` are string-enum-derived booleans.** Both compare against exactly one string and default to `false` otherwise. A boolean-typed replacement in the new API is the highest-risk silent-failure scenario.
3. **`templates` read/write key mismatch** (`template` vs `templates`, dto.go:156 vs dto.go:230) is unverified against any fixture — zero recordings exist.
4. **`enterprise_policies` overwrite bug** (models.go:404-455): multiple policies collapse to one. Also, the `vnets` DTO tag only matches the recorded `vNets` payload thanks to Go's case-insensitive JSON matching — a strict re-implementation would break it.
5. **`usedBy.type` asymmetry**: provider sends `"1"` (models.go:156), API returns `"User"` (fixtures). Unverified whether the new API accepts `"1"`.
6. **Locations `code` is `"NA"`, not `"NAM"`** — but `code` is not used for validation at all; `name` is (api_environment.go:39).
7. **Languages list `name` is a decimal LCID string** (`"1033"`), matched against `fmt.Sprintf("%d", baseLanguage)` (resource_environment.go:445). Any locale-tag format change breaks creates.
8. **No case normalization anywhere** for `location`/`azure_region`, both `RequiresReplace`.
9. **`release_cycle` is cloud-conditional**: GCC High / DoD / China / Ex / Rx return `nil` (config.go:89-101), so `release_cycle` is always `"Standard"` there and `"Early"` sends no cluster at all.
10. **List endpoint uses `$expand=properties/billingPolicy,properties/copilotPolicies`** (api_environment.go:768). Whether the new API supports an equivalent expand determines if `allow_moving_data_across_regions` / `allow_flex_routing` / `billing_policy_id` are populated on the data source.

## Clarifying questions requiring input or a live call

* Does the new API expose the environment GUID in a field distinct from the ARM-style `id` path (the provider maps `id` from `Name`, models.go:218)?
* Are the lifecycle/operation polling semantics (`Failed`/`Succeeded` state ids, `Location`/`Retry-After` headers, `OperationNotStartable` 409 code) preserved under `environmentmanagement`?
* Is there a supported `$expand` (or equivalent `select`) for copilot policies and billing policy on the new list endpoint?
