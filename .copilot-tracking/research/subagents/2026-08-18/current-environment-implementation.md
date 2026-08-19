<!-- markdownlint-disable-file -->
# Subagent Research: Current `powerplatform_environment` / `powerplatform_environments` Implementation (BAPI)

Research date: 2026-08-18
Repository: c:\tmp\terraform-provider-power-platform
Purpose: Establish the exact "as-is" surface of the environment resource + data source before migrating from BAPI (`api.bap.microsoft.com`) to the public Power Platform API (`api.powerplatform.com/environmentmanagement/...`).

Companion documents:

* .copilot-tracking/research/subagents/2026-08-18/powerplatform-api-environmentmanagement.md (the target API surface)
* .copilot-tracking/research/2026-08-18/powerplatform-api-migration-research.md (overall migration research)

## Research questions

1. What is the complete Terraform schema of `powerplatform_environment` (every attribute, type, required/optional/computed, plan modifiers, validators, DTO/JSON mapping, populating endpoint, source line)? — ANSWERED (section A)
2. Same for `powerplatform_environments` data source? — ANSWERED (section B)
3. What is the complete inventory of HTTP calls made by the environment service (method, URL, api-version, query params, request/response DTOs, async polling)? — ANSWERED (section C)
4. What exactly does `Update` change in place today, and which call handles each change? — ANSWERED (section D)
5. Which provider conventions must the migration preserve? — ANSWERED (section E)
6. Which services already call `api.powerplatform.com`, and what pattern do they use? — ANSWERED (section F)
7. What is the blast radius — which other services consume `environment.Client`? — ANSWERED (section G)

## Source file map

| File | Lines | Role |
| --- | --- | --- |
| internal/services/environment/resource_environment.go | 899 | Resource: schema, CRUD, update helpers |
| internal/services/environment/datasource_environments.go | ~300 | Data source: schema + Read |
| internal/services/environment/models.go | ~430 | TF models + `convert*` functions |
| internal/services/environment/dto.go | 314 | All DTOs + enum constants |
| internal/services/environment/api_environment.go | 757 | HTTP client (`Client`) — all BAPI calls |
| internal/services/environment/resource_environment_test.go | ~3690 | 59 unit + acceptance tests |
| internal/services/environment/datasource_environments_test.go | ~140 | 2 tests |
| internal/services/environment/tests/ | 149 JSON files | httpmock fixtures |

---

## A. `powerplatform_environment` resource — full attribute inventory

Schema is defined in `Schema()` at internal/services/environment/resource_environment.go:63-377.

Legend for Req/Opt/Comp: `R` = Required, `O` = Optional, `C` = Computed.

### A.1 Top-level attributes

| # | TF attribute | TF type | R/O/C | Plan modifiers | Validators | Go DTO field → JSON path | Populated by | Line |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `timeouts` | Object (timeouts) | O | — | — | n/a (framework-only, `SourceModel.Timeouts`, models.go:39) | n/a | resource_environment.go:103-108 |
| 2 | `id` | String | C | `stringplanmodifier.UseStateForUnknown()` | — | `EnvironmentDto.Name` → `name` (dto.go:41) | GET environment / GET environments | 109-115 |
| 3 | `environment_group_id` | String | O + C | `stringplanmodifier.UseStateForUnknown()` | `LengthAtLeast(1)`, `RegexMatches(helpers.GuidRegex)`, `AlsoRequires(dataverse)` | `ParentEnvironmentGroupDto.Id` → `properties.parentEnvironmentGroup.id` (dto.go:60,108-110) | GET environment; defaults to `ZERO_UUID` when absent (models.go:340-346, api_environment.go:270-272) | 116-130 |
| 4 | `description` | String | O + C, `Default: ""` | `stringplanmodifier.UseStateForUnknown()` | — | `EnviromentPropertiesDto.Description` → `properties.description` (dto.go:58) | GET environment | 131-139 |
| 5 | `cadence` | String | O + C, `Default: "Moderate"` | `stringplanmodifier.UseStateForUnknown()` | `OneOf(CadenceTypes)` = `Frequent`,`Moderate` (dto.go:18-19,33) | `UpdateCadenceDto.Id` → `properties.updateCadence.id` (dto.go:59,112-114) | GET environment | 140-151 |
| 6 | `release_cycle` | String | O + C | `RequiresReplace()`, `UseStateForUnknown()` | `OneOf(ReleaseCycleTypes)` = `Standard`,`Early` (dto.go:21-22,34) | `ClusterDto.Catergory` → `properties.cluster.category` (dto.go:62,84-86) | GET environment; derived by comparing to `config.GetCurrentCloudConfiguration(FirstReleaseClusterName)` → `"FirstRelease"` (public) / `"GovFR"` (gcc) (models.go:356-363, config.go:79-104) | 152-163 |
| 7 | `location` | String | **R** | `RequiresReplace()` | — | `EnvironmentDto.Location` → `location` (dto.go:40) | GET environment; validated against GET locations | 164-170 |
| 8 | `azure_region` | String | O + C (`Required: false` explicit) | `RequiresReplace()`, `UseStateForUnknown()` | — | `EnviromentPropertiesDto.AzureRegion` → `properties.azureRegion` (dto.go:46) | GET environment; validated against `LocationPropertiesDto.AzureRegions` (dto.go:313) | 171-180 |
| 9 | `environment_type` | String | **R** | — | `OneOf(EnvironmentTypes)` = `Developer`,`Sandbox`,`Production`,`Trial`,`Default` (dto.go:11-15,30); `validators.OtherFieldRequiredWhenValueOf(owner_id, nil, /^(Developer)$/)` | `EnviromentPropertiesDto.EnvironmentSku` → `properties.environmentSku` (dto.go:49) | GET environment | 181-188 |
| 10 | `owner_id` | String | O (no Computed) | `UseStateForUnknown()`, `RequiresReplace()` | `ConflictsWith(dataverse.security_group_id)`, `AlsoRequires(dataverse)`, `OtherFieldRequiredWhenValueOf(environment_type, /^(Developer)$/, nil)` | `UsedByDto.Id` → `properties.usedBy.id` (dto.go:63,88-92). Create also sends `type:"1"` and `tenantId` | GET environment (models.go:365-371); tenantId comes from GET tenant | 189-201 |
| 11 | `allow_bing_search` | Bool | O + C | — | — | `EnviromentPropertiesDto.BingChatEnabled` → `properties.bingChatEnabled` (dto.go:64) | GET environment | 202-206 |
| 12 | `allow_microsoft_365_services` | Bool | O + C | — | — | `EnviromentPropertiesDto.M365Enabled` → `properties.m365Enabled` (dto.go:65) | GET environment | 207-211 |
| 13 | `allow_moving_data_across_regions` | Bool | O + C | — | — | `CopilotPoliciesDto.CrossGeoCopilotDataMovementEnabled` → `properties.copilotPolicies.crossGeoCopilotDataMovementEnabled` (dto.go:66,80) | GET environment **with `$expand=properties/copilotPolicies`** | 212-216 |
| 14 | `allow_flex_routing` | Bool | O + C | — | — | `CopilotPoliciesDto.CrossBoundaryCopilotDataMovementEnabled` → `properties.copilotPolicies.crossBoundaryCopilotDataMovementEnabled` (dto.go:81) | GET environment **with `$expand=properties/copilotPolicies`** | 217-221 |
| 15 | `display_name` | String | **R** | — | `LengthAtLeast(1)` | `EnviromentPropertiesDto.DisplayName` → `properties.displayName` (dto.go:48) | GET environment | 222-228 |
| 16 | `billing_policy_id` | String | O + C, `Default: "00000000-0000-0000-0000-000000000000"` | `UseStateForUnknown()` | `LengthAtLeast(1)`, `RegexMatches(helpers.GuidRegex)` | `BillingPolicyDto.Id` → `properties.billingPolicy.id` (dto.go:56,116-118) | GET environment **with `$expand=properties/billingPolicy`**; also written via licensing API | 229-241 |
| 17 | `enterprise_policies` | SetNested | C | `setplanmodifier.UseStateForUnknown()` | — | `EnvironmentEnterprisePoliciesDto` → `properties.enterprisePolicies` (dto.go:61,94-98) | GET environment | 242-251 |
| 18 | `dataverse` | SingleNested | O + C | `modifiers.RequireReplaceObjectToEmptyModifier()` (internal/modifiers/require_replace_object_to_empty_modifier.go:12) | — | `LinkedEnvironmentMetadataDto` → `properties.linkedEnvironmentMetadata` (dto.go:52,148-159) | GET environment + Dataverse WebAPI (currency) | 252-375 |

