<!-- markdownlint-disable-file -->
# Attribute Impact Tables: BAPI to Power Platform API Migration

Per-attribute impact of migrating `powerplatform_environment`, `powerplatform_environments`, `powerplatform_locations`, `powerplatform_languages`, and `powerplatform_currencies` from the legacy BAPI (`api.bap.microsoft.com`) to the public Power Platform API (`api.powerplatform.com/environmentmanagement`, `api-version=2024-10-01`).

Extracted from .copilot-tracking/research/2026-08-18/powerplatform-api-migration-research.md section 3.0. See that document for endpoint mapping, BAPI source paths, the selected migration approach, and the 13 unverified items that block implementation.

Research date: 2026-08-18.

## Legend

### Status — where the data comes from

* **OK** — direct equivalent exists on the new API read model; only the JSON path changes.
* **RENAMED** — equivalent exists under a different name or shape; value semantics need confirmation.
* **WRITE-ONLY** — the field exists in the create/link request but is **not returned** by GET, so Terraform cannot detect drift or support import for it.
* **MISSING** — no equivalent anywhere in the new API.
* **NEW** — property exists only on the Power Platform API and is not currently exposed by the provider.

### Migration Result — what a practitioner experiences

* **None** — no user-visible change; values and plan behavior identical.
* **Behavior** — values identical, but plan-time or validation behavior changes. No configuration change required.
* **Verify** — mapping exists but the returned value may differ from BAPI. Reclassify after a live capture.
* **Avoidable** — nominally missing, but the value is derivable, so it can be preserved by synthesizing it.
* **Breaking: replace** — attribute becomes force-new; changing it destroys and recreates the environment.
* **Breaking: null** — attribute can no longer be read; returns `null` on refresh and import.
* **Breaking: removed** — attribute has no source at all and must be deprecated, then deleted.
* **Breaking: value** — attribute keeps working, but its value changes for existing state, producing a diff on the first refresh.
* **Additive** — new attribute; non-breaking.

> The Migration Result column describes the **end state of a full migration**. Under the recommended staged hybrid (provisioning first, BAPI retained for update and unmapped reads), Stages 1 and 2 introduce none of the `Breaking: replace` or `Breaking: null` outcomes — those land only when the read and update paths move in Stage 3.

### Resource versus data source semantics

The same status produces different outcomes depending on the surface:

| Status | In `powerplatform_environment` (resource) | In `powerplatform_environments` (data source) |
| --- | --- | --- |
| WRITE-ONLY | Value survives in state from the create request; only drift detection and import break | **Always `null`** — a data source never sends a create request, so there is no value to carry |
| OK, but no update endpoint | `Breaking: replace` — the attribute becomes force-new | **None** — data sources never plan replacements |

This makes the write-only group strictly worse for the data source than for the resource.

## Totals

| Surface | Attributes | None | Behavior | Verify | Avoidable | Breaking | NEW available |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `powerplatform_environment` | 36 | 14 | 1 | 7 | 1 | **13** | 15 |
| `powerplatform_environments` | 40 | 23 | 1 | 7 | 1 | **8** | 6 |
| `powerplatform_locations` | 10 | 6 | 0 | 1 | 0 | **3** | 5 |
| `powerplatform_languages` | 8 | 5 | 0 | 0 | 1 | **2** | 0 |
| `powerplatform_currencies` | 8 | 5 | 0 | 0 | 2 | **1** | 0 |

Breaking breakdown:

| Surface | replace | null | removed | value |
| --- | --- | --- | --- | --- |
| `powerplatform_environment` | 6 | 2 | 5 | 0 |
| `powerplatform_environments` | 0 | 3 | 5 | 0 |
| `powerplatform_locations` | 0 | 0 | 3 | 0 |
| `powerplatform_languages` | 0 | 0 | 2 | 0 |
| `powerplatform_currencies` | 0 | 0 | 1 | 0 |

These counts already credit the recovery paths verified in "Can Get Environment By Id For User fix the breaking rows?" below, which cut the environment resource from 19 breaking attributes to 13.

