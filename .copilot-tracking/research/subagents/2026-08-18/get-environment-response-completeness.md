<!-- markdownlint-disable-file -->
# Subagent Research: `GET /environmentmanagement/environments/{id}` response completeness

Research topic: for each of 13 Terraform attributes currently classified as "unavailable on the Power Platform API", determine definitively whether `GET https://api.powerplatform.com/environmentmanagement/environments/{environmentId}?api-version=2024-10-01` (Get Environment By Id For User) — or its list sibling — actually returns it, possibly renamed, possibly only via `$select`, possibly undocumented. Plus item 14: is `enterprisePolicies.*.resourceId` the same value as BAPI `properties.enterprisePolicies.*.systemId`?

Status: **COMPLETE for everything resolvable from published sources.** Three independent code-generated renderings of the same OpenAPI spec were cross-checked and agree exactly. Two sub-questions are flagged UNVERIFIED — requires live capture.

Research date: 2026-08-18.

Prior research consumed:

* .copilot-tracking/research/subagents/2026-08-18/powerplatform-api-environmentmanagement.md (section 4.3 — the Learn-derived `EnvironmentResponse`)
* .copilot-tracking/research/2026-08-18/powerplatform-api-attribute-impact-tables.md

---

## 1. Verdict table (items 1-14)