Note: `SourceModel` field order/tags: internal/services/environment/models.go:38-58.

### A.2 `enterprise_policies` nested attributes (`policyAttributeSchema`, resource_environment.go:67-86)

All five are `Computed` only, no plan modifiers, no validators. The set is built in `convertEnterprisePolicyModelFromDto` (models.go:379-430) — note it **overwrites** rather than appends, so only one policy survives if multiple are present (identity, then vnets, then customerManagedKeys — last one wins).

| TF attribute | TF type | R/O/C | JSON path | Line |
| --- | --- | --- | --- | --- |
| `enterprise_policies[*].type` | String | C | synthesized literal: `"Identity"` / `"NetworkInjection"` / `"Encryption"` (models.go:395,413,431 area) | 68-71 |
| `enterprise_policies[*].id` | String | C | `properties.enterprisePolicies.{identity\|vnets\|customerManagedKeys}.id` (dto.go:103) | 72-75 |
| `enterprise_policies[*].location` | String | C | `...location` (dto.go:102) | 76-79 |
| `enterprise_policies[*].system_id` | String | C | `...systemId` (dto.go:104) | 80-83 |
| `enterprise_policies[*].status` | String | C | `...linkStatus` (dto.go:105) | 84-87 |

`EnterprisePolicyDto.PolicyId` (`policyId`, dto.go:101) is parsed but **not surfaced** in the schema.

### A.3 `dataverse` nested attributes (resource_environment.go:258-374)

| # | TF attribute | TF type | R/O/C | Plan modifiers | Validators | Go DTO field → JSON path | Populated by | Line |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `dataverse.unique_name` | String | C | `modifiers.UseStateForUnknownKeepNonNullStateModifier()` | — | `LinkedEnvironmentMetadataDto.UniqueName` → `properties.linkedEnvironmentMetadata.uniqueName` (dto.go:158) | GET environment | 260-266 |
| 2 | `dataverse.administration_mode_enabled` | Bool | O + C, `Default: false` | `boolplanmodifier.UseStateForUnknown()` | — | **Not** in `linkedEnvironmentMetadata`. Derived from `StatesEnvironmentDto.Runtime.Id` → `properties.states.runtime.id == "AdminMode"` (dto.go:53,177-189). Written as `properties.states.runtime.id` = `"AdminMode"`/`"Enabled"` | GET environment | 267-275 |
| 3 | `dataverse.background_operation_enabled` | Bool | O + C, `Default: true` | `boolplanmodifier.UseStateForUnknown()` | — | `LinkedEnvironmentMetadataDto.BackgroundOperationsState` → `properties.linkedEnvironmentMetadata.backgroundOperationsState` (`"Enabled"`/`"Disabled"`) (dto.go:149) | GET environment | 276-284 |
| 4 | `dataverse.currency_code` | String | **R** | `modifiers.RequireReplaceStringFromNonEmptyPlanModifier()` | — | Write: `createCurrencyDto.Code` → `properties.linkedEnvironmentMetadata.currency.code` (dto.go:228,233-235). **Read: NOT returned by BAPI** — read from Dataverse WebAPI `transactioncurrencies.isocurrencycode` (dto.go:274) | Dataverse WebAPI (`GetDefaultCurrencyForEnvironment`) | 285-291 |
| 5 | `dataverse.url` | String | C | `UseStateForUnknownKeepNonNullStateModifier()`, `modifiers.SetStringAttributeUnknownOnlyIfSecondAttributeChange(path.Root("dataverse").AtName("domain"))` | — | `LinkedEnvironmentMetadataDto.InstanceURL` → `properties.linkedEnvironmentMetadata.instanceUrl` (dto.go:151) | GET environment | 292-299 |
| 6 | `dataverse.domain` | String | O + C | `UseStateForUnknownKeepNonNullStateModifier()` | `LengthAtLeast(1)`, `RegexMatches(helpers.DomainNameRegex)` | `LinkedEnvironmentMetadataDto.DomainName` → `properties.linkedEnvironmentMetadata.domainName` (dto.go:150) | GET environment; validated with `validateEnvironmentDetails` | 300-311 |
| 7 | `dataverse.organization_id` | String | C | `UseStateForUnknownKeepNonNullStateModifier()` | — | `LinkedEnvironmentMetadataDto.ResourceId` → `properties.linkedEnvironmentMetadata.resourceId` (dto.go:154) | GET environment | 312-318 |
| 8 | `dataverse.security_group_id` | String | O + C | — | `RegexMatches(helpers.GuidRegex)`, `validators.MakeFieldRequiredWhenOtherFieldDoesNotHaveValue(environment_type, /^(Sandbox\|Production\|Trial\|Default)$/)` | `LinkedEnvironmentMetadataDto.SecurityGroupId` → `properties.linkedEnvironmentMetadata.securityGroupId` (dto.go:153) | GET environment; empty string coerced to `ZERO_UUID` (api_environment.go:266-268); forced `null` for `Developer` SKU (models.go:262-265) | 319-327 |
| 9 | `dataverse.language_code` | Int64 | **R** | `modifiers.RequireReplaceIntAttributePlanModifier()` | — | `LinkedEnvironmentMetadataDto.BaseLanguage` → `properties.linkedEnvironmentMetadata.baseLanguage` (dto.go:152) | GET environment; validated against `environmentLanguages` | 328-334 |
| 10 | `dataverse.version` | String | C | `UseStateForUnknownKeepNonNullStateModifier()` | — | `LinkedEnvironmentMetadataDto.Version` → `properties.linkedEnvironmentMetadata.version` (dto.go:155) | GET environment | 335-341 |
| 11 | `dataverse.templates` | List(String) | O + C | — | — | Read: `LinkedEnvironmentMetadataDto.Templates` with JSON tag **`template`** (singular!) (dto.go:156). Write: `createLinkEnvironmentMetadataDto.Templates` with JSON tag **`templates`** (plural) (dto.go:230) | GET environment; **BAPI does not echo it back after create**, so the create payload is stitched back into state (api_environment.go:531-534; resource_environment.go:486-488) | 342-347 |
| 12 | `dataverse.template_metadata` | String (raw JSON) | **O only (not Computed)** | — | — | `createTemplateMetadataDto` → `properties.linkedEnvironmentMetadata.templateMetadata` with inner key `PostProvisioningPackages` (dto.go:157,237-244). Marshalled/unmarshalled as a JSON string (models.go:176-186, 302-315) | GET environment or echoed from plan | 348-351 |
| 13 | `dataverse.linked_app_type` | String | C | `UseStateForUnknownKeepNonNullStateModifier()` | — | `LinkedAppMetadataDto.Type` → `properties.linkedAppMetadata.type` (dto.go:173) | GET environment | 352-358 |
| 14 | `dataverse.linked_app_id` | String | C | `UseStateForUnknownKeepNonNullStateModifier()` | — | `LinkedAppMetadataDto.Id` → `properties.linkedAppMetadata.id` (dto.go:172) | GET environment | 359-365 |
| 15 | `dataverse.linked_app_url` | String | C | `UseStateForUnknownKeepNonNullStateModifier()` | — | `LinkedAppMetadataDto.Url` → `properties.linkedAppMetadata.url` (dto.go:174) | GET environment | 366-372 |

Attribute type map used when constructing the object value: models.go:222-238 (`attrTypesDataverseObject`) — must stay in sync with the schema.

### A.4 Resource-level config validators

```go
// internal/services/environment/resource_environment.go:379-386
func (d *Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        resourcevalidator.RequiredTogether(
            path.Root("dataverse").AtName("administration_mode_enabled").Expression(),
            path.Root("dataverse").AtName("background_operation_enabled").Expression(),
        ),
    }
}
```

### A.5 Imperative (non-schema) validation

| Validation | Where | Behavior |
| --- | --- | --- |
| Location + azure region exist | `LocationValidator` (api_environment.go:79-97), called at resource_environment.go:429 | GET `/locations`, then match `name` and `properties.azureRegions` |
| Language code valid for location | `languageCodeValidator` (api_environment.go:183-227), called at resource_environment.go:443 | GET `/locations/{loc}/environmentLanguages` |
| Currency code valid for location | `currencyCodeValidator` (api_environment.go:119-159), called at resource_environment.go:449 | GET `/locations/{loc}/environmentCurrencies` |
| Domain uniqueness on create | `ValidateCreateEnvironmentDetails` (api_environment.go:721-741), called from `createEnvironmentWithRetry` (api_environment.go:461-466) | POST `/validateEnvironmentDetails` |
| Domain uniqueness on update | `ValidateUpdateEnvironmentDetails` (api_environment.go:743-757), called from `updateEnvironmentWithRetry` (api_environment.go:611-616) | POST `/validateEnvironmentDetails` |
| Generative AI feature rules | `aiGenerativeFeaturesValidaor` (resource_environment.go:888-899), called in Create (435) and Update (582) | Non-public cloud → error; `unitedstates` + `allow_moving_data_across_regions` → error; non-US + any AI flag without `allow_moving_data_across_regions` → error |
| Azure region actually honored | resource_environment.go:493-496 | After create, if requested `azure_region` != returned, `AddAttributeError` |