## Can `Get Environment By Id For User` fix the breaking rows?

Short answer: **no, but three other endpoints can.**

The question rests on a reasonable hypothesis — that the `201`/`202` provisioning response is a minimal projection and a follow-up GET would return the full object. The provisioning response *is* minimal (`OperationExecutionResult.updatedEnvironment` carries only `environmentId`, `displayName`, `dataverseOrganizationUrl`), but that is not what produced the WRITE-ONLY and MISSING classifications. Those came from the **GET response model itself** (`EnvironmentResponse`), so a GET adds nothing.

That model was re-verified against three independent code-generated renderings of the same OpenAPI spec, which agree exactly on the same 26 properties with **zero differences**:

1. The Learn REST reference `Definitions` section.
2. The Kiota-generated `Microsoft.PowerPlatform.Management.Models.EnvironmentResponse` C# class.
3. The "Power Platform for Admins V2" connector reference.

Alphabetical ordering in the generated class corroborates the absences — the sequence runs `...geo, id, protectionLevel...`, exactly where `languageCode`, `location`, and `m365Enabled` would sort. One caveat: the C# class implements `IAdditionalDataHolder` with an `AdditionalData` bag for "additional data not described in the OpenAPI description", so the SDK tolerates undocumented wire fields without naming any. Only a live capture can rule them out.

Two structural findings close the remaining escape hatches:

* `OperationExecutionResult` has **no request-echo field**, so provisioning history cannot recover create-time inputs.
* There is **no `$expand`** on either environment read operation, and `$select` only projects a subset of the documented model. This is a direct break from the provider's current BAPI `$expand=permissions,properties.capacity,properties/billingPolicy,properties/copilotPolicies`.

### Verdict per breaking attribute

| Terraform attribute | Previous result | On GET environment? | Recovery path | New result | Confidence |
| --- | --- | --- | --- | --- | --- |
| `environment_type` | Breaking: replace | **Yes** (`type`) | n/a — this is an *update* gap, not a read gap. A GET cannot perform a PATCH. | Breaking: replace | High |
| `display_name` | Breaking: replace | **Yes** (`displayName`) | Same. | Breaking: replace | High |
| `dataverse.domain` | Breaking: replace | **Yes** (`domainName`) | Same. | Breaking: replace | High |
| `dataverse.security_group_id` | Breaking: replace | **Yes** (`securityGroupId`) | Same. | Breaking: replace | High |
| `dataverse.background_operation_enabled` | Breaking: replace | **Yes** (`backgroundOperationsState`) | Same. | Breaking: replace | High |
| `dataverse.administration_mode_enabled` | Breaking: replace | **Yes** (`adminMode`) | Same. Downgradable if the 404'd `environment-state` endpoints are confirmed. | Breaking: replace | Medium |
| `billing_policy_id` | Breaking: null | No | **`GET /licensing/environments/{environmentId}/billingPolicy?api-version=2024-10-01`** → `.id`. Environment-scoped, one call. Strictly better than a reverse lookup over all billing policies. | **None** | High |
| `dataverse.language_code` | Breaking: null | No | **Dataverse Web API `organizations.languagecode`.** The provider already queries that entity for `currency_code` at internal/services/environment/api_environment.go:701-720. | **None** | High |
| `dataverse.unique_name` | Breaking: removed | No | **Dataverse Web API `organizations.uniquename`.** `domainName` on the new API is a different value and is not a substitute. | **None** | Medium |
| `owner_id` | Breaking: null | No as `usedBy` | `createdFor.id` is the probable rename; `usedBy` survives only as a create-request field. | **Verify** | Medium |
| `enterprise_policies[].location` | Breaking: removed | No | Parse the region from `resourceId`, or resolve via ARM `GET management.azure.com{resourceId}` → `location` (different host and token audience). | **Avoidable** | Medium |
| `enterprise_policies[].system_id` | Breaking: value | Indirect | The new `resourceId` is a full ARM ID matching BAPI `id`, and BAPI `systemId` most likely matches the new `id` — that is, the pairing in earlier research was inverted. | **Verify** | Low |
| `description` | Breaking: null | **No** | None. Write-only on Provision and Reset. | Breaking: null | High |
| `dataverse.templates` | Breaking: null | **No** | Only a lossy proxy: `GET /appmanagement/environments/{id}/applicationPackages` filtered to `state=TemplateInstalled`. Not equivalent to the create-time `templates[]`. | Breaking: null | High |
| `dataverse.template_metadata` | None (resource) | **No** | None. | None (resource) / Breaking: null (data source) | High |
| `cadence` | Breaking: removed | **No** | None. The only `UpdateCadence` in the namespace is `FinOpsUpdateCadence`, which tracks Finance and Operations app versions and is unrelated. | Breaking: removed | High |
| `allow_bing_search` | Breaking: removed | **No** | None — absent from the 33-property `EnvironmentManagementSetting` model too. | Breaking: removed | High |
| `allow_microsoft_365_services` | Breaking: removed | **No** | None. | Breaking: removed | High |
| `allow_moving_data_across_regions` | Breaking: removed | **No** | None — no `copilotPolicies` anywhere in the Models namespace. | Breaking: removed | High |
| `allow_flex_routing` | Breaking: removed | **No** | None. | Breaking: removed | High |