| # | Terraform attribute | Available via GET environment? | Exact JSON path if yes | Alternative endpoint that could supply it | Confidence | Source |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | `description` | **No** | — | None on `api.powerplatform.com`. Write-only input on Provision New Environment (`description`) and Reset Environment (`description`). Speculative: ARG `resourcequery` free-form `properties` bag (UNVERIFIED). | High (absence), Low (ARG alt) | [EnvironmentResponse definition](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user); [`EnvironmentResponse` .NET class](https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.models.environmentresponse?view=power-platform-latest); [connector `EnvironmentResponse`](https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/) |
| 2 | `owner_id` (BAPI `usedBy.id`) | **No** — but a near-equivalent exists | `createdFor.id` (+ `createdFor.type`) | `createdFor` is the closest semantic match on the new model. `usedBy` itself exists only as a **request** field (`UserIdentity`) on Provision New Environment. | Medium (semantic equivalence UNVERIFIED) | [EnvironmentResponse definition](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user); [Provision New Environment](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/provision-new-environment) |
| 3 | `billing_policy_id` | **No** | — | **YES — solved.** `GET https://api.powerplatform.com/licensing/environments/{environmentId}/billingPolicy?api-version=2024-10-01` → `BillingPolicyResponseModel.id`. Also `GET /licensing/billingPolicies/{billingPolicyId}/environments`. | **High** | [Get Environment Billing Policy](https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/environment-billing-policy/get-environment-billing-policy) |
| 4 | `dataverse.language_code` | **No** | — | Dataverse Web API `GET {envUrl}/api/data/v9.0/organizations` → `languagecode` (int, e.g. `1033`). Provider already calls this endpoint in internal/services/environment_settings/api_environment_settings.go. | High | [EnvironmentResponse definition](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user); internal/services/environment_settings/tests/datasource/organisations.json (line 117 `"languagecode": 1033`) |
| 5 | `dataverse.templates` | **No** | — | Lossy: `GET /appmanagement/environments/{environmentId}/applicationPackages?api-version=2024-10-01` → `value[].uniqueName` where `value[].state` is `Installed` / `TemplateInstalled`. These are app-package unique names, **not** BAPI provisioning template codes (`D365_Sales` etc.). | Low | [Get Environment Application Package](https://learn.microsoft.com/en-us/rest/api/power-platform/appmanagement/applications/get-environment-application-package); connector `ApplicationPackage` / `InstancePackageState` |
| 6 | `dataverse.template_metadata` | **No** | — | **None.** Write-only (`linkedEnvironmentMetadata.templateMetadata` on Provision New Environment / Link Dataverse / Reset Environment). Not echoed by any read endpoint, including operations history. | High | [Provision New Environment](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/provision-new-environment); [OperationExecutionResult](https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/) |
| 7 | `dataverse.unique_name` | **No** | — | Dataverse Web API `GET {envUrl}/api/data/v9.0/organizations` → `uniquename` (e.g. `unqa7a3ad3827f8ee1190486045bd29e`). Note `EnvironmentResponse.domainName` is a **different** value (`org0fadb1dd`). | Medium (Dataverse column documented; not present in this repo's trimmed fixture) | [Dataverse `organization` table reference](https://learn.microsoft.com/en-us/power-apps/developer/data-platform/reference/entities/organization#BKMK_UniqueName); internal/services/environment/tests/resource/Validate_Create_Environment_And_Dataverse/get_lifecycle_new_dataverse.json |
| 8 | `cadence` (BAPI `updateCadence.id`) | **No** | — | **None.** The only `*UpdateCadence` in the whole Power Platform API surface is `FinOpsUpdateCadence` on `GET /dynamics/environments/{id}` maintenance settings — a Finance & Operations *application version* cadence, semantically unrelated to the Dataverse release channel (`Frequent` / `Moderate`). No `governance` equivalent. | High | [Microsoft.PowerPlatform.Management.Models namespace](https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.models?view=power-platform-latest); [Get Fin Ops Maintenance Settings](https://learn.microsoft.com/en-us/rest/api/power-platform/dynamics/finance-and-operations-maintenance-settings/get-fin-ops-maintenance-settings) |
| 9 | `allow_bing_search` (`bingChatEnabled`) | **No** | — | **None.** Not in `EnvironmentResponse`; not in `EnvironmentManagementSetting` (full 33-property list re-enumerated in §5.1 — no Bing/chat toggle). | High | [List Environment Management Settings](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/list-environment-management-settings); connector `EnvironmentManagementSetting` definition |
| 10 | `allow_microsoft_365_services` (`m365Enabled`) | **No** | — | **None.** Same as #9. | High | Same as #9 |
| 11 | `allow_moving_data_across_regions` (`copilotPolicies.crossGeoCopilotDataMovementEnabled`) | **No** | — | **None.** No `copilotPolicies` object anywhere in the `Microsoft.PowerPlatform.Management.Models` namespace. | High | [Models namespace listing](https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.models?view=power-platform-latest) |
| 12 | `allow_flex_routing` (`copilotPolicies.crossBoundaryCopilotDataMovementEnabled`) | **No** | — | **None.** Same as #11. | High | Same as #11 |
| 13 | `enterprise_policies[].location` | **No** | — | Derivable out-of-band: parse the Azure region out of `enterprisePolicies.*.resourceId`, or `GET https://management.azure.com{resourceId}?api-version=2020-10-30-preview` → `location`. Different host, different token audience. | Medium | [EnterprisePolicyLink definition](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user) |
| 14 | `enterprisePolicies.*.resourceId` == BAPI `systemId`? | n/a | `enterprisePolicies.{encryption\|identity\|networkInjection\|privateEndpoint}.resourceId` | **Almost certainly NOT `systemId`.** Learn describes `resourceId` as "The fully-qualified Azure resource ID of the enterprise policy" — that is BAPI's `enterprisePolicies.*.id`. BAPI's `systemId` is the enterprise policy's `properties.systemId` (a `/regions/{geo}/providers/Microsoft.PowerPlatform/enterprisePolicies/{guid}` path). The new `EnterprisePolicyLink.id` ("The ID of the enterprise policy") is the likelier `systemId`/`policyId` counterpart. | **Low — UNVERIFIED, requires live capture** | [EnterprisePolicyLink definition](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user); internal/services/environment/dto.go lines 100-106 |

**Net result: 1 of 13 (`billing_policy_id`) is fully recoverable from another `api.powerplatform.com` endpoint. 3 more (`dataverse.language_code`, `dataverse.unique_name`, `enterprise_policies[].location`) are recoverable from a non-Power-Platform-API surface the provider already touches or could touch. 1 (`owner_id`) has a near-equivalent field with a semantic caveat. The remaining 8 are genuinely unavailable.**

---

## 2. `EnvironmentResponse` — full property list from higher-fidelity sources

### 2.1 The three sources cross-checked

All three are generated from the same internal OpenAPI description. There is **no public OpenAPI/Swagger document** for this surface (independently confirmed by `microsoft/Employee-Self-Service-Agent-Developer-Kit`, which classes it "documented tier — no public OpenAPI and no no-auth `$metadata`" in `tests/fixtures/cassettes/INDEX.md`).

| Source | What it is | Link |
| --- | --- | --- |
| **A** — Learn REST "Definitions" | Rendered from the OpenAPI `definitions` block | https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user |
| **B** — .NET SDK class | Kiota-generated C# class in `Microsoft.PowerPlatform.Management` v2.0.3474.290 | https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.models.environmentresponse?view=power-platform-latest |
| **C** — Connector reference | "Power Platform for Admins V2" Swagger, rendered by Learn (updated 2026-08-07) | https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ |

### 2.2 Consolidated property list (26 top-level / leaf paths)

| # | JSON path | Type | In A | In B | In C |
| --- | --- | --- | :-: | :-: | :-: |
| 1 | `id` | string | Y | Y (`Id`) | Y |
| 2 | `displayName` | string | Y | Y | Y |
| 3 | `type` | string (SKU) | Y | Y | Y |
| 4 | `state` | string | Y | Y | Y |
| 5 | `tenantId` | string | Y | Y | Y |
| 6 | `geo` | string | Y | Y | Y |
| 7 | `azureRegion` | string | Y | Y | Y |
| 8 | `clusterCategory` | string | Y | Y | Y |
| 9 | `adminMode` | string | Y | Y | Y |
| 10 | `backgroundOperationsState` | string | Y | Y | Y |
| 11 | `protectionLevel` | string | Y | Y | Y |
| 12 | `environmentGroupId` | string | Y | Y | Y |
| 13 | `connectedGroupId` | string | Y | Y | Y |
| 14 | `securityGroupId` | string | Y | Y | Y |
| 15 | `scenarioName` | string | Y | Y | Y |
| 16 | `dataverseId` | string | Y | Y | Y |
| 17 | `url` | string | Y | Y | Y |
| 18 | `domainName` | string | Y | Y | Y |
| 19 | `version` | string | Y | Y | Y |
| 20 | `createdDateTime` | date-time | Y | Y | Y |
| 21 | `createdBy` → `{id, type}` | `EnvironmentPrincipal` | Y | Y | Y |
| 22 | `createdFor` → `{id, type}` | `EnvironmentPrincipal` | Y | Y | Y |
| 23 | `deletedDateTime` | date-time | Y | Y | Y |
| 24 | `retentionDetails` → `{retentionPeriod, availableFromDateTime}` | `RetentionDetails` | Y | Y | Y (flattened as `retentionDetails.*`) |
| 25 | `finOpsMetadata` → `{id, type, url}` | `FinOpsMetadata` | Y | Y | Y |
| 26 | `enterprisePolicies` → `{encryption, identity, networkInjection, privateEndpoint}` each `{id, resourceId, status, error}` | `EnterprisePolicies` | Y | Y | Y |

### 2.3 Diff against the Learn-derived list in prior research

**Zero differences. The SDK revealed no properties missing from the prior Learn-derived list, and no properties in the prior list are absent from the SDK.**

The one structural nuance worth recording: source B's C# class exposes an extra member `AdditionalData` — "Stores additional data not described in the OpenAPI description found when deserializing." This is the standard Kiota `IAdditionalDataHolder` escape hatch, present on **every** model in the package. It proves the SDK would tolerate undocumented wire properties, but it does **not** name any. It is therefore evidence that the SDK cannot rule out undocumented fields — not evidence that undocumented fields exist.

> Source: https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.models.environmentresponse?view=power-platform-latest ("Properties" table, `AdditionalData` row; class declared `public class EnvironmentResponse : IAdditionalDataHolder, IParsable`).

### 2.4 Explicit absence check for the 13 target names

Searched all three sources for: `description`, `usedBy`, `owner`, `billingPolicy`, `baseLanguage`, `baseLanguageCode`, `languageCode`, `template`, `templates`, `templateMetadata`, `uniqueName`, `updateCadence`, `cadence`, `releaseCadence`, `bingChatEnabled`, `bingChat`, `m365Enabled`, `copilotPolicies`, `crossGeoCopilotDataMovementEnabled`, `crossBoundaryCopilotDataMovementEnabled`, `location` (on `EnterprisePolicyLink`).

**None of these names appear on `EnvironmentResponse` or on any of its nested types (`EnvironmentPrincipal`, `RetentionDetails`, `FinOpsMetadata`, `EnterprisePolicies`, `EnterprisePolicyLink`) in any of the three sources.**

Corroborating negative evidence: the alphabetical property tables in sources A and B both run `...geo, id, protectionLevel...` with nothing between `id` and `protectionLevel` — which is exactly where `languageCode`, `location`, and `m365Enabled` would sort. Same for the `c`/`d` block: `clusterCategory, connectedGroupId, createdBy, createdDateTime, createdFor, dataverseId, deletedDateTime, displayName, domainName` — no `cadence`, `copilotPolicies`, or `description`.

### 2.5 Related models that are NOT `EnvironmentResponse` (do not confuse)

The connector reference exposes a second, much smaller model literally named `Environment`:

| Model | Properties | Used by |
| --- | --- | --- |
| `Environment` ("Power Platform environment") | `environmentId`, `displayName`, `dataverseOrganizationUrl` — **3 properties only** | `OperationExecutionResult.updatedEnvironment`, `EnvironmentPagedCollection.collection` (copy/restore candidates) |
| `EnvironmentResponse` | the 26 above | `GET /environmentmanagement/environments/{id}`, `EnvironmentList.value[]` |

> Source: https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ → Definitions → `Environment`, `EnvironmentPagedCollection`, `OperationExecutionResult`.

This matters for item 5/6/1/2/4 — see §5.4.

---

## 3. `$select` / `$expand` semantics

### 3.1 Query parameters, exhaustively

| Operation | Parameters |
| --- | --- |
| `GET /environmentmanagement/environments` (List Environments For User) | `api-version` (required), `ids`, `$filter`, `$select`, `$top`, `$skip`, `$orderby` |
| `GET /environmentmanagement/environments/{environmentId}` (Get Environment By Id For User) | `api-version` (required), `$select` |

Sources:

* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user — URI Parameters table lists exactly `environmentId`, `api-version`, `$select`.
* https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.environmentmanagement.environments.environmentsrequestbuilder.environmentsrequestbuildergetqueryparameters?view=power-platform-latest — generated query-parameter class exposes exactly `ApiVersion`, `Filter`, `Ids`, `Orderby`, `Select`, `Skip`, `Top`.
* https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ → "Retrieve a list of environments (preview)" and "Retrieves a single environment by ID (preview)" — same parameter sets.

### 3.2 Findings

* **`$select` is projection-only.** Its documented description in all three sources is verbatim "Comma-separated list of properties to include in the response." A projection operator narrows the default representation; there is no documented mechanism by which it could widen it. Compare with the `governance` namespace's `GET /governance/ruleSets` operation, which explicitly documents a **separate** `$expand` parameter alongside `$select` — proving the spec authors do model `$expand` when the endpoint supports it, and did not here.
* **There is no `$expand` on either environment read operation.** This is a direct behavioural break from BAPI, where the provider today issues `$expand=permissions,properties.capacity,properties/billingPolicy,properties/copilotPolicies` (internal/services/environment/api_environment.go line 253) to pull in `billingPolicy` and `copilotPolicies`. That expansion has no counterpart.
* **`$filter` supports only six properties**: `dataverseId`, `type`, `geo`, `state`, `environmentGroupId`, `domainName`. All six are in the documented 26. This is corroborating evidence that the server-side model *is* the documented 26 — a filterable-property allowlist drawn exclusively from the documented set.
* **UNVERIFIED**: whether passing `$select=description` returns `400 Bad Request`, silently ignores the token, or returns an undocumented value. Only a live call can settle this. It is the single cheapest experiment available and would definitively close items 1, 2, 4-12.

---

## 4. Alternative endpoints — exact URL templates and response shapes

### 4.1 Billing policy for an environment (solves item 3) — HIGH confidence

```http
GET https://api.powerplatform.com/licensing/environments/{environmentId}/billingPolicy?api-version=2024-10-01
```

Response `200 OK` — `BillingPolicyResponseModel`:

| JSON path | Type | Notes |
| --- | --- | --- |
| `id` | string | **This is `billing_policy_id`.** |
| `name` | string | |
| `status` | `BillingPolicyStatus` enum: `Enabled`, `Disabled` | |
| `location` | string | |
| `billingInstrument.subscriptionId` | string (uuid) | |
| `billingInstrument.resourceGroup` | string | |
| `billingInstrument.id` | string | |
| `createdOn` / `lastModifiedOn` | date-time | |
| `createdBy` / `lastModifiedBy` | `LicensingPrincipal` `{id, type}` where `type` ∈ `None`, `Application`, `User`, `DelegatedAdmin` | |

Other responses: `400`, `401`, `403`. Note **no documented `404`** — behaviour when an environment has no linked billing policy is UNVERIFIED (likely `404` or an empty/`null` body; the `microsoft/Employee-Self-Service-Agent-Developer-Kit` mock treats `404` on the sibling `billingPolicies/{id}/environments` route as "none linked").

* Doc: https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/environment-billing-policy/get-environment-billing-policy
* SDK confirmation of the path: https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.licensing.environments.item.billingpolicy.billingpolicyrequestbuilder?view=power-platform-latest — class summary reads verbatim `Builds and executes requests for operations under \licensing\environments{environmentId}\billingPolicy`, method `GetAsync` documented "Get the linked billing policy details for an environment."
* Connector confirmation: https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ → operation ID `GetEnvironmentBillingPolicy`, single parameter `environmentId`, returns `BillingPolicyResponseModel`.

**This is a strictly better fit than the reverse-lookup approach suggested in earlier research** (enumerate all billing policies, then list each policy's environments). One call, direct, environment-scoped.

Reverse-lookup fallbacks, if ever needed:

```http
GET https://api.powerplatform.com/licensing/billingPolicies/{billingPolicyId}/environments?api-version=2024-10-01
  → { "value": [ { "billingPolicyId": "...", "environmentId": "..." } ], "@odata.nextLink": "..." }

GET https://api.powerplatform.com/licensing/billingPolicies/{billingPolicyId}/environments/{environmentId}?api-version=2024-10-01
  → { "billingPolicyId": "...", "environmentId": "..." }        # 404 when not linked
```

Neither supports filtering by environment ID; both are keyed by policy ID. Docs: [List Billing Policy Environments](https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy-environment/list-billing-policy-environments), [Get Billing Policy Environment](https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy-environment/get-billing-policy-environment).

> The provider today already calls `POST /licensing/billingPolicies/{id}/environments/add` and `/remove` at `api-version=2022-03-01-preview` (internal/services/licensing/api_licensing.go lines 173, 200). The GA `api-version` for the whole `licensing` billing-policy group is now `2024-10-01`.

### 4.2 Environment Management Settings — does NOT carry items 9-12

```http
GET   https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/settings?api-version=2024-10-01
GET   .../settings?$top={$top}&$select={$select}&api-version=2024-10-01     # $top defaults to 500
PATCH .../settings?api-version=2024-10-01
```

Envelope `GetEnvironmentManagementSettingResponse`: `objectResult` (`EnvironmentManagementSetting[]`), `nextLink`, `responseMessage`, `errors` (`EnvironmentServiceErrorResponse` `{code, message, details[]}`).

Complete `EnvironmentManagementSetting` property list, re-enumerated from the connector reference (33 properties):

| Group | Properties |
| --- | --- |
| Identity | `id`, `tenantId` |
| Storage SAS | `enableIpBasedStorageAccessSignatureRule`, `allowedIpRangeForStorageAccessSignatures`, `ipBasedStorageAccessSignatureMode` (int32), `loggingEnabledForIpBasedStorageAccessSignature` |
| Power Pages | `powerPages_AllowMakerCopilotsForNewSites`, `powerPages_AllowMakerCopilotsForExistingSites`, `powerPages_AllowProDevCopilotsForSites`, `powerPages_AllowSiteCopilotForSites`, `powerPages_AllowSearchSummaryCopilotForSites`, `powerPages_AllowListSummaryCopilotForSites`, `powerPages_AllowIntelligentFormsCopilotForSites`, `powerPages_AllowSummarizationAPICopilotForSites`, `powerPages_AllowNonProdPublicSites`, `powerPages_AllowNonProdPublicSites_Exemptions`, `powerPages_AllowProDevCopilotsForEnvironment` |
| Power Apps | `powerApps_ChartVisualization`, `powerApps_FormPredictSmartPaste`, `powerApps_FormPredictAutomatic`, `powerApps_CopilotChat`, `powerApps_NLSearch`, `powerApps_EnableFormInsights`, `powerApps_AllowCodeApps` |
| Copilot Studio | `copilotStudio_ConnectedAgents`, `copilotStudio_CodeInterpreter`, `copilotStudio_ConversationAuditLoggingEnabled`, `copilotStudio_ComputerUseAppAllowlist`, `copilotStudio_ComputerUseWebAllowlist`, `copilotStudio_ComputerUseSharedMachines`, `copilotStudio_ComputerUseCredentialsAllowed` |
| D365 Customer Service | `d365CustomerService_Copilot`, `d365CustomerService_AIAgents` |

**No Bing/chat toggle. No M365-services toggle. No cross-geo or cross-boundary Copilot data-movement policy. Items 9, 10, 11, 12 are confirmed absent here.**

* Docs: [List Environment Management Settings](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/list-environment-management-settings), [Update Environment Management Settings](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/update-environment-management-settings)
* Connector `EnvironmentManagementSetting` definition: https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/
* Tutorial: https://learn.microsoft.com/en-us/power-platform/admin/programmability-tutorial-environmentmanagement-settings

### 4.3 `governance` namespace — nothing for item 8

Full operation inventory of the `governance` namespace, from the what's-new changelog:

| Operation group | Operations |
| --- | --- |
| Rule Based Policies | Create, List, Get By ID, Update By ID, Patch, Remove Rule From, Create Environment Rule Based Assignment, Create Environment Group Rule Based Assignment, List Rule Assignments, List Rule Assignments By Environment Id / Environment Group Id / Policy Id |
| Rule Sets | Get Rule Set, Create Rule Set, List Rule Sets For Tenant, Update Rule Set, Delete Rule Set |
| Cross-Tenant Connection Reports | Create, Get, List |

Models: `Policy` `{id, tenantId, name, lastModified, ruleSets[], ruleSetCount}`, `RuleSet` `{id, version, inputs (free-form object)}`, `RuleSetDto` `{id, lastModified, environmentFilter, parameters[]}`, `RuleSetParameters` `{type, resourceType, value[]}`, `MgGovRule` `{id, value}`.

**No release/update cadence concept anywhere.** Managed-environment enable/disable lives in `environmentmanagement` (`POST .../governanceSetting/enableManaged` / `disableManaged`) and takes no cadence input; its effect surfaces as `EnvironmentResponse.protectionLevel` (`Basic` / `Standard`).

* Source: https://learn.microsoft.com/en-us/power-platform/admin/programmability-whats-new-changed
* Models: https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ → Definitions

The only cadence in the entire API is F&O-specific and semantically different:

```http
GET  https://api.powerplatform.com/dynamics/environments/{environmentId}/finOpsMaintenanceSettings?api-version=2024-10-01
PATCH ... (same path)
  → FinOpsAdminSettingsResponse {
      maintenanceWindowDaysOfWeek: FinOpsDayOfWeek[],
      maintenanceWindowCadence:    FinOpsUpdateCadence
    }
```

`FinOpsUpdateCadence` is documented as "Cadence for major version application updates" — i.e. Finance & Operations application versions, **not** the Dataverse `updateCadence` (`Frequent` / `Moderate`) that BAPI's `properties.updateCadence.id` carries. Doc: [Get Fin Ops Maintenance Settings](https://learn.microsoft.com/en-us/rest/api/power-platform/dynamics/finance-and-operations-maintenance-settings/get-fin-ops-maintenance-settings).

### 4.4 Operations history does NOT echo the create request (kills items 1, 2, 4, 5, 6)

```http
GET https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/operations?api-version=2024-10-01
GET https://api.powerplatform.com/environmentmanagement/environments/{targetEnvironmentId}/operations/{operationId}?api-version=2024-10-01
GET https://api.powerplatform.com/environmentmanagement/operations/{operationId}?api-version=2024-10-01
```

`OperationExecutionResult` — complete property list:

| JSON path | Type | Notes |
| --- | --- | --- |
| `name` | string | Operation name |
| `status` | `OperationStatus` enum | |
| `operationId` | string | |
| `startTime` / `endTime` | date-time | |
| `updatedEnvironment` | `Environment` | **Only `{environmentId, displayName, dataverseOrganizationUrl}`** — see §2.5 |
| `requestedBy` | `UserIdentity` `{userId, type, displayName, tenantId}` | The caller, not `usedBy` |
| `errorDetail` | `OperationErrorDetail` `{code, fieldErrors}` | |
| `stageStatuses[]` | `StageStatus` `{name, status, startTime, endTime, errorDetail}` | |

**There is no request-echo field.** The original `CreateEnvironmentRequest` (which does contain `description`, `usedBy`, `linkedEnvironmentMetadata.baseLanguageCode`, `.templates`, `.templateMetadata`) is not retained in any readable form. `updatedEnvironment` is the 3-property `Environment`, not `EnvironmentResponse`.

Paged wrapper: `OperationExecutionResultPagedCollection` `{collection[], continuationToken}`; query params `limit`, `continuationToken`.

* Docs: [Get Operations For Environment](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-operations-for-environment), [Get Operation By ID](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-operation-by-id)
* Model: https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ → Definitions → `OperationExecutionResult`, `Environment`

### 4.5 Templates for an EXISTING environment (item 5/6) — only a lossy proxy

There is **no** endpoint that returns "which provisioning templates were applied to this environment." The provisioning-time list is location-scoped and environment-agnostic:

```http
GET https://api.powerplatform.com/environmentmanagement/provisioning/locations/{location}/templates?api-version=2024-10-01
  → EnvironmentTemplateResourceCollection { collection: EnvironmentTemplate[] }
  EnvironmentTemplate { name, displayName, isCustomerEngagement, isSupportedForResetOperation,
                        availability: TemplateAvailability[] { environmentSku, isDisabled, disabledReason{code,message} } }
```

The nearest per-environment read is the app-package inventory, which the provider already wraps as `powerplatform_environment_application_packages`:

```http
GET https://api.powerplatform.com/appmanagement/environments/{environmentId}/applicationPackages?api-version=2024-10-01
  → ApplicationPackageContinuationResponse { value: ApplicationPackage[], @odata.nextLink }
  ApplicationPackage { id, uniqueName, version, localizedName, localizedDescription, applicationId,
                       applicationName, applicationDescription, publisherName, publisherId,
                       state: InstancePackageState, catalogVisibility, applicationVisibility,
                       lastError, startDateUtc, endDateUtc, supportedCountries, ... }
```

`InstancePackageState` includes `Installed` and `TemplateInstalled` — the latter is the signal that a package landed via a provisioning template. But the identifier is a **package unique name** (`msdyn_SalesPatch`), not a BAPI template code (`D365_Sales`). Reconstructing `dataverse.templates` from this would require a hand-maintained mapping table and would not round-trip. **Not recommended.**

`template_metadata` (item 6) has no proxy at all — it is an opaque provisioning-time JSON blob (`CreateEnvironmentRequestLinkedMetadata_templateMetadata`, typed only as `object`).

* Docs: [Get Provisioning Templates](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-templates), [Get Environment Application Package](https://learn.microsoft.com/en-us/rest/api/power-platform/appmanagement/applications/get-environment-application-package)
* `TemplateAvailability` / `InstancePackageState` definitions: https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/

### 4.6 Licensing entitlements — richer environment metadata, but not the missing 13

```http
GET https://api.powerplatform.com/licensing/environments/{environmentId}/entitlements?$filter={$filter}&api-version=2024-10-01
  → EnvironmentEntitlementResponseModel[]
```

`EnvironmentEntitlementResponseModel` carries environment metadata not present on `EnvironmentResponse`: `isManagedEnvironment` (boolean), `location`, `scenario` (`EnvironmentScenario`), `disasterRecoveryState`, `disasterRecoveryLocation`, `addons[]` (`EnvironmentAddonResponseModel {addonType, allocated, addonUnit}`), `permissions[]` (`EnvironmentPermissionResponseModel {name, displayName}`), `cleanupOpportunitySize`, `recommendationCount`, plus `entitlementId`, `productCategories[]`, `environmentType`, `environmentName`.

**Worth flagging for the wider migration** — `addons[]` and `permissions[]` are the new-API replacements for BAPI's `$expand=permissions,properties.capacity`. But **none of the 13 target attributes are here.**

* Doc: [Get Many Environment Entitlements](https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/entitlement/get-many-environment-entitlements)
* Model: https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ → Definitions → `EnvironmentEntitlementResponseModel`

### 4.7 Azure Resource Graph passthrough — the only remaining lead for item 1

```http
POST https://api.powerplatform.com/resourcequery/resources/query?api-version=2024-10-01
Content-Type: application/json

{
  "TableName": "PowerPlatformResources",
  "Clauses": [ { "$type": "WhereClause", ... }, { "$type": "ProjectClause", ... } ],
  "Options": { "Top": 100, "Skip": 0, "SkipToken": null }
}
```

Response `ResourceQueryResponse { totalRecords, count, resultTruncated (0=truncated, 1=not), skipToken, data: ResourceItem[] }`.

`ResourceItem` is explicitly documented as "ARG row with Power Platform-specific fields. **Arbitrary properties may exist under `properties`**" — a genuinely free-form bag. Typed siblings: `id`, `name`, `type`, `tenantId`, `kind`, `location`, `resourceGroup`, `subscriptionId`, `managedBy`, `sku`, `plan`, `tags`, `identity`, `zones`, `extendedLocation`, `environmentId`, `environmentId1`, `environmentName`, `environmentRegion`, `environmentType`, `isManagedEnvironment`.

**This is the only surface in the entire Power Platform API where an undocumented `description` (or `updateCadence`, or `copilotPolicies`) could plausibly surface**, because `properties` is deliberately unmodelled. Whether `PowerPlatformResources` rows actually carry the BAPI environment property bag is **UNVERIFIED — requires live capture.** Even if they do, the ARG projection is a KQL query surface with its own throttling (`429 – ARG throttling` is a documented response), a different consistency model, and preview-grade stability. **Not recommended as a Terraform read path**, but worth one exploratory query to close the question.

* Doc: https://learn.microsoft.com/en-us/rest/api/power-platform/resourcequery/resource-query/query-resources
* Clause types: `WhereClause`, `ProjectClause`, `ExtendClause`, `SummarizeClause`, `JoinClause`, `OrderByClause`, `TakeClause`, `CountClause`, `DistinctClause` (see `Clause` composed-type wrapper in the Models namespace).

---

## 5. Additional findings relevant to the migration

### 5.1 Enterprise policy shape change (items 13 & 14 in detail)

| | BAPI (`properties.enterprisePolicies.*`) | Power Platform API (`enterprisePolicies.*`) |
| --- | --- | --- |
| Keys | `vnets`, `customerManagedKeys`, `identity` | `networkInjection`, `encryption`, `identity`, **`privateEndpoint`** (new) |
| Per-policy fields | `policyId`, `location`, `id`, `systemId`, `linkStatus` | `id`, `resourceId`, `status`, **`error`** (new) |
| Status enum | (undocumented on BAPI) | `Linking`, `Unlinking`, `Linked`, `Failed`, `LinkingOnline`, `UnlinkingOnline` |

Provider's current DTO for reference: internal/services/environment/dto.go lines 94-106 (`EnvironmentEnterprisePoliciesDto` / `EnterprisePolicyDto`); model mapping at internal/services/environment/models.go lines 379-451; rendered schema at docs/resources/environment.md lines 126-134 (`id`, `location`, `status`, `system_id`, `type`).

Likely field correspondence — **UNVERIFIED, requires side-by-side live capture of the same environment through both APIs**:

| BAPI | New API | Rationale |
| --- | --- | --- |
| `vnets` | `networkInjection` | Both are the VNet-injection policy link |
| `customerManagedKeys` | `encryption` | Both are the CMK policy link |
| `identity` | `identity` | Name-identical |
| — | `privateEndpoint` | No BAPI counterpart |
| `linkStatus` | `status` | Both "link status" |
| `id` | `resourceId` | New API doc: "The fully-qualified Azure resource ID of the enterprise policy" — this is the ARM `/subscriptions/{s}/resourceGroups/{rg}/providers/Microsoft.PowerPlatform/enterprisePolicies/{name}` form |
| `systemId` **or** `policyId` | `id` | New API doc: "The ID of the enterprise policy" — bare ID, no "fully-qualified Azure resource ID" qualifier |
| `location` | **dropped** | No counterpart |
| — | `error` | New: "Error details when the link status is Failed" |

**Answer to question 14: `resourceId` is almost certainly NOT BAPI's `systemId`. It is the full ARM resource ID, which BAPI exposes as `id`.** The naive rename `systemId → resourceId` would produce a state diff on every existing `powerplatform_environment` with a linked enterprise policy. Treat this as a breaking-shape change until a live capture proves otherwise.

For item 13 (`location`): the Azure region is embedded in the ARM resource ID's resource group / can be read via `GET https://management.azure.com{resourceId}?api-version=2020-10-30-preview` → `location`. That is a cross-API call with a different token audience (`https://management.azure.com/.default`) — a significant new dependency for one cosmetic attribute.

### 5.2 `usedBy` vs `createdFor` (item 2 in detail)

* BAPI `properties.usedBy` = `{id, tenantId, type}` — for Developer SKU environments this is the user the environment belongs to. Confirmed in this repo: internal/services/environment/tests/resource/Validate_Create_Dev_Env/get_environment_00000000-0000-0000-0000-000000000001.json lines 19-23, alongside `"environmentSku": "Developer"`.
* New API request model `CreateEnvironmentRequest.usedBy` = `UserIdentity {userId, type, displayName, tenantId}` — accepted on **input** to Provision New Environment. So the concept survives on the write path.
* New API response model `EnvironmentResponse.createdFor` = `EnvironmentPrincipal {id, type}` — "Represents a principal (user or application)."

The pairing `createdBy` (who called the API) + `createdFor` (who it was created for) maps naturally onto BAPI's `createdBy` + `usedBy`. **Hypothesis: `createdFor.id` is populated from the create request's `usedBy.userId`.** Not stated anywhere. **UNVERIFIED — requires live capture**: provision a Developer environment with an explicit `usedBy`, then read it back and compare `createdFor.id`.

If confirmed, `owner_id` maps to `createdFor.id` and item 2 moves from "unavailable" to "renamed", which materially reduces the breaking-change surface.

### 5.3 `description` is accepted on write but never read back

Three write operations accept a `description`:

| Operation | Field |
| --- | --- |
| Provision New Environment (`POST /environmentmanagement/provisioning/environments`) | `description` — "An optional description for the environment." |
| Reset Environment (`POST /environmentmanagement/environments/{id}/reset`) | `description` — "An optional description for the environment to reset to." |
| Create/Update Environment Group | `description` (different resource) |

No read operation returns it. This is a genuine write-only attribute on the new API. For Terraform this means: either keep the BAPI read for `description`, drop the attribute, or mark it as a non-refreshable input (which breaks drift detection and import).

Sources: [Provision New Environment](https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/provision-new-environment), connector "Resets the environment" parameter table.

### 5.4 There is no Update Environment operation

Worth restating because it constrains every "can we PATCH it back" workaround: the `environmentmanagement` namespace has **no** generic environment update. The only mutations are `ModifyEnvironmentSku`, `EnableManaged`/`DisableManaged`, `Enable`/`Disable` (admin mode), `Reset`, `Copy`, `Restore`, `Recover`, `Delete`, `LinkDataverse`, plus the `/settings` PATCH. Attributes like `description`, `securityGroupId`, and `displayName` can only be changed via `Reset` (destructive) — or by staying on BAPI.

---

## 6. Explicitly UNVERIFIED — requires live capture

| # | Question | Cheapest experiment |
| --- | --- | --- |
| U1 | Does `GET /environmentmanagement/environments/{id}` return any property beyond the documented 26? | One authenticated GET against a real environment; dump the raw JSON body. |
| U2 | Does `$select=description` (or `usedBy`, `updateCadence`, `bingChatEnabled`) return `400`, silently drop the token, or return a value? | Four GETs with different `$select` values; compare status codes and bodies. |
| U3 | Is `createdFor.id` populated from the create request's `usedBy.userId`? | Provision a Developer environment with explicit `usedBy`, read back, compare. |
| U4 | Is `enterprisePolicies.*.resourceId` the ARM resource ID (== BAPI `id`) or the `systemId`? | Read the same CMK/VNet-linked environment through BAPI and the new API; diff the four fields. |
| U5 | Does `PowerPlatformResources` in `resourcequery` expose the BAPI environment property bag (incl. `description`, `updateCadence`, `copilotPolicies`)? | One `POST /resourcequery/resources/query` with `TableName: "PowerPlatformResources"` and no `ProjectClause`; inspect `data[].properties`. |
| U6 | What does `GET /licensing/environments/{id}/billingPolicy` return when no policy is linked? (`404`? empty body? `null`?) | GET against an environment with no billing policy. |
| U7 | Does Dataverse `organizations` return `uniquename` in the provider's current unfiltered GET? (Not present in the trimmed fixture at internal/services/environment_settings/tests/datasource/organisations.json.) | Inspect a live `GET {envUrl}/api/data/v9.0/organizations` response, or add `?$select=uniquename,languagecode`. |

Capture tooling for all of the above is already documented in this repo: devdocs/adr/mitmproxy.md (recording) and devdocs/adr/httpmocks.md (turning captures into fixtures).

---

## 7. Source links (complete)

**Learn REST reference**

* Get Environment By Id For User — https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user
* List Environments For User — https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/list-environments-for-user
* Provision New Environment — https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/provision-new-environment
* Get Provisioning Templates — https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-templates
* List Environment Management Settings — https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/list-environment-management-settings
* Get Operations For Environment — https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-operations-for-environment
* Get Environment Billing Policy — https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/environment-billing-policy/get-environment-billing-policy
* Environment Billing Policy (group) — https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/environment-billing-policy
* Billing Policy Environment (group) — https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy-environment
* Get Billing Policy Environment — https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy-environment/get-billing-policy-environment
* List Billing Policy Environments — https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/billing-policy-environment/list-billing-policy-environments
* Get Many Environment Entitlements — https://learn.microsoft.com/en-us/rest/api/power-platform/licensing/entitlement/get-many-environment-entitlements
* Get Fin Ops Maintenance Settings — https://learn.microsoft.com/en-us/rest/api/power-platform/dynamics/finance-and-operations-maintenance-settings/get-fin-ops-maintenance-settings
* Query Resources (ARG) — https://learn.microsoft.com/en-us/rest/api/power-platform/resourcequery/resource-query/query-resources
* Get Environment Application Package — https://learn.microsoft.com/en-us/rest/api/power-platform/appmanagement/applications/get-environment-application-package
* API reference root — https://learn.microsoft.com/en-us/rest/api/power-platform/

**.NET SDK (`Microsoft.PowerPlatform.Management`, code-generated via Kiota)**

* `EnvironmentResponse` class — https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.models.environmentresponse?view=power-platform-latest
* Models namespace (full class/enum inventory) — https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.models?view=power-platform-latest
* `EnvironmentsRequestBuilderGetQueryParameters` — https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.environmentmanagement.environments.environmentsrequestbuilder.environmentsrequestbuildergetqueryparameters?view=power-platform-latest
* `BillingPolicyRequestBuilder` (`\licensing\environments{environmentId}\billingPolicy`) — https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.licensing.environments.item.billingpolicy.billingpolicyrequestbuilder?view=power-platform-latest
* NuGet package — https://www.nuget.org/packages/Microsoft.PowerPlatform.Management
* Python SDK — https://pypi.org/project/powerplatform-management/

**Connector reference (third independent rendering of the same spec)**

* Power Platform for Admins V2 — https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ (updated 2026-08-07; git commit `a270b1b3363b7bb5c91cd4a5bc28ab0a2736c530`)

**Changelog / versioning**

* Programmability what's new or changed — https://learn.microsoft.com/en-us/power-platform/admin/programmability-whats-new-changed
* Versioning and support — https://learn.microsoft.com/en-us/power-platform/admin/programmability-versioning-support
* Permission reference — https://learn.microsoft.com/en-us/power-platform/admin/programmability-permission-reference

**Third-party corroboration**

* `microsoft/Employee-Self-Service-Agent-Developer-Kit` — `tests/mocks/powerplatform.py` (builds `EnvironmentResponse` mocks strictly from the Learn field tables; declares `MOCK_STATUS = "documented"`), `tests/fixtures/cassettes/INDEX.md` (API tier registry: "No public OpenAPI and no no-auth `$metadata`, so `validatable` does not apply") — https://github.com/microsoft/Employee-Self-Service-Agent-Developer-Kit

**Dataverse (for items 4 and 7)**

* `organization` table reference — https://learn.microsoft.com/en-us/power-apps/developer/data-platform/reference/entities/organization

**Workspace files consulted**

* internal/services/environment/dto.go (BAPI `EnvironmentPropertiesDto`, `EnterprisePolicyDto`)
* internal/services/environment/models.go (`convertEnterprisePolicyModelFromDto`)
* internal/services/environment/api_environment.go (line 253 — BAPI `$expand` string)
* internal/services/licensing/api_licensing.go (existing billing-policy calls)
* internal/services/environment_settings/api_environment_settings.go (Dataverse `organizations` calls)
* internal/services/environment_settings/tests/datasource/organisations.json
* internal/services/environment/tests/resource/Validate_Create_Dev_Env/get_environment_00000000-0000-0000-0000-000000000001.json
* docs/resources/environment.md (rendered `enterprise_policies` schema)
* devdocs/adr/mitmproxy.md, devdocs/adr/httpmocks.md

---

## 8. Recommended next research

- [ ] Run the seven UNVERIFIED experiments in §6 against a real tenant with mitmproxy (devdocs/adr/mitmproxy.md), starting with U1 and U2 — those two alone would move most "No" verdicts from "documented absence" to "observed absence".
- [ ] Capture the same environment through BAPI and the Power Platform API side by side to settle U4 (enterprise policy field correspondence) and U3 (`usedBy` → `createdFor`).
- [ ] Confirm U6 (billing-policy-not-linked behaviour) before wiring `GET /licensing/environments/{id}/billingPolicy` into the environment read path — the provider needs to distinguish "no policy" from "permission denied".
- [ ] Decide the policy for the 8 genuinely-unrecoverable attributes: hybrid read (keep BAPI for those fields), deprecate-and-remove, or block the migration until Microsoft adds them. This is a product decision, not a research one.
- [ ] Separately assess whether `EnvironmentEntitlementResponseModel.addons[]` / `.permissions[]` (§4.6) can replace BAPI's `$expand=permissions,properties.capacity` — out of scope here but adjacent and likely needed by the same migration.
- [ ] File a Microsoft feedback item requesting `description`, `usedBy`, `updateCadence`, and the Copilot policy flags on `EnvironmentResponse` — the API is explicitly described as reaching "full parity ... for what an administrator can perform in Power Platform admin center", so these are gaps rather than intentional removals.

## 9. Clarifying questions

1. **Is a live tenant available to this workstream?** Seven of the open questions collapse to a handful of authenticated GETs. Without one, items 1, 2, and 14 stay at Medium/Low confidence indefinitely.
2. **Is a hybrid BAPI + Power Platform API read acceptable for one release?** Keeping the BAPI GET solely to populate `description`, `owner_id`, `cadence`, and the four AI/Copilot toggles would make the migration non-breaking, at the cost of two auth audiences and two round-trips per read.
3. **Is `cadence` (`Frequent` / `Moderate`) still settable at all on the new provisioning path?** It is absent from `CreateEnvironmentRequest` as well as from `EnvironmentResponse` — so it may be write-impossible, not just read-impossible, which changes it from a state-refresh problem to a functionality-removal problem.
4. **How should `enterprise_policies` be versioned?** The key set changes (`vnets`→`networkInjection`, `customerManagedKeys`→`encryption`, plus new `privateEndpoint`), `location` disappears, and `system_id`/`resourceId` may not be the same value. This looks like a major-version schema change for that nested block regardless of what U4 finds.