---

## B. `powerplatform_environments` data source — full attribute inventory

Schema at internal/services/environment/datasource_environments.go:48-232. Model: `ListDataSourceModel` (models.go:33-36) reusing the **same** `SourceModel` per element, so the element shape is identical to the resource model.

| # | TF attribute | TF type | R/O/C | Plan modifiers | Validators | JSON path | Line |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `timeouts` | Object (Read only) | O | — | — | n/a | 78-80 |
| 2 | `environments` | ListNested | C | — | — | `value[]` from `environmentArrayDto` (dto.go:195-197) | 81-231 |
| 3 | `environments[*].timeouts` | Object (all false) | O | — | — | n/a — present only because `SourceModel` carries `Timeouts` (models.go:39); populated with `timeouts.Value{}` (datasource_environments.go:290) | 86-91 |
| 4 | `environments[*].id` | String | C | — | — | `name` | 92-95 |
| 5 | `environments[*].location` | String | C | — | — | `location` | 96-99 |
| 6 | `environments[*].azure_region` | String | C | — | — | `properties.azureRegion` | 100-103 |
| 7 | `environments[*].environment_type` | String | C | — | — | `properties.environmentSku` | 104-107 |
| 8 | `environments[*].display_name` | String | C | — | — | `properties.displayName` | 108-111 |
| 9 | `environments[*].description` | String | C | — | — | `properties.description` | 112-115 |
| 10 | `environments[*].cadence` | String | C | — | `OneOf("Frequent","Moderate")` (unusual on a data source) | `properties.updateCadence.id` | 116-122 |
| 11 | `environments[*].release_cycle` | String | C | — | — | derived from `properties.cluster.category` | 123-126 |
| 12 | `environments[*].billing_policy_id` | String | C | — | — | `properties.billingPolicy.id` | 127-130 |
| 13 | `environments[*].environment_group_id` | String | C | — | — | `properties.parentEnvironmentGroup.id` | 131-134 |
| 14 | `environments[*].enterprise_policies` | SetNested | C | — | — | `properties.enterprisePolicies` | 135-141 |
| 15 | `environments[*].enterprise_policies[*].type` | String | C | — | — | synthesized | 53-56 |
| 16 | `environments[*].enterprise_policies[*].id` | String | C | — | — | `...id` | 57-60 |
| 17 | `environments[*].enterprise_policies[*].location` | String | C | — | — | `...location` | 61-64 |
| 18 | `environments[*].enterprise_policies[*].system_id` | String | C | — | — | `...systemId` | 65-68 |
| 19 | `environments[*].enterprise_policies[*].status` | String | C | — | — | `...linkStatus` | 69-72 |
| 20 | `environments[*].owner_id` | String | C | — | — | `properties.usedBy.id` | 142-145 |
| 21 | `environments[*].allow_bing_search` | Bool | C | — | — | `properties.bingChatEnabled` | 146-149 |
| 22 | `environments[*].allow_microsoft_365_services` | Bool | C | — | — | `properties.m365Enabled` | 150-153 |
| 23 | `environments[*].allow_moving_data_across_regions` | Bool | C | — | — | `properties.copilotPolicies.crossGeoCopilotDataMovementEnabled` | 154-157 |
| 24 | `environments[*].allow_flex_routing` | Bool | C | — | — | `properties.copilotPolicies.crossBoundaryCopilotDataMovementEnabled` | 158-161 |
| 25 | `environments[*].dataverse` | SingleNested | C | — | — | `properties.linkedEnvironmentMetadata` | 162-229 |
| 26 | `environments[*].dataverse.unique_name` | String | C | — | — | `...uniqueName` | 166-169 |
| 27 | `environments[*].dataverse.administration_mode_enabled` | Bool | C | — | — | `properties.states.runtime.id == "AdminMode"` | 170-173 |
| 28 | `environments[*].dataverse.background_operation_enabled` | Bool | C | — | — | `...backgroundOperationsState` | 174-177 |
| 29 | `environments[*].dataverse.url` | String | C | — | — | `...instanceUrl` | 178-181 |
| 30 | `environments[*].dataverse.domain` | String | C | — | — | `...domainName` | 182-185 |
| 31 | `environments[*].dataverse.organization_id` | String | C | — | — | `...resourceId` | 186-189 |
| 32 | `environments[*].dataverse.security_group_id` | String | C | — | — | `...securityGroupId` | 190-193 |
| 33 | `environments[*].dataverse.language_code` | Int64 | C | — | — | `...baseLanguage` | 194-197 |
| 34 | `environments[*].dataverse.version` | String | C | — | — | `...version` | 198-201 |
| 35 | `environments[*].dataverse.linked_app_type` | String | C | — | — | `properties.linkedAppMetadata.type` | 202-205 |
| 36 | `environments[*].dataverse.linked_app_id` | String | C | — | — | `properties.linkedAppMetadata.id` | 206-209 |
| 37 | `environments[*].dataverse.linked_app_url` | String | C | — | — | `properties.linkedAppMetadata.url` | 210-213 |
| 38 | `environments[*].dataverse.currency_code` | String | C | — | — | Dataverse WebAPI `transactioncurrencies.isocurrencycode` | 214-217 |
| 39 | `environments[*].dataverse.templates` | List(String) | C | — | — | `...template` (singular JSON tag) | 218-222 |
| 40 | `environments[*].dataverse.template_metadata` | String | C | — | — | `...templateMetadata` | 223-226 |

Data source has **no filter parameters** — it always lists everything.

### B.1 Data source Read behavior (datasource_environments.go:243-299)

1. `GetEnvironments(ctx)` — one BAPI list call.
2. **For each environment**, `GetDefaultCurrencyForEnvironment(ctx, env.Name)` — which itself does `GetEnvironment` (to resolve the host) + 2 Dataverse WebAPI calls. So reading N environments costs **1 + 3N** HTTP calls. This is a known perf hot spot and shows up in the fixtures (`tests/datasource/Validate_Read/get_environment_*.json`).
3. Errors other than `ErrEnvironmentUrlNotFound` are **fatal** in the data source (datasource_environments.go:278-284), unlike the resource which downgrades them to a warning (resource_environment.go:531-535).
4. `convertSourceModelFromEnvironmentDto(env, &currencyCode, nil, nil, nil, timeouts.Value{}, config)`.

---

## C. Exact BAPI call inventory

All calls use `client.Api.GetConfig().Urls.BapiUrl` (public: `api.bap.microsoft.com`, constants.go:29) unless noted. Scope is inferred from the URL by `tryGetScopeFromURL` (internal/api/client.go:253-270) → `PowerAppsScope` (`https://service.powerapps.com/.default`) for BAPI hosts.

API versions (internal/constants/constants.go:194-203):

```go
BAP_API_VERSION      = "2023-06-01"
BAP_2021_API_VERSION = "2021-04-01"
BAP_2022_API_VERSION = "2022-05-01"   // defined but not used by environment service
```

### C.1 Call table