### Net effect

| | Before this analysis | After |
| --- | --- | --- |
| Breaking attributes on `powerplatform_environment` | 19 | **13** |
| of which `replace` | 6 | 6 (unchanged — read paths cannot fix update gaps) |
| of which `null` | 5 | 2 |
| of which `removed` | 7 | 5 |
| of which `value` | 1 | 0 (reclassified to Verify pending a live capture) |

The three recovered attributes cost two extra calls per environment read, one of which the provider already makes:

1. `GET /licensing/environments/{environmentId}/billingPolicy` — new call.
2. Dataverse Web API `organizations` — **already called** for `currency_code`; extend the existing `$select` to add `languagecode` and `uniquename`.

Eight attributes remain genuinely unrecoverable without BAPI: `description`, `cadence`, `dataverse.templates`, `dataverse.template_metadata`, and the four `allow_*` flags. `cadence` is the worst of these because it is absent from `CreateEnvironmentRequest` as well as the response, making it a **capability removal** rather than a refresh gap.

## `powerplatform_environment` (resource)

Schema: internal/services/environment/resource_environment.go:63-377

| Status | Terraform attribute | Migration Result | Power Platform API path | Note |
| --- | --- | --- | --- | --- |
| OK | `id` | None | `id` | Same GUID; ARM-style `id` disappears. |
| OK | `location` | None | `geo` | Renamed; verify `unitedstates`/`europe` values still match. |
| OK | `azure_region` | Behavior | `azureRegion` | Value survives, but plan-time validation is lost with `azureRegions[]`. Invalid values now fail at apply instead of plan. |
| OK | `environment_type` | **Breaking: replace** | `type` | No `modifySku` equivalent, so a SKU change destroys and recreates. New API documents 12 SKUs vs the 5 allowed today. |
| OK | `display_name` | **Breaking: replace** | `displayName` | No PATCH, so a rename destroys and recreates the environment. Highest-impact break in the set. |
| OK | `release_cycle` | None | `clusterCategory` (read), `cluster.category` (create) | Already `RequiresReplace`, so no update loss. |
| OK | `environment_group_id` | None | `environmentGroupId` (read), `parentEnvironmentGroup.id` (create) | Update preserved via `environmentGroups/{g}/addEnvironment/{e}`. |
| OK | `enterprise_policies[].type` | None | synthesized | Gains a fourth value, `PrivateEndpoint`. |
| OK | `enterprise_policies[].id` | Verify | `enterprisePolicies.*.id` | **Likely swapped with `system_id`** — new `resourceId` is the full ARM ID, which matches BAPI `id`, while BAPI `systemId` matches the new `id`. UNVERIFIED. |
| OK | `enterprise_policies[].status` | Verify | `enterprisePolicies.*.status` | Was `linkStatus`; enum now published (6 values). Confirm the value set matches. |
| OK | `dataverse.url` | None | `url` | Critical — `GetEnvironmentHostById` plus 6+ services. |
| OK | `dataverse.domain` | **Breaking: replace** | `domainName` | In-place update lost. |
| OK | `dataverse.organization_id` | None | `dataverseId` | |
| OK | `dataverse.version` | None | `version` | |
| OK | `dataverse.security_group_id` | **Breaking: replace** | `securityGroupId` (top level) | In-place update lost. |
| OK | `dataverse.currency_code` | None | unchanged (Dataverse WebAPI) | Read already bypasses BAPI. |
| OK | `dataverse.background_operation_enabled` | **Breaking: replace** | `backgroundOperationsState` | In-place update lost. |
| OK | `dataverse.administration_mode_enabled` | **Breaking: replace** | `adminMode` | Downgradable to None if the 404'd `environment-state/enable` and `disable` endpoints are confirmed. Enum UNVERIFIED. |
| OK | `timeouts` | None | n/a | Framework-only. |
| RENAMED | `enterprise_policies[].system_id` | Verify | `enterprisePolicies.*.resourceId` (likely `enterprisePolicies.*.id`) | The new `resourceId` is a full ARM resource ID, which corresponds to BAPI `id`, not `systemId`. The BAPI `systemId` GUID most likely maps to the new `id`. Resolve the pairing with a live capture before wiring either. |
| RENAMED | `dataverse.linked_app_type` | Verify | `finOpsMetadata.type` | Value equivalence UNVERIFIED. Becomes `Breaking: value` if the values differ. |
| RENAMED | `dataverse.linked_app_id` | Verify | `finOpsMetadata.id` | Value equivalence UNVERIFIED. |
| RENAMED | `dataverse.linked_app_url` | Verify | `finOpsMetadata.url` | Value equivalence UNVERIFIED. |
| WRITE-ONLY | `description` | **Breaking: null** | `description` (create request) | Not returned by GET, so drift is undetectable, update is impossible, and import returns `null`. |
| WRITE-ONLY | `owner_id` | Verify | `createdFor.id` | Not on the GET model as `usedBy`, but `createdFor` is the probable rename. Confirm the value matches BAPI `properties.usedBy.id` before relying on it. |
| WRITE-ONLY | `billing_policy_id` | None | — (recovered) | **Recovered**: `GET /licensing/environments/{environmentId}/billingPolicy?api-version=2024-10-01` returns `id`. Environment-scoped, one extra call. |
| WRITE-ONLY | `dataverse.language_code` | None | — (recovered) | **Recovered**: Dataverse Web API `organizations.languagecode`. The provider already calls that entity for `currency_code`, so this is nearly free. |
| WRITE-ONLY | `dataverse.templates` | **Breaking: null** | `templates` (create/link) | Create path already stitches from plan, so only import breaks. |
| WRITE-ONLY | `dataverse.template_metadata` | None | `templateMetadata` (create/link, bare `object`) | Already `Optional`-only (not `Computed`), so it is never read back today. `PostProvisioningPackages` shape UNVERIFIED. |
| MISSING | `cadence` | **Breaking: removed** | — | No read and no create field. Frequent/Moderate cadence becomes unreachable entirely. |
| MISSING | `allow_bing_search` | **Breaking: removed** | — | Not on `/environments/{id}/settings` either. |
| MISSING | `allow_microsoft_365_services` | **Breaking: removed** | — | Same. |
| MISSING | `allow_moving_data_across_regions` | **Breaking: removed** | — | Same. Also removes the input to `aiGenerativeFeaturesValidaor`. |
| MISSING | `allow_flex_routing` | **Breaking: removed** | — | Same. |
| MISSING | `enterprise_policies[].location` | Avoidable | — | Parse the Azure region out of `enterprisePolicies.*.resourceId`, or resolve it with an ARM `GET management.azure.com{resourceId}` (different host and token audience). |
| MISSING | `dataverse.unique_name` | None | — (recovered) | **Recovered**: Dataverse Web API `organizations.uniquename`. Note `domainName` on the new API is a *different* value and cannot be substituted. |

### New attributes available to add

| Status | Proposed attribute | Migration Result | Power Platform API path | Kind |
| --- | --- | --- | --- | --- |
| NEW | `macro_region` | Additive | `macroRegion` (create request) | Optional + RequiresReplace; needs `ConflictsWith(location)`; not read back. |
| NEW | `connected_group_id_for_teams_environment` | Additive | create request | Optional + RequiresReplace. |
| NEW | `currency_name`, `currency_symbol`, `currency_precision` | Additive | `linkedEnvironmentMetadata.currency.{name,symbol,precision}` | Optional; BAPI create accepted only `code`. |
| NEW | `state` | Additive | `state` | Computed; enum UNVERIFIED. |
| NEW | `tenant_id` | Additive | `tenantId` | Computed; in BAPI but never surfaced. |
| NEW | `connected_group_id` | Additive | `connectedGroupId` | Computed. |
| NEW | `scenario_name` | Additive | `scenarioName` | Computed. |
| NEW | `protection_level` | Additive | `protectionLevel` | Computed; overlaps `powerplatform_managed_environment`. |
| NEW | `created_date_time` | Additive | `createdDateTime` | Computed. |
| NEW | `created_by` (`id`, `type`) | Additive | `createdBy` | Computed. |
| NEW | `created_for` (`id`, `type`) | Additive | `createdFor` | Computed; **no BAPI equivalent**. |
| NEW | `deleted_date_time` | Additive | `deletedDateTime` | Computed; **no BAPI equivalent**. |
| NEW | `retention_details` (`retention_period`, `available_from_date_time`) | Additive | `retentionDetails` | Computed. |
| NEW | `enterprise_policies[].type = "PrivateEndpoint"` | Additive | `enterprisePolicies.privateEndpoint` | Computed; **no BAPI equivalent**. |
| NEW | `enterprise_policies[].error` | Additive | `enterprisePolicies.*.error` | Computed; **no BAPI equivalent**. |