| # | Op | Method | URL template | api-version | Query params | Request DTO | Response DTO | Async? | Line |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | List locations | GET | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/locations` | `2023-06-01` | — | — | `LocationArrayDto` (dto.go:295) | No | api_environment.go:61-77 |
| 2 | List currencies for location | GET | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/locations/{location}/environmentCurrencies` | `2023-06-01` | — | — | `currencyCodeValidatorArrayDto` (api_environment.go:112) | No | api_environment.go:119-159 |
| 3 | List languages for location | GET | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/locations/{location}/environmentLanguages` | `2023-06-01` | — | — | `languageCodeValidatorArrayDto` (api_environment.go:166) | No | api_environment.go:183-227 |
| 4 | **Get environment** | GET | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environmentId}` | `2023-06-01` | `$expand=permissions,properties.capacity,properties/billingPolicy,properties/copilotPolicies` | — | `EnvironmentDto` | No | api_environment.go:246-275 |
| 5 | **List environments** | GET | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments` | `2023-06-01` | `$expand=properties/billingPolicy,properties/copilotPolicies` | — | `environmentArrayDto` | No | api_environment.go:681-699 |
| 6 | **Create environment** | POST | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/environments` | `2023-06-01` | — | `environmentCreateDto` (dto.go:199) | `lifecycleCreatedDto` on 201; lifecycle poll on 202 | **Yes** | api_environment.go:460-539 |
| 7 | **Update environment** | PATCH | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environmentId}` | `2021-04-01` | `$expand=permissions,properties.capacity,properties/billingPolicy` | `EnvironmentDto` | — (poll + re-GET) | **Yes** | api_environment.go:610-679 |
| 8 | **Update AI features** | PATCH | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environmentId}` | `2021-04-01` | — | `GenerativeAiFeaturesDto` (dto.go:69) | — | **Yes** | api_environment.go:545-591 |
| 9 | **Modify SKU (env type)** | POST | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environmentId}/modifySku` | `2021-04-01` | — | `modifySkuDto` (dto.go:204) | — | **Yes** | api_environment.go:419-454 |
| 10 | **Provision Dataverse** | POST | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/environments/{environmentId}/provisionInstance` | `2021-04-01` | — | `createLinkEnvironmentMetadataDto` (dto.go:225) | `EnvironmentDto` (polled) | **Yes, custom poll** | api_environment.go:350-413 |
| 11 | **Delete environment** | DELETE | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environmentId}` | `2023-06-01` | — | `enironmentDeleteDto` `{code:"7", message:"Deleted using Power Platform Terraform Provider"}` (dto.go:246) | — | **Yes** | api_environment.go:281-348 |
| 12 | Validate create details | POST | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails` | `2021-04-01` | — | `validateCreateEnvironmentDetailsDto` `{domainName, environmentLocation}` (dto.go:285) | — | No | api_environment.go:721-741 |
| 13 | Validate update details | POST | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails` | `2021-04-01` | — | `validateUpdateEnvironmentDetailsDto` `{domainName, environmentName}` (dto.go:290) | — | No | api_environment.go:743-757 |
| 14 | Get org base currency | GET | `https://{envHost}/api/data/v9.2/organizations` (Dataverse WebAPI) | n/a | — | — | `organizationSettingsArrayDto` (dto.go:260) | No | api_environment.go:701-720 → solution/api_solution.go:356-378 |
| 15 | Get currency by id | GET | `https://{envHost}/api/data/v9.2/transactioncurrencies` (Dataverse WebAPI) | n/a | `$filter=transactioncurrencyid eq {baseCurrencyId}` | — | `transactionCurrencyArrayDto` (dto.go:281) | No | api_environment.go:706-716 |
| 16 | Get tenant (for `owner_id`) | GET | `https://{bapi}/providers/Microsoft.BusinessAppPlatform/tenant` | `2021-04-01` | — | — | `TenantDto` | No | tenant/api_tenant.go:26-40, called from models.go:145 |
| 17 | Add env to billing policy | POST | `https://{powerplatform}/licensing/billingPolicies/{billingId}/environments/add` | `2022-03-01-preview` | — | `{environmentIds:[...]}` | — | No | licensing/api_licensing.go:167-192, called from resource_environment.go:877-886 |
| 18 | Remove env from billing policy | POST | `https://{powerplatform}/licensing/billingPolicies/{billingId}/environments/remove` | `2022-03-01-preview` | — | `{environmentIds:[...]}` | — | No | licensing/api_licensing.go:194-219, called from resource_environment.go:866-875 |

**Explicitly NOT implemented today** (no code found):

* `validateDelete` — the provider does not pre-validate deletion.
* Soft-delete / recover — no `recover` endpoint call anywhere in the repo.
* Reset / copy environment — not implemented.
* "Add currency" / "add language" as separate operations — currency and language are only set at Dataverse provisioning time (`currency.code`, `baseLanguage`). There are no `addCurrency`/`addLanguage` calls.
* `Microsoft.BusinessAppPlatform/scopes/admin/...` admin-scope endpoints beyond #4, #5, #7, #8, #9, #11 above.

### C.2 Async / lifecycle polling

Generic helper: `(*api.Client).DoWaitForLifecycleOperationStatus` at internal/api/lifecycle.go:53-103.

```go
// internal/api/lifecycle.go:53-70
locationHeader := response.GetHeader(constants.HEADER_LOCATION)      // "Location"
if locationHeader == "" {
    locationHeader = response.GetHeader(constants.HEADER_OPERATION_LOCATION) // "Operation-Location"
}
if locationHeader == "" { return nil, nil }   // treated as synchronous success
waitFor := retryAfter(ctx, response.HttpResponse)                    // reads "Retry-After"
```

Polling loop semantics:

* GET the `Location` URL, accepting `200`, `409`, `404`.
* `404` → synthesizes `State.Id = "Succeeded"` and returns (lifecycle.go:80-86) — "resource already gone" is success.
* `409` → sleep `Retry-After` and retry (lifecycle.go:88-94).
* Terminal states: `"Succeeded"` or `"Failed"` (lifecycle.go:96-98). Anything else → sleep and loop.

`LifecycleDto` shape (lifecycle.go:15-51):

```go
type LifecycleDto struct {
    Id                 string
    Links              LifecycleLinksDto  // .Self.Path, .Environment.Path
    State              LifecycleStateDto  // .Id: "Succeeded"|"Failed"|...
    Type               LifecycleStateDto
    CreatedDateTime    string
    LastActionDateTime string
    RequestedBy        LifecycleRequestedByDto
    Stages             []LifecycleStageDto
}
```

Per-operation async handling:

| Op | Accepted status codes | Polling | Post-poll behavior |
| --- | --- | --- | --- |
| Create (api_environment.go:475) | `202, 201, 500, 409` | `202` → `DoWaitForLifecycleOperationStatus`; environment id parsed from `lifecycleResponse.Links.Environment.Path` last segment (api_environment.go:503-513). `201` → read `lifecycleCreatedDto.name` + `properties.provisioningState == "Succeeded"` (515-524) | GET environment; re-stitch `templates`/`templateMetadata` from the create payload (531-534). `500` → `ErrEnvironmentCreation` provider error (496-498). `409` → `handleHttpConflict` + retry (487-495) |
| Update (api_environment.go:628) | `202, 409` | `DoWaitForLifecycleOperationStatus`; `Failed` → sleep + retry up to `MAX_RETRY_COUNT` | Then busy-loop GET environment until `properties.states.management.id == "Ready"`; `"Running"` → keep polling; anything else → error (api_environment.go:659-678) |
| Update AI features (api_environment.go:555) | `202, 409, 204` | `204` → no-op return (557-560). `409` → `handleHttpConflict` + retry. Else `DoWaitForLifecycleOperationStatus` | `Failed` → retry | 
| Modify SKU (api_environment.go:434) | `202, 200, 409` | `DoWaitForLifecycleOperationStatus` | `Failed` → sleep + retry |
| Delete (api_environment.go:292-296) | `204, 202, 409, 404` | `404` → treat as deleted (299-302). `409` with empty body → sleep + retry (309-316); `409` with body → `handleHttpConflict`. Else `DoWaitForLifecycleOperationStatus` | `Failed` → sleep + retry |
| Provision Dataverse (api_environment.go:359) | `202` only | **Custom loop** (api_environment.go:361-412): reads `Location` + `Retry-After` headers itself, GETs the location URL expecting `200/202/409`, unmarshals into `EnvironmentDto`, and waits for `properties.provisioningState == "Succeeded"`; anything other than `LinkedDatabaseProvisioning`/`Succeeded` → error | Returns the polled `EnvironmentDto` |

Conflict handling helper (api_environment.go:593-604):

```go
func (client *Client) handleHttpConflict(ctx context.Context, apiResponse *api.Response) error {
    body := string(apiResponse.BodyAsBytes)
    if body == "" { return errors.New("environment failed with HTTP 409. No body in response") }
    if !strings.Contains(body, "OperationNotStartable") {
        return errors.New("environment failed with HTTP 409. Body: " + body)
    }
    return client.Api.SleepWithContext(ctx, api.DefaultRetryAfter())
}
```

Retry budget: `constants.MAX_RETRY_COUNT = 10` (constants.go:141). `api.DefaultRetryAfter()` = random 10–20 s (internal/api/client.go:216-218). `SleepWithContext` is a no-op in test contexts (client.go:221-235).

---

## D. Update semantics — what changes in place today

`Update` is at internal/services/environment/resource_environment.go:571-682. Sequence:

```text
 1. aiGenerativeFeaturesValidaor(plan)                        (582)
 2. build envProp := EnviromentPropertiesDto{
        DisplayName, EnvironmentSku, BingChatEnabled, M365Enabled }  (588-593)
 3. updateEnvironmentType(plan, state)   -> POST modifySku IF changed (599)
 4. updateDescription(plan, &dto)                             (604)
 5. updateCadence(plan, &dto)                                 (605)
 6. updateEnvironmentGroupId(plan, &dto)                      (607)
 7. updateBillingPolicyId(plan, &dto)                         (608)
 8. updateDataverse(plan, state, &dto)   -> in-place OR POST provisionInstance (610)
 9. removeBillingPolicy(state)           -> POST licensing .../environments/remove (616)
10. addBillingPolicy(plan)               -> POST licensing .../environments/add    (621)
11. UpdateEnvironment(id, dto)           -> PATCH admin/environments/{id}          (630)
12. IF display_name changed: PATCH AGAIN (BAPI propagation bug workaround)         (636-642)
13. updateGenerativeAiFeatures(plan)     -> PATCH with GenerativeAiFeaturesDto     (645)
14. GetEnvironment(id) and rebuild state                                            (651)
```

### D.1 Field-by-field update mapping

| TF attribute | Updatable in place? | Mechanism | Call | Source |
| --- | --- | --- | --- | --- |
| `display_name` | Yes | `envProp.DisplayName` | PATCH `admin/environments/{id}` — **sent twice** because of a BAPI propagation bug | resource_environment.go:589, 636-642 |
| `environment_type` | Yes | dedicated call **before** the PATCH; also mirrored into `envProp.EnvironmentSku` | POST `.../modifySku` | resource_environment.go:590, 599, 802-810; api_environment.go:419-454 |
| `description` | Yes (only when non-null and non-empty) | `properties.description` | PATCH | resource_environment.go:604, 812-816 |
| `cadence` | Yes (only when non-null and non-empty) | `properties.updateCadence.id` | PATCH | resource_environment.go:605, 818-824 |
| `environment_group_id` | Yes (only when not null) | `properties.parentEnvironmentGroup.id` | PATCH | resource_environment.go:607, 836-842 |
| `billing_policy_id` | Yes | Three-way: PATCH body `properties.billingPolicy.id` **plus** licensing remove-then-add | PATCH + POST licensing add/remove | resource_environment.go:608, 616, 621, 844-850, 866-886 |
| `allow_bing_search` | Yes | `envProp.BingChatEnabled` in PATCH **and** `GenerativeAiFeaturesDto.Properties.BingChatEnabled` in the second PATCH | PATCH ×2 | resource_environment.go:591, 645, 684-706 |
| `allow_microsoft_365_services` | Yes | `envProp.M365Enabled` + `GenerativeAiFeaturesDto...M365Enabled` | PATCH ×2 | resource_environment.go:592, 645, 684-706 |
| `allow_moving_data_across_regions` | Yes | `GenerativeAiFeaturesDto.Properties.CopilotPolicies.CrossGeoCopilotDataMovementEnabled` only | PATCH (AI features) | resource_environment.go:691-698 |
| `allow_flex_routing` | Yes | `...CopilotPolicies.CrossBoundaryCopilotDataMovementEnabled` only | PATCH (AI features) | resource_environment.go:692-698 |
| `dataverse.domain` | Yes | `properties.linkedEnvironmentMetadata.domainName`, guarded by `ValidateUpdateEnvironmentDetails` | PATCH | resource_environment.go:721-773 (748-750); api_environment.go:611-616 |
| `dataverse.security_group_id` | Yes | `properties.linkedEnvironmentMetadata.securityGroupId` | PATCH | resource_environment.go:729 |
| `dataverse.administration_mode_enabled` | Yes | `properties.states.runtime.id` = `"AdminMode"` / `"Enabled"` | PATCH | resource_environment.go:731-745 |
| `dataverse.background_operation_enabled` | Yes | `properties.linkedEnvironmentMetadata.backgroundOperationsState` = `"Enabled"` / `"Disabled"` | PATCH | resource_environment.go:747-753 |
| `dataverse.linked_app_id/type/url` | Yes (mirrored from plan; all three are `Computed`, so this only fires on import/drift) | `properties.linkedAppMetadata` set or nulled | PATCH | resource_environment.go:758-768 |
| `dataverse` null → set (add Dataverse to an existing env) | Yes | `addDataverse` | POST `.../provisionInstance` | resource_environment.go:708-719, 852-864 |
| `dataverse` set → null (remove Dataverse) | **No — forces replace** | `modifiers.RequireReplaceObjectToEmptyModifier()` | n/a | resource_environment.go:255-257 |
| `owner_id` | **No — RequiresReplace** | — | n/a | resource_environment.go:192-195 |
| `location` | **No — RequiresReplace** | — | n/a | resource_environment.go:167-169 |
| `azure_region` | **No — RequiresReplace** | — | n/a | resource_environment.go:175-178 |
| `release_cycle` | **No — RequiresReplace** | — | n/a | resource_environment.go:158-161 |
| `dataverse.currency_code` | **No — RequireReplaceStringFromNonEmptyPlanModifier** (replace when changed from a non-empty value) | — | n/a | resource_environment.go:288-290 |
| `dataverse.language_code` | **No — RequireReplaceIntAttributePlanModifier** | — | n/a | resource_environment.go:331-333 |
| `dataverse.templates` / `template_metadata` | Not updatable (create-only inputs; drift is masked by re-stitching plan values into state) | — | n/a | resource_environment.go:485-488, 546-558, 663-674 |
| `id`, `enterprise_policies`, `dataverse.url/organization_id/version/unique_name` | Read-only | — | n/a | — |

### D.2 Update quirks worth carrying into the migration

* **Double PATCH for rename**: resource_environment.go:635-643 — `// This is a temporary fix for the issue in BAPI where the display name is not propagated correctly on environment update`.
* **Older api-version deliberately pinned on update**: api_environment.go:625-627 — `// Due to a bug in BAPI that triggers managed environment on update of a description field, we need to use the older API version`, so update uses `2021-04-01`, not `2023-06-01`.
* **`environment_group_id` has no default**: resource_environment.go:127-128 — `// TODO: because of the bug on the backend default value would trigger managed environment creation`.
* **`updateExistingDataverse` always sets `DomainName` and `SecurityGroupId`** (resource_environment.go:725-729), even when unchanged — so a PATCH always carries them.
* **Post-update readiness spin**: `updateEnvironmentWithRetry` re-GETs the environment until `properties.states.management.id == "Ready"` (api_environment.go:659-678). This is an extra readiness gate beyond lifecycle success.

---

## E. Provider conventions to preserve

### E.1 Resource / data source struct + request context

```go
// internal/services/environment/models.go:22-31
type EnvironmentsDataSource struct {
    helpers.TypeInfo
    EnvironmentClient Client
}

type Resource struct {
    helpers.TypeInfo
    EnvironmentClient Client
    LicensingClient   licensing.Client
}
```

Every interface method opens a request context:

```go
// internal/services/environment/resource_environment.go:412-413 (identical pattern at 51-57, 63-65, 388-390, 508-510, 571-573, 775-777, 795-797)
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
```

`EnterRequestContext` is at internal/helpers/contexts.go:63. `Metadata` also sets `r.ProviderTypeName = req.ProviderTypeName` **before** entering the context (resource_environment.go:53) and returns `r.FullTypeName()`.

Factories: `NewEnvironmentResource()` (resource_environment.go:42-49) and `NewEnvironmentsDataSource()` (datasource_environments.go:28-34), both seeding `helpers.TypeInfo{TypeName: "..."}`.

`Configure` casts `req.ProviderData.(*api.ProviderClient)` and tolerates `nil` ProviderData (resource_environment.go:388-409; datasource_environments.go:234-254).

### E.2 Client factory

```go
// internal/services/environment/api_environment.go:24-35
func NewEnvironmentClient(apiClient *api.Client) Client {
    return Client{
        tenantClient:   tenant.NewTenantClient(apiClient),
        solutionClient: solution.NewSolutionClient(apiClient),
        Api:            apiClient,
    }
}

type Client struct {
    tenantClient   tenant.Client
    solutionClient solution.Client
    Api            *api.Client
}
```

Convention: `New<Service>Client(apiClient *api.Client) Client`, exported `Api *api.Client` field, unexported dependency clients.

### E.3 URL construction + Execute