## `powerplatform_environments` (data source)

Schema: internal/services/environment/datasource_environments.go:48-232

Element shape is the identical `SourceModel` (internal/services/environment/models.go:33-36), so every row in the resource table applies to `environments[*].<attribute>` — but with **two systematic reclassifications**:

| Resource rows | Resource result | Data source result | Why |
| --- | --- | --- | --- |
| `environment_type`, `display_name`, `dataverse.domain`, `dataverse.security_group_id`, `dataverse.background_operation_enabled`, `dataverse.administration_mode_enabled` | Breaking: replace | **None** | Data sources never plan replacements, so the missing update endpoint is irrelevant. |
| `dataverse.template_metadata` | None | **Breaking: null** | It is `Computed` in the data source (datasource_environments.go:223-226) but `Optional`-only in the resource, so the data source does read it back today. |

All five remaining `Breaking: removed` rows apply unchanged. `description` and `dataverse.templates` stay `Breaking: null` and are strictly worse here, since a data source has no create request to carry the value forward. The three recovered attributes (`billing_policy_id`, `dataverse.language_code`, `dataverse.unique_name`) are `None` on both surfaces — though `billing_policy_id` costs one extra call **per environment**, which compounds the data source's existing 1 + 3N amplification.

Data-source-only rows:

| Status | Terraform attribute | Migration Result | Power Platform API path | Note |
| --- | --- | --- | --- | --- |
| OK | `environments` | None | `value[]` | Envelope unchanged (`value` plus `@odata.nextlink`). |
| OK | `environments[*].timeouts` | None | n/a | Vestigial — present only because `SourceModel` carries `Timeouts`. |
| OK | `timeouts` | None | n/a | Read-only timeouts. |
| NEW | `ids` | Additive | `ids` query param | Optional filter input. |
| NEW | `filter` | Additive | `$filter` on `dataverseId`, `type`, `geo`, `state`, `environmentGroupId`, `domainName` | Optional filter input. |
| NEW | `top`, `skip`, `orderby` | Additive | `$top`, `$skip`, `$orderby` | Optional paging inputs; `@odata.nextlink` must be followed. |

## `powerplatform_locations`

Schema: internal/services/locations/datasource_locations.go:49-104

| Status | Terraform attribute | Migration Result | Power Platform API path | Note |
| --- | --- | --- | --- | --- |
| OK | `locations[].name` | None | `collection[].name` | |
| OK | `locations[].display_name` | None | `collection[].displayName` | |
| OK | `locations[].code` | Verify | `collection[].code` | Value may change (`"NA"` vs `"NAM"`) — UNVERIFIED. Becomes `Breaking: value` if it differs. |
| OK | `locations[].is_default` | None | `collection[].isDefault` | |
| OK | `locations[].is_disabled` | None | `collection[].isDisabled` | |
| OK | `locations[].can_provision_database` | None | `collection[].canProvisionDatabase` | |
| OK | `timeouts` | None | n/a | |
| MISSING | `locations[].id` | **Breaking: removed** | — | Synthesizable, but fabricating a BAPI ARM path on a non-ARM API is not advisable. |
| MISSING | `locations[].can_provision_customer_engagement_database` | **Breaking: removed** | — | |
| MISSING | `locations[].azure_regions` | **Breaking: removed** | — | Also removes `powerplatform_environment.azure_region` plan-time validation, which is why that attribute is classified `Behavior`. |
| NEW | `locations[].has_first_release_island_available_for_provisioning` | Additive | `collection[].hasFirstReleaseIslandAvailableForProvisioning` | Directly relevant to `release_cycle = "Early"`. |
| NEW | `location_selection_mode` | Additive | `locationSelectionMode` | `Region` or `MacroRegion`. |
| NEW | `macro_regions[].macro_region_id` | Additive | `macroRegions[].macroRegionId` | |
| NEW | `macro_regions[].display_name` | Additive | `macroRegions[].displayName` | |
| NEW | `macro_regions[].data_residency_note` | Additive | `macroRegions[].dataResidencyNote` | |

## `powerplatform_languages`

Schema: internal/services/languages/datasource_languages.go:44-90

| Status | Terraform attribute | Migration Result | Power Platform API path | Note |
| --- | --- | --- | --- | --- |
| OK | `location` (input) | None | `{location}` path segment | |
| OK | `languages[].localized_name` | None | `collection[].localizedName` | |
| OK | `languages[].locale_id` | None | `collection[].localeId` | |
| OK | `languages[].is_tenant_default` | None | `collection[].isTenantDefault` | |
| OK | `timeouts` | None | n/a | |
| MISSING | `languages[].id` | **Breaking: removed** | — | Synthesizable, but fabricating a BAPI ARM path on a non-ARM API is not advisable. |
| MISSING | `languages[].name` | Avoidable | — | Preserve by synthesizing `string(localeId)`. `environment.languageCodeValidator` matches on this (internal/services/environment/api_environment.go:183-227) and must be rewritten against `localeId`. |
| MISSING | `languages[].display_name` | **Breaking: removed** | — | Only `localizedName` survives, and the two values differ for non-English locales, so it cannot be aliased. |

No new properties — the new language object is strictly smaller.

## `powerplatform_currencies`

Schema: internal/services/currencies/datasource_currencies.go:49-95

| Status | Terraform attribute | Migration Result | Power Platform API path | Note |
| --- | --- | --- | --- | --- |
| OK | `location` (input) | None | `{location}` path segment | |
| OK | `currencies[].code` | None | `collection[].code` | |
| OK | `currencies[].symbol` | None | `collection[].symbol` | |
| OK | `currencies[].is_tenant_default` | None | `collection[].isTenantDefault` | |
| OK | `timeouts` | None | n/a | |
| MISSING | `currencies[].id` | **Breaking: removed** | — | Synthesizable, but fabricating a BAPI ARM path on a non-ARM API is not advisable. |
| MISSING | `currencies[].name` | Avoidable | — | Equal to `code` in practice, so it can be preserved. `environment.currencyCodeValidator` matches on this (internal/services/environment/api_environment.go:119-163) and must be rewritten against `code`. |
| MISSING | `currencies[].type` | Avoidable | — | Was the constant ARM discriminator, so it can be preserved as a literal — though it is meaningless on the new API and is a good deprecation candidate. |