```go
// canonical shape, e.g. internal/services/environment/api_environment.go:246-258
apiUrl := &url.URL{
    Scheme: constants.HTTPS,
    Host:   client.Api.GetConfig().Urls.BapiUrl,
    Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
}
values := url.Values{}
values.Add("$expand", "...")
values.Add(constants.API_VERSION_PARAM, constants.BAP_API_VERSION)
apiUrl.RawQuery = values.Encode()

resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
```

`Execute` signature: `Execute(ctx, scopes []string, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any) (*Response, error)` (internal/api/client.go:113-115). Passing `nil` scopes triggers URL-based scope inference. There is also `ExecuteWithoutRetry` (client.go:123-125) for non-idempotent mutation starts. There is a `helpers.BuildApiUrl(host, path, query)` helper (internal/helpers/uri.go:36+) but the environment service does not use it.

### E.4 Conversion function conventions (`models.go`)

The repo-wide guidance names them `convertDtoToModel` / `convertModelToDto`; the environment service uses more descriptive variants:

| Function | Direction | Line |
| --- | --- | --- |
| `convertCreateEnvironmentDtoFromSourceModel(ctx, *SourceModel, *Resource) (*environmentCreateDto, error)` | model → DTO | models.go:93-168 |
| `convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, types.Object) (*createLinkEnvironmentMetadataDto, error)` | model → DTO | models.go:170-205 |
| `convertSourceModelFromEnvironmentDto(EnvironmentDto, currencyCode, ownerId *string, templateMetadata, templates, timeouts.Value, config.ProviderConfig) (*SourceModel, error)` | DTO → model | models.go:207-338 |
| `convertEnvironmentGroupFromDto` | DTO → model (field) | models.go:340-346 |
| `convertBillingPolicyModelFromDto` | DTO → model (field) | models.go:348-354 |
| `convertReleaseCycleModelFromDto` | DTO → model (field) | models.go:356-363 |
| `convertOwnerIdFromDto` | DTO → model (field) | models.go:365-371 |
| `convertCopilotPoliciesFromDto` | DTO → model (field) | models.go:373-377 |
| `convertEnterprisePolicyModelFromDto` | DTO → model (field) | models.go:379-430 |
| `isDataverseEnvironmentEmpty(ctx, *SourceModel) bool` | predicate | models.go:86-91 |

Note the DTO→model function takes extra out-of-band inputs (`currencyCode`, `ownerId`, `templateMetadata`, `templates`) precisely because BAPI does not round-trip those values.

### E.5 Error handling

Sentinels (internal/customerrors/provider_error.go:16-23):

```go
ErrObjectNotFound            = ProviderError{ErrorCode: ErrorCode(constants.ERROR_OBJECT_NOT_FOUND)}
ErrEnvironmentUrlNotFound    = ProviderError{ErrorCode: ErrorCode(constants.ERROR_ENVIRONMENT_URL_NOT_FOUND)}
ErrEnvironmentsInEnvGroup    = ...
ErrPolicyAssignedToEnvGroup  = ...
ErrEnvironmentSettingsFailed = ...
ErrEnvironmentCreation       = ProviderError{ErrorCode: ErrorCode(constants.ERROR_ENVIRONMENT_CREATION)}
```

Usage in the environment service:

* `customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("environment '%s' not found", environmentId))` on HTTP 404 (api_environment.go:260-263).
* `WrapIntoProviderError(nil, ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, ...")` when `instanceUrl` is empty (api_environment.go:236-238).
* `WrapIntoProviderError(nil, ERROR_ENVIRONMENT_CREATION, string(body))` on HTTP 500 during create (api_environment.go:497).
* Resource `Read`/`Update` do `if errors.Is(err, customerrors.ErrObjectNotFound) { resp.State.RemoveResource(ctx); return }` (resource_environment.go:520-524, 652-656, 470-474).
* Diagnostics phrasing: `resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())`.
* Currency read failure is a **warning** in the resource (resource_environment.go:533) but an **error** in the data source (datasource_environments.go:279-283).

### E.6 Unit test structure + fixture conventions

Test package is `environment_test` (resource_environment_test.go:4). Naming: `TestUnitEnvironmentsResource_<Scenario>` / `TestAccEnvironmentsResource_<Scenario>`; data source: `TestUnitEnvironmentsDataSource_Validate_Read`, `TestAccEnvironmentsDataSource_Basic`.

Canonical unit test skeleton (resource_environment_test.go:1523-1597):

```go
func TestUnitEnvironmentsResource_Validate_Create(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    mocks.ActivateEnvironmentHttpMocks()

    httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
        func(req *http.Request) (*http.Response, error) {
            id := httpmock.MustGetSubmatch(req, 1)
            return httpmock.NewStringResponse(http.StatusOK,
                httpmock.File(fmt.Sprintf("tests/resource/Validate_Create/get_environment_%s.json", id)).String()), nil
        })

    httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
        func(req *http.Request) (*http.Response, error) {
            resp := httpmock.NewStringResponse(http.StatusAccepted, "")
            resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-.../?api-version=2023-06-01")
            return resp, nil
        })

    resource.Test(t, resource.TestCase{
        IsUnitTest:               true,
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{ { Config: `...`, Check: resource.ComposeTestCheckFunc(...) } },
    })
}
```

Shared mock bootstrap `mocks.ActivateEnvironmentHttpMocks()` (internal/mocks/mocks.go:46-113) registers:

* a `RegisterNoResponder` that **fails** on any unmocked call (mocks.go:47-49),
* Dataverse `transactioncurrencies` and `organizations` on `*.crm4.dynamics.com` (mocks.go:51-66),
* BAPI `environmentLanguages` / `environmentCurrencies` / `locations` (mocks.go:68-81),
* `validateEnvironmentDetails` → 200 (mocks.go:83-86),
* Dataverse `WhoAmI` (mocks.go:88-97),
* BAPI `tenant` (mocks.go:99-112).

Fixture layout: `internal/services/environment/tests/{resource|datasource}/<Scenario>/<method>_<object>[_<n>].json`. Scenario folder name == test name minus the `TestUnitEnvironmentsResource_` / `TestUnitEnvironmentsDataSource_` prefix. 149 files exist. Representative folders:

```text
tests/resource/Validate_Create/                       get_environment_00000000-...-000000000001.json, get_lifecycle.json, get_lifecycle_delete.json
tests/resource/Validate_Create_And_Update/            get_environment_0..3.json, get_environments_1.json, get_lifecycle_1.json, get_lifecycle_delete.json
tests/resource/Validate_Create_And_Force_Recreate/    get_environment_*.json, get_lifecycle_*.json, get_transactioncurrencies_*.json
tests/resource/Validate_Create_Environment_And_Dataverse/  get_environment_1..8.json, get_lifecycle_new_dataverse.json
tests/resource/Validate_Update_Environment_Type/      get_environment_0..10.json
tests/resource/Validate_Update_With_Billing_Policy/   get_environment_1..17.json, get_environments_1..3.json
tests/resource/Validate_Update_Security_Group_Id/     get_environment_0..1.json, get_environments_0..1.json
tests/resource/Validate_Update_Generative_Ai_Features/get_environment.json, get_lifecycle.json
tests/resource/Validate_Update_Microsoft_365_Services/get_environment_enabled.json, get_environment_disabled.json
tests/resource/Validate_Create_With_D365_Template/    get_environment_*.json, get_environments.json
tests/resource/Validate_Create_No_Dataverse/          ...
tests/resource/Validate_Create_Them_Try_Remove_Dataverse/ get_environment_1..8.json
tests/resource/Validate_Create_Dev_Env/               ...
tests/resource/Validate_Create_Early_Release_Cycle/   ...
tests/resource/Validate_Locations_And_Azure_Regions/  ...
tests/resource/Create_Environment_With_Env_Group/     get_environment_group.json
tests/resource/Create_Environment_And_Add_Env_Group/  post_environment_group_1..2.json, get_environment_group_1..2.json
tests/resource/Validate_Retry_LidecycleOperation/     get_lifecycle_delete_1..3.json   (note the typo in the folder name)
tests/resource/Validate_Retry_On_Running_LifecycleOperation/ post_environment_operation_in_progress.json
tests/resource/Validate_Delete_Retry_On_EmptyBody_Conflict/ ...
tests/resource/Validate_Delete_With_Lifecycle_404/    ...
tests/resource/Validate_Do_Not_Retry_On_NoCapacity/   post_environment.json
tests/resource/Validate_Create_No_Dataverse_Region_Not_Available/ ...
tests/resource/Validate_Attribute_Validators/         ...
tests/datasource/Validate_Read/                       get_environments.json, get_environment_00000000-...-000000000001.json, ...-000000000002.json
```

Full unit test list (59 tests total across both files) — the ones that will need re-mocking after migration:

```text
TestUnitEnvironmentsResource_Validate_Attribute_Validators                                      (:39)
TestUnitEnvironmentsResource_Validate_Do_Not_Retry_On_NoCapacity                                (:138)
TestUnitEnvironmentsResource_Validate_Retry_On_Running_LifecycleOperation                       (:175)
TestUnitEnvironmentsResource_Validate_Delete_Retry_On_EmptyBody_Conflict                        (:244)
TestUnitEnvironmentsResource_Validate_Retry_LifecycleOperation                                  (:311)
TestUnitEnvironmentsResource_Validate_Create_Error_Check_Environment_Group                      (:709)
TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_Security_Group   (:728)
TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_No_Dataverse     (:752)
TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_No_Developer_Env (:771)
TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_No_OwnerId       (:794)
TestUnitEnvironmentsResource_Validate_CreateDevelopmentEnvironment_Error_Check_Empty_OwnerId    (:816)
TestUnitEnvironmentsResource_Validate_CreateDeveloperEnvironment                                (:839)
TestUnitEnvironmentsResource_Validate_Create_Early_Release_Cycle                                (:937)
TestUnitEnvironmentsResource_Validate_Create_And_Force_Recreate                                 (:1171)
TestUnitEnvironmentsResource_Validate_Create_And_Update                                         (:1314)
TestUnitEnvironmentsResource_Validate_Update_Security_Group_Id                                  (:1418)
TestUnitEnvironmentsResource_Validate_Create                                                    (:1523)
TestUnitEnvironmentsResource_Validate_Create_With_Billing_Policy                                (:1598)
TestUnitEnvironmentsResource_Validate_Update_With_Billing_Policy                                (:1662)
TestUnitEnvironmentsResource_Validate_Create_With_D365_Template                                 (:1834)
TestUnitEnvironmentsResource_Validate_Taken_Domain_Name                                         (:1930)
TestUnitEnvironmentsResource_Validate_Domain_Format_Valid                                       (:1971)
TestUnitEnvironmentsResource_Validate_Domain_Format_Invalid_Characters                          (:1996)
TestUnitEnvironmentsResource_Validate_Create_No_Dataverse                                       (:2020)
TestUnitEnvironmentsResource_Validate_Update_Generative_Ai_Features                             (:2094)
TestUnitEnvironmentsResource_Validate_Microsoft_365_Services_Requires_Moving_Data_Across_Regions(:2231)
TestUnitEnvironmentsResource_Validate_Flex_Routing_Requires_Moving_Data_Across_Regions           (:2258)
TestUnitEnvironmentsResource_Validate_Create_Them_Try_Remove_Dataverse                          (:2353)
TestUnitEnvironmentsResource_Validate_Create_Environment_And_Dataverse                          (:2435)
TestUnitEnvironmentsResource_Validate_Locations_And_Azure_Regions                               (:2662)
TestUnitEnvironmentsResource_Create_Environment_With_Env_Group                                  (:2838)
TestUnitEnvironmentsResource_Create_Environment_And_Add_Env_Group                               (:2965)
TestUnitEnvironmentsResource_Validate_Update_Environment_Type                                   (:3381)
TestUnitEnvironmentsResource_Validate_Create_No_Dataverse_Region_Not_Available                  (:3558)
TestUnitEnvironmentsResource_Validate_Delete_With_Lifecycle_404                                 (:3621)
TestUnitEnvironmentsDataSource_Validate_Read                                (datasource_environments_test.go:71)
```

Acceptance tests (24): `TestAccEnvironmentsResource_Validate_Update_Name_Field` (:374), `..._CreateGenerativeAiFeatures_Non_US_Region_Update` (:413), `..._US_Region_Update` (:496), `..._US_Region_Expect_Fail` (:560), `..._Non_US_Region_Expect_Fail` (:582), `..._Non_US_Region_Microsoft_365_Services_Expect_Fail` (:604), `..._Non_US_Region_Flex_Routing_Expect_Fail` (:627), `..._CreateDeveloperEnvironment` (:650), `..._Create_Early_Release_Cycle` (:916), `..._Validate_Update` (:1000), `..._Validate_Domain_Uniqueness_On_Update` (:1070), `..._Validate_Create` (:1126), `..._Validate_Create_No_Dataverse` (:2285), `..._Validate_Create_Them_Try_Remove_Dataverse` (:2315), `..._Validate_Create_Environment_And_Dataverse` (:2566), `..._Validate_Locations_And_Azure_Regions` (:2621), `..._Validate_Enable_Admin_Mode` (:2738), `..._Create_Environment_With_Env_Group` (:2933), `..._Create_Environment_And_Add_Env_Group` (:3181), `..._Create_Environment_No_Dataverse_Add_Dataverse_And_Add_Env_Group` (:3272), `..._Create_Environment_No_Dataverse_Add_Env_Group` (:3357), `..._Validate_Update_Environment_Type` (:3487), plus `TestAccEnvironmentsDataSource_Basic` (datasource_environments_test.go:20).

### E.7 Changie changelog requirement

Config: .changie.yaml. Kinds: `breaking`, `added`, `changed`, `deprecated`, `removed`, `fixed`, `security`, `documentation` (.changie.yaml:9-32). Custom field `Issue` is an `int` with `minInt: 1` (.changie.yaml:44-48). Change line format renders as a link to `https://github.com/microsoft/terraform-provider-power-platform/issues/{Issue}` (.changie.yaml:8).

Exactly one entry per PR:

```bash
changie new --kind changed --body "<concise summary>" --custom Issue=<issue_number>
```

### E.8 Docs

`docs/resources/environment.md` and `docs/data-sources/environments.md` are generated from `MarkdownDescription` via `make userdocs`. Do not hand-edit. The generated resource doc groups attributes as Required (`display_name`, `environment_type`, `location`), Optional (11), Read-Only (`enterprise_policies`, `id`) — docs/resources/environment.md:59-134.

---

## F. Existing `api.powerplatform.com` usage in this repo

Host comes from `client.Api.GetConfig().Urls.PowerPlatformUrl` (internal/config/config.go:71), populated per cloud in internal/provider/provider.go:511-608 from `constants.*_POWERPLATFORM_API_DOMAIN`.

Scope is **never passed explicitly** — `Execute(ctx, nil, ...)` lets `tryGetScopeFromURL` (internal/api/client.go:253-270) match the host against `cloudConfig.PowerPlatformUrl` and return `cloudConfig.PowerPlatformScope`:

```go
// internal/api/client.go:258-259
case strings.LastIndex(url, cloudConfig.PowerPlatformUrl) != -1:
    return cloudConfig.PowerPlatformScope, nil
```

Public cloud values (internal/constants/constants.go:32-33):

```go
PUBLIC_POWERPLATFORM_API_DOMAIN = "api.powerplatform.com"
PUBLIC_POWERPLATFORM_API_SCOPE  = "https://api.powerplatform.com/.default"
```

Because `BuildEnvironmentHostUri`/`BuildTenantHostUri` append `PowerPlatformUrl` as a suffix, the `LastIndex` match also covers the sharded hosts.

### F.1 Three host-building patterns in use

| Pattern | Host expression | Example URL | Used by |
| --- | --- | --- | --- |
| **Flat / global** | `client.Api.GetConfig().Urls.PowerPlatformUrl` | `https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview` | licensing, role_based_access, capacity, application, environment_group_rule_set (one call) |
| **Environment-sharded** | `helpers.BuildEnvironmentHostUri(environmentId, ...PowerPlatformUrl)` (internal/helpers/uri.go:17-23) | `https://000000000000000000000000000000.01.environment.api.powerplatform.com/connectivity/connections?api-version=1` | connection, managed_solution |
| **Tenant-sharded** | `helpers.BuildTenantHostUri(tenantId, ...PowerPlatformUrl)` (internal/helpers/uri.go:28-34) | `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/{id}/ruleSets?api-version=2021-10-01-preview` | analytics_data_export, environment_group_rule_set, environment_groups |

### F.2 Per-service inventory