No new properties on the list endpoint. `name`, `symbol`, and `precision` exist only on the create-request currency model (`EnvironmentRequestCurrency`).

## Cross-cutting consequences

* Two imperative validators are hard-coupled to attributes that vanish. `languageCodeValidator` matches on `languages[].name` and `currencyCodeValidator` matches on `currencies[].name`. Both are rewritable against `localeId` and `code` respectively.
* `LocationValidator`'s azure-region check has no replacement — `azureRegions[]` does not exist on the new locations endpoint. Either retain a single BAPI `/locations` call for that check, or drop client-side validation and rely on `errorDetail.fieldErrors.{field}.suggestedValue` from the provisioning call.
* `dataverse.url` survives as `url`, so `GetEnvironmentHostById` (internal/services/environment/api_environment.go:229-244) and its 6+ dependent services keep working.
* internal/mocks/mocks.go:69-82 registers shared responders for the three BAPI reference endpoints, and `RegisterNoResponder` at :47 fails fast on any unregistered URL. Data source URL changes must land in the same change as the mock updates.

## Breaking changes ranked by practitioner impact

1. **`display_name` becomes force-new** — renaming an environment would destroy and recreate it, taking the Dataverse database with it. The single most dangerous outcome in the set, and one that no read endpoint can fix.
2. **`environment_type` becomes force-new** — a Sandbox-to-Production promotion would destroy and recreate.
3. **`cadence` removed as a capability** — absent from both the response and `CreateEnvironmentRequest`, so Frequent/Moderate release cadence becomes unsettable, not merely unmanageable.
4. **The four `allow_*` flags removed** — no equivalent anywhere in the namespace, including the 33-property `EnvironmentManagementSetting` model.
5. **`dataverse.domain` and `dataverse.security_group_id` become force-new** — both are routinely changed in place today.
6. **`description` and `dataverse.templates` return `null` on refresh and import** — always `null` in `powerplatform_environments`.
7. **Reference data source removals** — `locations[].azure_regions` has the widest blast radius because it also disables environment region validation.

Mitigations that reduce the count, in order of value:

* Keep `Update` on BAPI (the recommended staged hybrid) — eliminates all six `Breaking: replace` outcomes.
* Adopt the three verified recovery paths — `GET /licensing/environments/{id}/billingPolicy` plus two extra columns on the Dataverse `organizations` query the provider already makes. Eliminates three outcomes for the cost of one new call.
* Keep one supplemental BAPI GET for the eight genuinely unrecoverable fields — eliminates the remaining two `null` and five `removed` outcomes.
* Synthesize `languages[].name`, `currencies[].name`, and `currencies[].type` — eliminates three reference-data removals.
* Deprecate rather than delete anything that still has to go, returning `null` with a `DeprecationMessage` for one major cycle.

## Open verification items from this analysis

| # | Item | Blocks | How to verify |
| --- | --- | --- | --- |
| 1 | Dump a raw `GET /environmentmanagement/environments/{id}` body and diff against the 26 documented properties | Upgrades every "No" verdict from documented absence to observed absence | mitmproxy capture |
| 2 | Does `$select=description` error, or silently return nothing? | Confirms `$select` cannot reach undocumented fields | Live probe |
| 3 | Does `createdFor.id` equal BAPI `properties.usedBy.id`? | `owner_id` verdict | Side-by-side BAPI and new-API capture |
| 4 | Does new `id` equal BAPI `systemId`, and new `resourceId` equal BAPI `id`? | `enterprise_policies[].id` and `.system_id` pairing | Capture one CMK or VNet-linked environment on both APIs |
| 5 | What does `GET /licensing/environments/{id}/billingPolicy` return when no policy is linked? | Must distinguish "no policy" from "permission denied" before wiring it in | Live probe |