| Service | File:lines | Namespace / path | api-version | Host pattern |
| --- | --- | --- | --- | --- |
| licensing | api_licensing.go:32,49,72,101,127,143,173,200 | `/licensing/billingPolicies[...]`, `/licensing/billingPolicies/{id}/environments/{add\|remove}` | `2022-03-01-preview` | flat |
| capacity | api_capacity.go:29 | `/licensing/tenantCapacity` | `2022-03-01-preview` | flat |
| role_based_access | api_role_based_access.go:82,140,176,191 | `/authorization/roleDefinitions`, `/authorization/roleAssignments`, `/authorization/environments/{id}/roleAssignments`, `/authorization/environmentGroups/{id}/roleAssignments` | **`2024-10-01`** | flat |
| application | api_application.go:556,576,596 | `/appmanagement/applicationPackages`, `/appmanagement/environments/{id}/applicationPackages[/{name}/install]`, `/appmanagement/environments/{id}/operations/{opId}` | `2022-03-01-preview` (operations poll uses `api-version=1`) | flat |
| environment_groups | api_environment_group.go:161-163 | **`/environmentmanagement/environmentGroups/{groupId}/removeEnvironment/{environmentId}`** | `1` | tenant-sharded |
| environment_group_rule_set | api_environment_group_rule_set.go:39,76,108,134 (tenant-sharded), :155 (flat) | `/governance/environmentGroups/{id}/ruleSets`, `/governance/ruleSets/{id}` | `2021-10-01-preview` | tenant + flat |
| analytics_data_export | api_analytics_data_exports.go:41 | analytics export paths | — | tenant-sharded |
| connection | api_connection.go:34,59,83,105,125,143,191,246,264 | `/connectivity/connections`, `/connectivity/connectors/{connector}/connections/{id}/...` | `1` | environment-sharded |
| managed_solution | api_managed_solution.go:232 | `/connectivity/connections` | `1` | environment-sharded |
| environment_wave | api_environment_wave.go:34,65,114 | uses `Urls.AdminPowerPlatformUrl` (`api.admin.powerplatform.microsoft.com`) → `constants.PPAC_SCOPE`, **not** `api.powerplatform.com` | — | flat (admin host) |
| copilot_studio_application_insights | api_copilot_studio_application_insights.go:32,66,99 | **does not** call `api.powerplatform.com` — it only consumes `environment.Client.GetEnvironment` to resolve the Dataverse host | — | n/a |

Key takeaways for the migration:

* **`environment_groups` already calls a real `/environmentmanagement/...` path** (api_environment_group.go:161-163), but with `api-version=1` and against the **tenant-sharded** host — worth verifying against the documented `api.powerplatform.com/environmentmanagement/...?api-version=2024-10-01` flat form before copying it.
* `role_based_access` is the closest precedent for a flat `api.powerplatform.com` service on `api-version=2024-10-01` — follow internal/services/role_based_access/api_role_based_access.go:82 for host/query construction and internal/services/role_based_access/resource_environment_role_based_access_assignment_test.go:91-96 for the httpmock URL shape.
* No auth plumbing changes are needed: the scope resolution and the constants for all seven clouds already exist.

---

## G. Blast radius — other consumers of `environment.Client`

`environment.Client` is constructed by 12 other services. Any change to the exported method set or to `EnvironmentDto` ripples into all of them.

| Consumer | Constructs client | Methods used | Line |
| --- | --- | --- | --- |
| application (env app admin) | resource_environment_application_admin.go:115 | `GetEnvironment` | :195 |
| authorization | api_user.go:25 | `GetEnvironment` | api_user.go:35; resource_user.go:258,277,294 |
| connection | resource_connection.go:143, resource_connection_share.go:140 | `GetEnvironment` | :220 / :198 |
| copilot_studio_application_insights | api_copilot_studio_application_insights.go:22 | `GetEnvironment` | :32,:66,:99 |
| data_record | resource_data_record.go:125 | `GetEnvironment` | :182 |
| disaster_recovery | api_disaster_recovery.go:21 | `GetEnvironment` | :88; resource_disaster_recovery.go:170 |
| enterprise_policy | api_enterprise_policy.go:21 | `GetEnvironment` | resource_enterprise_policy.go:157 |
| environment_wave | api_environment_wave.go:27 | `GetEnvironment` | :44; resource_environment_wave.go:163 |
| managed_environment | api_managed_environment.go:22 | `GetEnvironment` | :122; resource_managed_environment.go:199,218,243,287,309,337 |
| powerapps | api_powerapps.go:20 | `GetEnvironments` | :30 |
| publisher | api_publisher.go:28 | `GetEnvironmentHostById` | :33,:56,:81,:101,:130,:148 |
| solution_checker_rules | client.go:25 | `GetEnvironment` | :43 |

Exported surface that must remain source-compatible (or be migrated everywhere at once):

```go
func NewEnvironmentClient(apiClient *api.Client) Client
func (client *Client) GetEnvironment(ctx, environmentId string) (*EnvironmentDto, error)
func (client *Client) GetEnvironments(ctx) ([]EnvironmentDto, error)
func (client *Client) GetEnvironmentHostById(ctx, environmentId string) (string, error)
func (client *Client) GetLocations(ctx) (*LocationArrayDto, error)
func (client *Client) LocationValidator(ctx, location, azureRegion string) error
func (client *Client) GetDefaultCurrencyForEnvironment(ctx, environmentId string) (*TransactionCurrencyDto, error)
func (client *Client) CreateEnvironment(ctx, environmentCreateDto) (*EnvironmentDto, error)
func (client *Client) UpdateEnvironment(ctx, environmentId string, EnvironmentDto) (*EnvironmentDto, error)
func (client *Client) UpdateEnvironmentAiFeatures(ctx, environmentId string, GenerativeAiFeaturesDto) error
func (client *Client) AddDataverseToEnvironment(ctx, environmentId string, createLinkEnvironmentMetadataDto) (*EnvironmentDto, error)
func (client *Client) ModifyEnvironmentType(ctx, environmentId, environmentType string) error
func (client *Client) DeleteEnvironment(ctx, environmentId string) error
func (client *Client) ValidateCreateEnvironmentDetails(ctx, location, domain string) error
func (client *Client) ValidateUpdateEnvironmentDetails(ctx, environmentId, domain string) error
```

`GetEnvironmentHostById` (api_environment.go:229-244) depends on `properties.linkedEnvironmentMetadata.instanceUrl` — any replacement API must supply an equivalent field or 6+ services break.

Other services also hard-code BAPI environment URLs independently of `environment.Client`:

* internal/services/solution/api_solution.go:400-401 — its own `GET /providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{id}`.
* internal/services/environment_groups/api_environment_group.go:132-137 — `GET .../scopes/admin/environments?$filter=properties/parentEnvironmentGroup/id eq {id}`.

---

## H. Gaps and open questions

1. **No update endpoint on the new API.** Section D shows 16 attributes that are updatable in place today. If `api.powerplatform.com/environmentmanagement` has no PATCH environment operation, the migration must either (a) keep BAPI PATCH for update while moving reads/creates, or (b) mark those attributes `RequiresReplace`, which is a **breaking change** requiring a `breaking` Changie entry.
2. **`enterprise_policies`, `cluster`/`release_cycle`, `usedBy`/`owner_id`, `updateCadence`/`cadence`, `copilotPolicies`** — need confirmation that the new API surfaces equivalents. If not, these attributes either become permanently unknown or must retain a BAPI call.
3. **`dataverse.currency_code` read path** already bypasses BAPI (Dataverse WebAPI) and is unaffected by the migration, but the 1+3N call amplification in the data source is worth fixing at the same time.
4. **`templates` JSON tag asymmetry** (`template` on read vs `templates` on write, dto.go:156 vs dto.go:230) is a latent bug that the migration should not carry forward blindly.
5. **`environment_groups` uses `api-version=1` on a tenant-sharded `/environmentmanagement/` path** — inconsistent with the documented `2024-10-01` flat form. Needs verification.
6. `BAP_2022_API_VERSION = "2022-05-01"` is defined (constants.go:198) but commented out in the update path (api_environment.go:626) — historical context for why update is pinned to `2021-04-01`.

## Clarifying questions for the user

1. Is the migration expected to be **wholesale** (drop all BAPI calls) or **hybrid** (new API for create/read/delete, BAPI retained for update and for fields the new API does not expose)?
2. Are breaking schema changes acceptable (e.g. converting today's in-place-updatable attributes to `RequiresReplace`), or must the migration be behavior-preserving for existing configurations?
3. Should the 12 downstream consumers of `environment.Client` migrate in the same PR, or should `GetEnvironment`/`GetEnvironments`/`GetEnvironmentHostById` keep their current signatures with an internal implementation swap?
4. Does the target scope include `powerplatform_locations`, `powerplatform_languages`, and `powerplatform_currencies` (they share the same BAPI `/locations` family via `internal/mocks`), or only the environment resource + data source?
