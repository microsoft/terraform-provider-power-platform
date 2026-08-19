<!-- markdownlint-disable-file -->
# Subagent Research: Power Platform API (api.powerplatform.com) Environment Management

Research topic: exhaustively document the public Power Platform API Environment Management surface (especially environment provisioning) to support migrating the Terraform provider away from legacy BAPI (`api.bap.microsoft.com`).

Status: COMPLETE for the documented public surface. Several items flagged UNVERIFIED (see "Gaps and UNVERIFIED items").

Research date: 2026-08-18. Docs snapshot: Microsoft Learn `powerplatform-rest` repo, git commit `92813e8ead95f02f8f28e750c0f54f4cbb5eb00f` (updated 2026-08-06) for most `environmentmanagement` pages; `9d5f2974d5b8a7c5df33a967224d25e0d0868f5d` (2026-05-08) for Environment Groups.

## Research Questions

1. What Environment Management endpoints exist on `api.powerplatform.com`? Full URL templates, methods, `api-version` values. — ANSWERED
2. Full request/response schemas for environment provisioning (create), link Dataverse, delete, recover, reset, copy. — ANSWERED
3. Full `Environment` model returned by GET/LIST environments on the new API (every JSON path). — ANSWERED
4. Locations / Currencies / Languages / Templates list endpoints and their full response schemas. — ANSWERED
5. Async operation behavior (202, polling, terminal states). — PARTIALLY ANSWERED (polling endpoints documented; response headers NOT documented — flagged)
6. Auth scope and required permissions/roles. — ANSWERED (scope); permissions PARTIAL
7. Comparison to legacy BAPI: what is present in one but not the other. — ANSWERED

---

## 1. Fundamentals

### 1.1 Request shape

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/

```http
{HTTP method} https://api.powerplatform.com/{namespace}/{resource}?api-version={version}
```

* `{namespace}` — logical grouping (`environmentmanagement`, `licensing`, `appmanagement`, `authorization`, `governance`, `connectivity`, `resourcequery`, `powerapps`, `powerautomate`, `powerpages`, `copilotstudio`, `dynamics`, `workflowsagent`).
* `{resource}` — resource path. Tenant-level resources infer `tenantId` from the OAuth bearer token and route to the region matching the tenant's physical address. Environment-level resources require an `environmentId` in the path.
* HTTP verbs: `GET` (read), `POST` (create or action), `PATCH` (update), `PUT` (replace), `DELETE` (remove). `GET`/`DELETE` take no body.

### 1.2 Authentication scope

Source: every `environmentmanagement` operation page, "Security" section.

* Flow: Microsoft Entra ID OAuth2, implicit.
* Authorization URL: `https://login.microsoftonline.com/common/oauth2/authorize?resource=https://api.powerplatform.com`
* Scope: `.default`
* Effective client-credentials scope: `https://api.powerplatform.com/.default`

Contrast with legacy BAPI, which uses `https://service.powerapps.com/.default` (or the `https://api.bap.microsoft.com/` resource).

The provider already has both constants defined:

* internal/constants/constants.go — `PUBLIC_POWERPLATFORM_API_DOMAIN = "api.powerplatform.com"`, `PUBLIC_POWERPLATFORM_API_SCOPE = "https://api.powerplatform.com/.default"`
* Verified in internal/constants/constants_test.go.

So no new auth plumbing is required — the provider already authenticates against `api.powerplatform.com` for the `licensing`, `appmanagement`, and `authorization` namespaces.

### 1.3 API versions

Source: https://learn.microsoft.com/en-us/power-platform/admin/programmability-versioning-support

| Version | Type | Notes |
| --- | --- | --- |
| `2024-10-01` | **Stable** | Current stable version. This is what ALL `environmentmanagement` operations document. |
| `2020-10-01` | Generally available | "specific to environment management and is also commonly referred to as Business Application Platform (BAP) API. The functionality of this set of endpoints are made available in the newer versions of Power Platform API." |

Other `api-version` values observed in the wild / in this repo (namespace-specific, NOT `environmentmanagement`):

* `2022-03-01-preview` — `licensing/*`, `appmanagement/*` (used today by the provider).
* `2021-10-01-preview` — `governance/environmentGroups/{id}/ruleSets` (tenant-scoped host `*.tenant.api.powerplatform.com`).
* `2024-10-01` — `authorization/roleAssignments`, `authorization/roleDefinitions` (used today by the provider).
* `1` — environment-scoped `connectivity/*` and `appmanagement/environments/{id}/operations/{opId}` (used today by the provider).

Deprecation policy: incrementing the API version immediately deprecates the prior GA version; retirement occurs 12 months after announcement. Breaking changes (removed/renamed property, changed type, changed URL, new required request parameter) force a version bump; additive/nullable properties, new enum members, added paging, changed error codes, reordered properties do NOT.

**Conclusion for migration: target `api-version=2024-10-01` for everything under `environmentmanagement`.**

### 1.4 Permissions (Entra granular permissions)

Source: https://learn.microsoft.com/en-us/power-platform/admin/programmability-permission-reference

Naming convention: `{namespace}.{resourceType}.{action}` where action is derived from the HTTP method:

| HTTP Method | Path Structure | Action Name |
| --- | --- | --- |
| GET or HEAD | Any | `Read` |
| DELETE | Any | `Delete` |
| PATCH | Any | `Update` |
| PUT | Any | `Create` and `Update` |
| POST | `/{namespace}/.../{resourceType}` | `Create` |
| POST | `/{namespace}/.../{resourceType}/{resourceId}/{action}` | `{action}` |

Documented `EnvironmentManagement.*` permissions:

| Permission | Display name | Description |
| --- | --- | --- |
| `EnvironmentManagement.Environments.Read` | Read environments | Allows reading of environments. |
| `EnvironmentManagement.Groups.Read` | Read environment groups | Allows reading of environment groups. |
| `EnvironmentManagement.Groups.ReadWrite` | Read and write environment groups | Allows reading and writing of environment groups. |
| `EnvironmentManagement.Settings.Read` | Read environment management settings | Allows reading of Environment Management Settings. |
| `EnvironmentManagement.Settings.ReadWrite` | Update environment management settings | Allows update of environment management settings. |

> **UNVERIFIED — needs confirmation**: there is NO documented granular permission for **creating/deleting/recovering/resetting/copying environments** (no `EnvironmentManagement.Environments.Create` / `.Delete` / `.ReadWrite` in the published permission reference as of its 2025-09-26 revision). The write-side provisioning operations most likely require a **tenant admin role** (Power Platform Administrator / Global Administrator / Dynamics 365 Administrator) rather than a granular app permission, matching how BAPI `scopes/admin/environments` works. This must be confirmed before switching the provider's create/delete path, especially for service-principal auth.

### 1.5 Related SDKs (useful for reverse-engineering exact payloads)

* .NET: `Microsoft.PowerPlatform.Management` NuGet — namespace `Microsoft.PowerPlatform.Management.Environmentmanagement` exposes `EnvironmentmanagementRequestBuilder` ("Builds and executes requests for operations under \environmentmanagement"). Latest observed `2.0.3503.299` (July 2026).
  * https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.environmentmanagement?view=power-platform-latest
* Python: `powerplatform-management` on PyPI (initial release Sept 2025, latest observed `2.0.3503.299`).
* Connector: "Power Platform for Admins V2" — https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/ — includes an **Environment Management MCP Server** (Sept 2025).
* `Microsoft.PowerApps.CLI` (pac cli).

These generated SDKs are the best source for exact wire payloads/headers where the Learn reference omits examples.

---

## 2. Complete `environmentmanagement` operation inventory

Derived from https://learn.microsoft.com/en-us/power-platform/admin/programmability-whats-new-changed plus each operation-group page. All use `api-version=2024-10-01`. Host is `https://api.powerplatform.com` unless noted.

| Operation group | Operation | Method | Path |
| --- | --- | --- | --- |
| Environment Provisioning | Provision New Environment | POST | `/environmentmanagement/provisioning/environments` |
| Environment Provisioning | Link Dataverse | PATCH | `/environmentmanagement/provisioning/environments/{environmentId}/link` |
| Environment Provisioning | Get Supported Locations | GET | `/environmentmanagement/provisioning/locations` |
| Environment Provisioning | Get Provisioning Currencies | GET | `/environmentmanagement/provisioning/locations/{location}/currencies` |
| Environment Provisioning | Get Provisioning Languages | GET | `/environmentmanagement/provisioning/locations/{location}/languages` |
| Environment Provisioning | Get Provisioning Templates | GET | `/environmentmanagement/provisioning/locations/{location}/templates` |
| Environments | List Environments For User | GET | `/environmentmanagement/environments` |
| Environments | Get Environment By Id For User | GET | `/environmentmanagement/environments/{environmentId}` |
| Environment Delete | Delete Environment By ID | DELETE | `/environmentmanagement/environments/{environmentId}` |
| Environment Recover | Recover Environment | POST | `/environmentmanagement/environments/{environmentId}/recover` |
| Environment Reset | Reset Environment | POST | `/environmentmanagement/environments/{environmentId}/reset` |
| Environment Copy | Copy Environment | POST | `/environmentmanagement/environments/{targetEnvironmentId}/copy` |
| Environment Copy | Get Environment Copy Candidates | GET | *(slug documented; path UNVERIFIED)* |
| Environment Restore | Restore Environment | POST | *(slug documented; path UNVERIFIED)* |
| Environment Restore | Get Restore Candidates | GET | *(slug documented; path UNVERIFIED)* |
| Environment Backup | Create / Get / Delete Environment Backup | POST/GET/DELETE | *(slugs documented; paths UNVERIFIED)* |
| Environment State | Enable Environment / Disable Environment | POST | *(slugs listed in what's-new but doc pages 404 as of this research; paths UNVERIFIED)* |
| Environment Managed Governance | Enable Managed Environment | POST | `/environmentmanagement/environments/{environmentId}/governanceSetting/enableManaged` |
| Environment Managed Governance | Disable Managed Environment | POST | `/environmentmanagement/environments/{environmentId}/governanceSetting/disableManaged` |
| Operation | Get Operation By ID | GET | `/environmentmanagement/operations/{operationId}` |
| Operation | Get Environment Operation By ID | GET | `/environmentmanagement/environments/{targetEnvironmentId}/operations/{operationId}` |
| Operation | Get Operations For Environment | GET | `/environmentmanagement/environments/{environmentId}/operations` |
| Environment Management Settings | List Environment Management Settings | GET | `/environmentmanagement/environments/{environmentId}/settings` |
| Environment Management Settings | Update Environment Management Settings | PATCH | `/environmentmanagement/environments/{environmentId}/settings` |
| Environment Management Settings | Create Environment Management Settings | POST/PUT | *(slug documented; path UNVERIFIED)* |
| Environment Groups | Create Environment Group | POST | `/environmentmanagement/environmentGroups` *(UNVERIFIED — inferred)* |
| Environment Groups | Get Environment Group | GET | `/environmentmanagement/environmentGroups/{groupId}` |
| Environment Groups | List Environment Groups | GET | `/environmentmanagement/environmentGroups` *(UNVERIFIED — inferred)* |
| Environment Groups | Update Environment Group | PATCH/PUT | `/environmentmanagement/environmentGroups/{groupId}` *(UNVERIFIED — inferred)* |
| Environment Groups | Delete Environment Group | DELETE | `/environmentmanagement/environmentGroups/{groupId}` *(UNVERIFIED — inferred)* |
| Environment Groups | Add Environment To Group | POST | `/environmentmanagement/environmentGroups/{groupId}/addEnvironment/{environmentId}` |
| Environment Groups | Remove Environment From Group | POST | `/environmentmanagement/environmentGroups/{groupId}/removeEnvironment/{environmentId}` *(UNVERIFIED — inferred)* |
| Environment Groups | Get Environment Group Operation | GET | `/environmentmanagement/environmentGroupOperations/{operationId}` |
| Failover | Enable/Disable Disaster Recovery, Perform Failback, Perform Force Failover, Perform DR Drill, Get Business Continuity State Full Snapshot | POST/GET | *(slugs documented; paths UNVERIFIED)* |

Release timeline (from the what's-new page):

* **Jan 2025** — `List Environments For User`, `Get Environment By Id For User` (preview).
* **Feb 2025** — `Delete Environment By ID`, `Get Operation By ID`, `List Operations For Environment`.
* **Jul 2025** — `Enable/Disable Managed Environment` (preview).
* **Oct 2025** — `Recover Environment`, `Restore Environment`, `Get Restore Candidates`, `Copy Environment`, `Get Environment Copy Candidates`, `Enable/Disable Environment`, environment backup ops, failover/DR ops.
* **Apr 2026** — **`Provision New Environment`, `Link Dataverse`, `Get Supported Locations`, `Get Provisioning Currencies`, `Get Provisioning Languages`, `Get Provisioning Templates`** ← the provisioning surface is brand new.
* **May 2026** — `Get Environment Operation By ID` (preview).
* **Jul 2026** — `Reset Environment`; bug fixes to `Provision New Environment`.

> **Migration implication**: the provisioning endpoints are roughly four months old at the time of this research and still receiving "various bug fixes" (July 2026). Expect churn.

---

## 3. Environment Provisioning — full detail

Operation group: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning

### 3.1 Provision New Environment

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/provision-new-environment

```http
POST https://api.powerplatform.com/environmentmanagement/provisioning/environments?api-version=2024-10-01
```

**URI parameters**

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `api-version` | query | yes | string | The API version. |

**Request body** — media types `application/json`, `text/json`, `application/*+json`. Model: `CreateEnvironmentRequest`.

| JSON path | Required | Type | Description |
| --- | --- | --- | --- |
| `displayName` | **yes** | string (minLength 1) | The display name of the environment. |
| `environmentSku` | **yes** | `EnvironmentSku` enum | The environment SKU. |
| `billingPolicy` | no | `CreateEnvironmentRequestBillingPolicy` | Billing policy for the environment. |
| `billingPolicy.id` | no | string | The billing policy ID. |
| `cluster` | no | `CreateEnvironmentRequestCluster` | Cluster configuration. |
| `cluster.category` | no | string | The cluster category. E.g. `FirstRelease`. |
| `connectedGroupIdForTeamsEnvironment` | no | string | Microsoft 365 Group ID linked to a Teams environment during provisioning. Not applicable to non-Teams environments. |
| `databaseType` | no | string | The type of database to create (for example, `CommonDataService`). |
| `description` | no | string | An optional description for the environment. |
| `finOpsMetadata` | no | `CreateEnvironmentRequestFinOpsMetadata` | FinOps metadata for environment provisioning. |
| `finOpsMetadata.id` | no | string | The FinOps environment ID. |
| `finOpsMetadata.type` | no | string | The FinOps environment link type. |
| `finOpsMetadata.url` | no | string | The FinOps environment URL. |
| `governanceConfiguration` | no | `CreateEnvironmentRequestGovernance` | Governance configuration. |
| `governanceConfiguration.protectionLevel` | no | `ProtectionLevel` enum | The environment governance protection level. |
| `linkedEnvironmentMetadata` | no | `CreateEnvironmentRequestLinkedMetadata` | Metadata for the linked Dataverse environment. |
| `linkedEnvironmentMetadata.baseLanguageCode` | no | integer (int32) | The base language code (for example, `1033` for English). |
| `linkedEnvironmentMetadata.currency` | no | `EnvironmentRequestCurrency` | Currency settings for an environment. |
| `linkedEnvironmentMetadata.currency.code` | no | string | The currency code (for example, `USD`). |
| `linkedEnvironmentMetadata.currency.name` | no | string | The currency name. |
| `linkedEnvironmentMetadata.currency.precision` | no | integer (int32) | The currency precision. |
| `linkedEnvironmentMetadata.currency.symbol` | no | string | The currency symbol. |
| `linkedEnvironmentMetadata.domainName` | no | string | The domain name. |
| `linkedEnvironmentMetadata.securityGroupId` | no | string | The security group ID. |
| `linkedEnvironmentMetadata.templateMetadata` | no | object | A JSON object payload customized for the selected templates. |
| `linkedEnvironmentMetadata.templates` | no | string[] | The templates to apply. |
| `location` | no | string | The location where the environment will be provisioned. **Mutually exclusive with `macroRegion`.** |
| `macroRegion` | no | string | The macro region where the environment will be provisioned. |
| `parentEnvironmentGroup` | no | `CreateEnvironmentRequestParentGroup` | Parent environment group. |
| `parentEnvironmentGroup.id` | no | string | The environment group ID. |
| `usedBy` | no | `UserIdentity` | Represents the identity of a user. |
| `usedBy.displayName` | no | string | The display name of the user. |
| `usedBy.tenantId` | no | string | The tenant ID of the user. |
| `usedBy.type` | no | string | The type of the user identity (for example, `User`). |
| `usedBy.userId` | no | string | The ID of the user. |

**Dataverse vs no-Dataverse.** The docs do not spell it out explicitly, but the model implies:

* **With Dataverse** → send `databaseType: "CommonDataService"` **and** a populated `linkedEnvironmentMetadata` object (currency, baseLanguageCode, domainName, securityGroupId, templates).
* **Without Dataverse** → omit `databaseType` and omit `linkedEnvironmentMetadata`; then later call **Link Dataverse** (see 3.2) to add the database.

> **UNVERIFIED — needs confirmation**: the exact required combination for a no-Dataverse environment (is `databaseType` truly optional, or must it be explicitly absent vs. null?). Also unverified whether `location` accepts BAPI-style location names (`unitedstates`, `europe`) — highly likely given `Get Supported Locations` returns a `name` field, but not stated.

**Example request body (reconstructed from the schema — the Learn page publishes no example payload):**

```json
{
  "displayName": "Contoso Dev",
  "environmentSku": "Sandbox",
  "description": "Terraform managed environment",
  "location": "unitedstates",
  "databaseType": "CommonDataService",
  "cluster": {
    "category": "FirstRelease"
  },
  "billingPolicy": {
    "id": "00000000-0000-0000-0000-000000000001"
  },
  "governanceConfiguration": {
    "protectionLevel": "Standard"
  },
  "parentEnvironmentGroup": {
    "id": "00000000-0000-0000-0000-000000000002"
  },
  "linkedEnvironmentMetadata": {
    "baseLanguageCode": 1033,
    "domainName": "contoso-dev",
    "securityGroupId": "00000000-0000-0000-0000-000000000003",
    "currency": {
      "code": "USD",
      "name": "US Dollar",
      "symbol": "$",
      "precision": 2
    },
    "templates": ["D365_Sales"],
    "templateMetadata": {
      "PostProvisioningPackages": [
        {
          "applicationUniqueName": "msdyn_SalesPatch",
          "parameters": "DisableSalesTrial=true"
        }
      ]
    }
  },
  "usedBy": {
    "userId": "00000000-0000-0000-0000-000000000004",
    "type": "User",
    "tenantId": "00000000-0000-0000-0000-000000000005",
    "displayName": "Jane Doe"
  }
}
```

> **UNVERIFIED**: `templateMetadata` is typed only as `object` ("A JSON object payload customized for the selected templates"). The `PostProvisioningPackages` shape above is carried over from BAPI (internal/services/environment/dto.go → `createTemplateMetadataDto`). It is very likely identical since the new API is a facade over the same backend, but not documented.

**Responses**

| Status | Type | Notes |
| --- | --- | --- |
| 201 Created | `OperationExecutionResult` | Media types `text/plain`, `application/json`, `text/json`. |
| 202 Accepted | *(no body type)* | Media types `text/plain`, `application/json`, `text/json`. |
| 400 Bad Request | `ValidationResponse` | |
| 401 Unauthorized | *(none)* | |
| 403 Forbidden | *(none)* | |
| 429 Too Many Requests | *(none)* | |
| Other (409) | `ValidationResponse` | Conflict. |

> Note the dual success codes: **201 with an `OperationExecutionResult` body** (result embedded) **or 202 with no body** (async). A client must handle both.

> **UNVERIFIED — CRITICAL for the provider**: the Learn reference does **not** document which response header carries the polling URL for the 202 case (`Location` vs `Operation-Location`), nor the `Retry-After` behavior. BAPI uses `Location`. The existing provider code for `appmanagement` uses `Operation-Location` (see internal/services/application/resource_environment_application_package_install_test.go). Given the `Operation - Get Operation By ID` endpoint exists at `/environmentmanagement/operations/{operationId}`, the 202 almost certainly returns either that URL or a bare `operationId`. **Must be confirmed by a live capture (mitmproxy per devdocs/adr/mitmproxy.md) before implementing.**

### 3.2 Link Dataverse (add a database to an existing environment)

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/link-dataverse

```http
PATCH https://api.powerplatform.com/environmentmanagement/provisioning/environments/{environmentId}/link?api-version=2024-10-01
```

This is the replacement for BAPI `POST .../environments/{id}/provisionInstance`.

**URI parameters**

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `environmentId` | path | yes | string | The ID of the environment to link Dataverse to. |
| `api-version` | query | yes | string | The API version. |

**Request body** — model `CreateEnvironmentRequestLinkedMetadata` (identical to `linkedEnvironmentMetadata` in 3.1). All fields optional.

| JSON path | Type | Description |
| --- | --- | --- |
| `baseLanguageCode` | integer (int32) | The base language code (for example, `1033` for English). |
| `currency` | `EnvironmentRequestCurrency` | Currency settings for an environment. |
| `currency.code` | string | The currency code (for example, `USD`). |
| `currency.name` | string | The currency name. |
| `currency.precision` | integer (int32) | The currency precision. |
| `currency.symbol` | string | The currency symbol. |
| `domainName` | string | The domain name. |
| `securityGroupId` | string | The security group ID. |
| `templateMetadata` | object | A JSON object payload customized for the selected templates. |
| `templates` | string[] | The templates to apply. |

**Example request body (reconstructed):**

```json
{
  "baseLanguageCode": 1033,
  "domainName": "contoso-dev",
  "securityGroupId": "00000000-0000-0000-0000-000000000003",
  "currency": {
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$",
    "precision": 2
  },
  "templates": []
}
```

**Responses**: `202 Accepted` (no body), `400` `ValidationResponse`, `401`, `403`, `429`, other `ValidationResponse` (Conflict).

Note: unlike Provision New Environment, Link Dataverse has **no 201 path** — it is always async.

### 3.3 Get Supported Locations

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-supported-locations

```http
GET https://api.powerplatform.com/environmentmanagement/provisioning/locations?api-version=2024-10-01
```

**URI parameters**: `api-version` (query, required).

**Response 200** — `ProvisioningLocations`:

| JSON path | Type | Description |
| --- | --- | --- |
| `collection` | `Location[]` | The list of provisioning locations available to the tenant. |
| `locationSelectionMode` | `LocationSelectionMode` enum | Describes how a tenant selects a provisioning location. Tells callers whether to pick a specific location or a macro region. |
| `macroRegions` | `MacroRegion[]` | The list of macro regions available to the tenant. |

`Location` object — "Represents a location/geo for environment provisioning":

| JSON path | Type | Description |
| --- | --- | --- |
| `collection[].name` | string | The location name. |
| `collection[].displayName` | string | The display name. |
| `collection[].code` | string | The location code. |
| `collection[].isDefault` | boolean | Whether this is the default location. |
| `collection[].isDisabled` | boolean | Whether this location is disabled. |
| `collection[].canProvisionDatabase` | boolean | Whether database provisioning is allowed. |
| `collection[].hasFirstReleaseIslandAvailableForProvisioning` | boolean | Whether a first-release island is available for provisioning in this location. |

`MacroRegion` object — "Represents a macro region that groups one or more provisioning locations":

| JSON path | Type | Description |
| --- | --- | --- |
| `macroRegions[].macroRegionId` | string | The macro region identifier. |
| `macroRegions[].displayName` | string | The display name of the macro region. |
| `macroRegions[].dataResidencyNote` | string | The data residency note shown to customers for this macro region. |

`LocationSelectionMode` enum: `Region`, `MacroRegion`.

**Example response (reconstructed from the schema):**

```json
{
  "collection": [
    {
      "name": "unitedstates",
      "displayName": "United States",
      "code": "NAM",
      "isDefault": true,
      "isDisabled": false,
      "canProvisionDatabase": true,
      "hasFirstReleaseIslandAvailableForProvisioning": true
    },
    {
      "name": "europe",
      "displayName": "Europe",
      "code": "EUR",
      "isDefault": false,
      "isDisabled": false,
      "canProvisionDatabase": true,
      "hasFirstReleaseIslandAvailableForProvisioning": true
    }
  ],
  "locationSelectionMode": "Region",
  "macroRegions": [
    {
      "macroRegionId": "americas",
      "displayName": "Americas",
      "dataResidencyNote": "Data is stored within the Americas."
    }
  ]
}
```

**Responses**: `200 OK` (`ProvisioningLocations`), `400` (`ValidationResponse`), `401`, `403`, `404`, `429`.

### 3.4 Get Provisioning Currencies

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-currencies

```http
GET https://api.powerplatform.com/environmentmanagement/provisioning/locations/{location}/currencies?api-version=2024-10-01
```

**URI parameters**

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `location` | path | yes | string | The location name. |
| `api-version` | query | yes | string | The API version. |

**Response 200** — `EnvironmentCurrencyResourceCollection` ("A non-paginated collection of resources returned in full (no continuation token). Used for finite reference lists where partial results would be incorrect."):

| JSON path | Type | Description |
| --- | --- | --- |
| `collection` | `EnvironmentCurrency[]` | |
| `collection[].code` | string | The currency code (for example, `"USD"`). |
| `collection[].symbol` | string | The currency symbol (for example, `"$"`). |
| `collection[].isTenantDefault` | boolean | Whether this is the tenant's default currency. |

**Example response (reconstructed):**

```json
{
  "collection": [
    { "code": "USD", "symbol": "$", "isTenantDefault": true },
    { "code": "EUR", "symbol": "€", "isTenantDefault": false },
    { "code": "GBP", "symbol": "£", "isTenantDefault": false }
  ]
}
```

**Responses**: `200`, `400` (`ValidationResponse`), `401`, `403`, `404`, `429`.

> **Gap vs BAPI**: the new API's currency object has **no `name`** and **no `isDisabled`** — only `code`, `symbol`, `isTenantDefault`. BAPI `environmentCurrencies` returns `properties.code`, `properties.symbol`, `properties.name`, `properties.isTenantDefault`. Note that `EnvironmentRequestCurrency` (the *request* model) does have `name` and `precision`, but the *list* endpoint does not return them.

### 3.5 Get Provisioning Languages

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-languages

```http
GET https://api.powerplatform.com/environmentmanagement/provisioning/locations/{location}/languages?api-version=2024-10-01
```

**URI parameters**: `location` (path, required), `api-version` (query, required).

**Response 200** — `EnvironmentLanguageResourceCollection`:

| JSON path | Type | Description |
| --- | --- | --- |
| `collection` | `EnvironmentLanguage[]` | |
| `collection[].localeId` | integer (int32) | The locale identifier (LCID, for example, `1033` for English). |
| `collection[].localizedName` | string | The language name, localized for display. |
| `collection[].isTenantDefault` | boolean | Whether this is the tenant's default language. |

**Example response (reconstructed):**

```json
{
  "collection": [
    { "localeId": 1033, "localizedName": "English", "isTenantDefault": true },
    { "localeId": 1031, "localizedName": "Deutsch", "isTenantDefault": false },
    { "localeId": 1036, "localizedName": "Français", "isTenantDefault": false }
  ]
}
```

**Responses**: `200`, `400` (`ValidationResponse`), `401`, `403`, `404`, `429`.

> **Gap vs BAPI**: no `displayName` (BAPI returns both `localizedName` and `displayName`), no `isDisabled`, no `localeName`/`code`.

### 3.6 Get Provisioning Templates

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-templates

```http
GET https://api.powerplatform.com/environmentmanagement/provisioning/locations/{location}/templates?api-version=2024-10-01
```

**URI parameters**: `location` (path, required), `api-version` (query, required).

**Response 200** — `EnvironmentTemplateResourceCollection`:

| JSON path | Type | Description |
| --- | --- | --- |
| `collection` | `EnvironmentTemplate[]` | |
| `collection[].name` | string | The template name (identifier). |
| `collection[].displayName` | string | The template name, localized for display. |
| `collection[].isCustomerEngagement` | boolean | Whether the template is a Customer Engagement template. |
| `collection[].isSupportedForResetOperation` | boolean | Whether this template is supported for the reset operation. |
| `collection[].availability` | `TemplateAvailability[]` | The per-SKU availability of this template. |
| `collection[].availability[].environmentSku` | `EnvironmentSku` enum | The environment SKU. |
| `collection[].availability[].isDisabled` | boolean | Whether the template is disabled for this SKU. |
| `collection[].availability[].disabledReason` | `DisabledReason` | Explains why a template is unavailable for a given SKU. |
| `collection[].availability[].disabledReason.code` | string | The reason code. |
| `collection[].availability[].disabledReason.message` | string | The reason message. |

**Example response (reconstructed):**

```json
{
  "collection": [
    {
      "name": "D365_Sales",
      "displayName": "Dynamics 365 Sales",
      "isCustomerEngagement": true,
      "isSupportedForResetOperation": true,
      "availability": [
        { "environmentSku": "Sandbox", "isDisabled": false },
        { "environmentSku": "Production", "isDisabled": false },
        {
          "environmentSku": "Developer",
          "isDisabled": true,
          "disabledReason": {
            "code": "SkuNotSupported",
            "message": "This template is not available for Developer environments."
          }
        }
      ]
    }
  ]
}
```

**Responses**: `200`, `400` (`ValidationResponse`), `401`, `403`, `404`, `429`.

> Replacement for the provider's `powerplatform_environment_templates` data source (BAPI `.../locations/{location}/environmentTemplates`), plus richer per-SKU availability data.

---

## 4. Environments (read) — full `EnvironmentResponse` model

Sources:

* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/list-environments-for-user
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user

> **CRITICAL STRUCTURAL DIFFERENCE**: The new API's environment model is **FLAT**. There is **no `properties` wrapper** and **no `linkedEnvironmentMetadata` sub-object**. Dataverse fields (`url`, `domainName`, `version`, `dataverseId`, `securityGroupId`) are promoted to the top level. Every existing JSON path in internal/services/environment/dto.go changes.

### 4.1 List Environments For User

```http
GET https://api.powerplatform.com/environmentmanagement/environments?api-version=2024-10-01
```

With optional parameters:

```http
GET https://api.powerplatform.com/environmentmanagement/environments?ids={ids}&$filter={$filter}&$select={$select}&$top={$top}&$skip={$skip}&$orderby={$orderby}&api-version=2024-10-01
```

**URI parameters**

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `api-version` | query | yes | string | The API version. |
| `ids` | query | no | string[] | Comma-separated list of environment IDs to retrieve. When specified, only environments matching these IDs are returned. |
| `$filter` | query | no | string | OData filter expression. **Supported filter properties: `dataverseId`, `type`, `geo`, `state`, `environmentGroupId`, `domainName`.** |
| `$orderby` | query | no | string | OData order-by expression for sorting the results. |
| `$select` | query | no | string | Comma-separated list of properties to include in the response. |
| `$skip` | query | no | integer (min 0) | Number of environments to skip before returning results. |
| `$top` | query | no | integer (min 0) | Maximum number of environments to return. |

Note: **no `$expand`**. The response is flat and complete; there is no capacity/addons/permissions expansion.

**Response 200** — `EnvironmentList`:

| JSON path | Type | Description |
| --- | --- | --- |
| `value` | `EnvironmentResponse[]` | |
| `@odata.nextlink` | string (uri) | Opaque URL to retrieve the next page of results. Present only when additional pages are available. |

### 4.2 Get Environment By Id For User

```http
GET https://api.powerplatform.com/environmentmanagement/environments/{environmentId}?api-version=2024-10-01
GET https://api.powerplatform.com/environmentmanagement/environments/{environmentId}?$select={$select}&api-version=2024-10-01
```

**URI parameters**: `environmentId` (path, required, string), `api-version` (query, required), `$select` (query, optional).

**Response 200**: `EnvironmentResponse` (bare object, no envelope).

### 4.3 `EnvironmentResponse` — every property

| JSON path | Type | Description |
| --- | --- | --- |
| `id` | string | The ID of the environment. |
| `displayName` | string | The display name of the environment. |
| `type` | string | The type (SKU) of the environment. |
| `state` | string | The current state of the environment. |
| `tenantId` | string | The ID of the tenant that the environment belongs to. |
| `geo` | string | The geographical region of the environment. |
| `azureRegion` | string | The Azure region of the environment. |
| `clusterCategory` | string | The cluster category the environment is in. |
| `adminMode` | string | Indicates whether admin-only mode is enabled or disabled for the environment. |
| `backgroundOperationsState` | string | Indicates whether background operations are enabled or disabled for the environment. |
| `protectionLevel` | string | The protection level applied to the environment. |
| `environmentGroupId` | string | The ID of the environment group to which this environment belongs. |
| `connectedGroupId` | string | The ID of the AAD group connected to the environment. |
| `securityGroupId` | string | The security group that controls access to the environment. |
| `scenarioName` | string | The scenario name associated with the environment (for example, singleton scenario type). |
| `dataverseId` | string | The ID of the Dataverse database (organization) associated with the environment. |
| `url` | string | The URL of the Dataverse database associated with the environment. |
| `domainName` | string | The domain name of the Dataverse database associated with the environment. |
| `version` | string | The version of the Dataverse database associated with the environment. |
| `createdDateTime` | string (date-time) | The creation date and time of the environment. |
| `createdBy` | `EnvironmentPrincipal` | Represents a principal (user or application). |
| `createdBy.id` | string | The principal ID. |
| `createdBy.type` | string | The principal type. |
| `createdFor` | `EnvironmentPrincipal` | Represents a principal (user or application). |
| `createdFor.id` | string | The principal ID. |
| `createdFor.type` | string | The principal type. |
| `deletedDateTime` | string (date-time) | The deletion date and time of the environment. |
| `retentionDetails` | `RetentionDetails` | The retention details of the environment. |
| `retentionDetails.retentionPeriod` | string | The retention period of the environment. |
| `retentionDetails.availableFromDateTime` | string (date-time) | The date and time from which the environment is available for recovery. |
| `finOpsMetadata` | `FinOpsMetadata` | Metadata describing a linked FinOps environment. |
| `finOpsMetadata.id` | string | The linked FinOps environment ID. |
| `finOpsMetadata.type` | string | The linked FinOps environment type. |
| `finOpsMetadata.url` | string | The linked FinOps environment URL. |
| `enterprisePolicies` | `EnterprisePolicies` | The set of enterprise policies linked to the environment. |
| `enterprisePolicies.encryption` | `EnterprisePolicyLink` | |
| `enterprisePolicies.identity` | `EnterprisePolicyLink` | |
| `enterprisePolicies.networkInjection` | `EnterprisePolicyLink` | |
| `enterprisePolicies.privateEndpoint` | `EnterprisePolicyLink` | |
| `enterprisePolicies.*.id` | string | The ID of the enterprise policy. |
| `enterprisePolicies.*.resourceId` | string | The fully-qualified Azure resource ID of the enterprise policy. |
| `enterprisePolicies.*.status` | `EnterprisePolicyLinkStatus` enum | The status of the link. |
| `enterprisePolicies.*.error` | string | Error details when the link status is `Failed`. |

`EnterprisePolicyLinkStatus` enum: `Linking`, `Unlinking`, `Linked`, `Failed`, `LinkingOnline`, `UnlinkingOnline`.

**Example response (reconstructed from the schema — the Learn pages publish no example payload):**

```json
{
  "value": [
    {
      "id": "00000000-0000-0000-0000-000000000001",
      "displayName": "Contoso Dev",
      "type": "Sandbox",
      "state": "Ready",
      "tenantId": "00000000-0000-0000-0000-000000000005",
      "geo": "unitedstates",
      "azureRegion": "westus",
      "clusterCategory": "Prod",
      "adminMode": "Disabled",
      "backgroundOperationsState": "Enabled",
      "protectionLevel": "Standard",
      "environmentGroupId": "00000000-0000-0000-0000-000000000002",
      "connectedGroupId": null,
      "securityGroupId": "00000000-0000-0000-0000-000000000003",
      "scenarioName": null,
      "dataverseId": "a0a0a0a0-bbbb-cccc-dddd-e1e1e1e1e1e1",
      "url": "https://org0fadb1dd.crm.dynamics.com/",
      "domainName": "org0fadb1dd",
      "version": "9.2.21013.00152",
      "createdDateTime": "2020-10-22T04:38:17.8550157Z",
      "createdBy": {
        "id": "0f747967-84c4-4f29-84c2-682fb00390c8",
        "type": "ServicePrincipal"
      },
      "createdFor": {
        "id": "0f747967-84c4-4f29-84c2-682fb00390c8",
        "type": "User"
      },
      "deletedDateTime": null,
      "retentionDetails": {
        "retentionPeriod": "P7D",
        "availableFromDateTime": "2021-02-16T05:42:52.2822636Z"
      },
      "finOpsMetadata": null,
      "enterprisePolicies": {
        "encryption": {
          "id": "00000000-0000-0000-0000-000000000010",
          "resourceId": "/subscriptions/.../providers/Microsoft.PowerPlatform/enterprisePolicies/cmk-policy",
          "status": "Linked",
          "error": null
        },
        "identity": null,
        "networkInjection": null,
        "privateEndpoint": null
      }
    }
  ],
  "@odata.nextlink": "https://api.powerplatform.com/environmentmanagement/environments?api-version=2024-10-01&$skip=100"
}
```

> **UNVERIFIED**: enum value sets for `state`, `adminMode`, `backgroundOperationsState`, `protectionLevel`, `type`, `clusterCategory` are **not published** — all six are typed as bare `string` with no enumeration. From BAPI experience the likely values are:
>
> * `state`: `Ready`, `NotSpecified`, `Disabled`, `AdminMode`, `Deleted`, `Soft Deleted`
> * `adminMode`: `Enabled` / `Disabled`
> * `backgroundOperationsState`: `Enabled` / `Disabled`
> * `protectionLevel`: `Basic` / `Standard` (matches the `ProtectionLevel` enum in the provisioning models)
> * `type`: the `EnvironmentSku` values (see 4.4)
>
> These must be confirmed against live responses before writing validators.

### 4.4 `EnvironmentSku` enum (12 values)

Documented on Provision New Environment and Get Provisioning Templates:

`Standard`, `Premium`, `Developer`, `Basic`, `Production`, `Sandbox`, `Trial`, `Default`, `Support`, `SubscriptionBasedTrial`, `Teams`, `Platform`

> The provider today only allows `Developer`, `Sandbox`, `Production`, `Trial`, `Default` (internal/services/environment/dto.go → `EnvironmentTypes`). The new API documents seven additional SKUs.

### 4.5 `ProtectionLevel` enum

`Basic`, `Standard`

---

## 5. Lifecycle operations

### 5.1 Delete Environment By ID

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-delete/delete-environment-by-id

```http
DELETE https://api.powerplatform.com/environmentmanagement/environments/{environmentId}?api-version=2024-10-01
DELETE https://api.powerplatform.com/environmentmanagement/environments/{environmentId}?ValidateOnly={ValidateOnly}&ValidateProperties={ValidateProperties}&api-version=2024-10-01
```

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `environmentId` | path | yes | string | The ID of the environment. |
| `api-version` | query | yes | string | The API version. |
| `ValidateOnly` | query | no | boolean | When true, validates the request without executing it. Use with `validateProperties` to validate only specific fields. If `validateProperties` is empty, the entire request is validated. Defaults to `false` (validate and execute). |
| `ValidateProperties` | query | no | string | A comma-separated list of property names to validate (for example, `"property1,property2"`). Applies only when `validateOnly` is true. |

Responses: `202 Accepted` (no body), `400` `ValidationResponse`, `401`, `403`, `404`, `429`, other `ValidationResponse` (Conflict).

> **`ValidateOnly` is the direct replacement for BAPI `POST .../environments/{id}/validateDelete`.** This `ValidateOnly`/`ValidateProperties` pair appears on **all** write-side lifecycle operations (delete, recover, reset, copy, enable/disable managed) — a consistent pre-flight-validation pattern, and a strong fit for Terraform `ValidateConfig` / plan-time checks.

> **UNVERIFIED**: `ValidateOnly` and `ValidateProperties` are documented with capitalized first letters in the URI-parameter table but lowercase in the description text ("when `validateOnly` is true"). Query-string casing must be confirmed empirically.

### 5.2 Recover Environment

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-recover/recover-environment

```http
POST https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/recover?api-version=2024-10-01
POST https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/recover?ValidateOnly={ValidateOnly}&ValidateProperties={ValidateProperties}&api-version=2024-10-01
```

No request body documented. Responses: `202`, `400` `ValidationResponse`, `401`, `403`, `429`, other `ValidationResponse` (Conflict). Note: no `404` listed.

### 5.3 Reset Environment

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-reset/reset-environment

```http
POST https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/reset?api-version=2024-10-01
POST https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/reset?ValidateOnly={ValidateOnly}&ValidateProperties={ValidateProperties}&api-version=2024-10-01
```

**Request body** — model `ResetRequest`, all fields optional:

| JSON path | Type | Description |
| --- | --- | --- |
| `displayName` | string | The display name for the environment to reset to. |
| `description` | string | An optional description for the environment to reset to. |
| `domainName` | string | Domain name for the environment to reset to. |
| `securityGroupId` | string | Security group ID for the environment to reset to. |
| `baseLanguageCode` | integer (int32) | The base language code (for example, `1033` for English) for the environment to reset to. |
| `currency` | `EnvironmentRequestCurrency` | Currency settings for an environment (`code`, `name`, `precision`, `symbol`). |
| `templates` | string[] | Templates to apply for the environment after reset. |

Responses: `202`, `400`, `401`, `403`, `429`, Conflict.

### 5.4 Copy Environment

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-copy/copy-environment

```http
POST https://api.powerplatform.com/environmentmanagement/environments/{targetEnvironmentId}/copy?api-version=2024-10-01
POST https://api.powerplatform.com/environmentmanagement/environments/{targetEnvironmentId}/copy?ValidateOnly={ValidateOnly}&ValidateProperties={ValidateProperties}&api-version=2024-10-01
```

`targetEnvironmentId` = "The ID of the target environment that will be overwritten."

**Request body** — model `CopyRequest`:

| JSON path | Required | Type | Description |
| --- | --- | --- | --- |
| `copyType` | **yes** | `CopyType` enum | Represents the type of copy operation. |
| `sourceEnvironmentId` | **yes** | string | Source environment ID to copy from. |
| `copyOptions` | no | `CopyRequestOptions` | Optional inputs for copy request. |
| `copyOptions.environmentNameToOverride` | no | string | Environment name to override on target environment. |
| `copyOptions.securityGroupIdToOverride` | no | string | Security group ID to override on target environment. |
| `copyOptions.skipAuditData` | no | boolean | Boolean flag to skip audit data for copy. |
| `copyOptions.executeAdvancedCopyForFinanceAndOperations` | no | boolean | Boolean flag to execute advanced copy for Finance and Operations data. |

`CopyType` enum: `Minimal`, `Full`.

Responses: `202`, `400`, `401`, `403`, `429`, Conflict.

### 5.5 Enable / Disable Managed Environment

Sources:

* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-managed-governance/enable-managed-environment
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-managed-governance/disable-managed-environment

```http
POST https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/governanceSetting/enableManaged?api-version=2024-10-01
POST https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/governanceSetting/disableManaged?api-version=2024-10-01
```

Both support `ValidateOnly` / `ValidateProperties`. No request body documented. Responses: `202`, `400`, `401`, `403`, `429`, Conflict.

> **Gap**: these toggle managed-environment on/off but expose **no way to set the managed-environment extended settings** (`limitSharingMode`, `maxLimitUserSharing`, `solutionCheckerMode`, `suppressValidationEmails`, `solutionCheckerRuleOverrides`, `isGroupSharingDisabled`, `excludeEnvironmentFromAnalysis`, `includeOnHomepageInsights`, `disableAiGeneratedDescriptions`, the `solutionCloudFlows-*` and `bot-*` keys). The provider currently reads/writes all of these via BAPI `properties.governanceConfiguration.settings.extendedSettings` (internal/services/environment/dto.go → `ExtendedSettingsDto`). **The new API's `EnvironmentResponse` exposes only a flat `protectionLevel` string.** This is the single biggest functional gap for the provider's managed-environment behavior.

### 5.6 Enable / Disable Environment (admin mode)

Listed in the Oct 2025 what's-new entry as `environmentmanagement/environment-state/enable-environment` and `.../disable-environment`, but **both doc pages return HTTP 404 as of 2026-08-18**, as does the `environment-state` operation-group page.

> **UNVERIFIED — needs confirmation**: exact paths. By analogy with the managed-governance pattern the likely shapes are `POST /environmentmanagement/environments/{environmentId}/enable` and `POST /environmentmanagement/environments/{environmentId}/disable`, or `.../state/enable` / `.../state/disable`. The `EnvironmentResponse.adminMode` and `EnvironmentResponse.state` fields are the read-side counterparts. Must be confirmed via the .NET/Python SDK or a live capture.

### 5.7 Update Environment — DOES NOT EXIST

> **KEY FINDING (blocking for the migration).** There is **no `PATCH /environmentmanagement/environments/{environmentId}`** and no "Update Environment" / "Update Environment Properties" operation anywhere in the `environmentmanagement` namespace. Confirmed by:
>
> * The full what's-new endpoint list (Dec 2024 → Jul 2026) contains no such entry.
> * The `Environments` operation group lists only `Get Environment By Id For User` and `List Environments For User`.
> * There is no `Environment-Update` operation group.
>
> Consequences: **changing `display_name`, `description`, `security_group_id`, `domain_name`, `environment_sku`/type, `cadence`, `billing_policy_id`, or `environment_group_id` on an existing environment has no documented Power Platform API equivalent.** The provider today does all of this via BAPI `PATCH https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{id}?api-version=2021-04-01`.
>
> Partial workarounds available on the new API:
>
> * `environmentGroupId` → `POST /environmentmanagement/environmentGroups/{groupId}/addEnvironment/{environmentId}` and the corresponding remove.
> * `protectionLevel` → `governanceSetting/enableManaged` / `disableManaged`.
> * `displayName`, `description`, `domainName`, `securityGroupId`, `baseLanguageCode`, `currency`, `templates` → only via the **destructive** `POST .../reset` (which wipes the Dataverse database — not an acceptable Terraform in-place update).
> * `billingPolicy` → `POST https://api.powerplatform.com/licensing/billingPolicies/{id}/environments/add|remove?api-version=2022-03-01-preview` (already used by the provider).
>
> **Recommendation: a full migration is not currently possible. Plan a hybrid — reads + create + delete on the new API, in-place updates still on BAPI — or defer until an Update Environment endpoint ships.**

---

## 6. Async operations and polling

### 6.1 Polling endpoints

| Operation | Method | Path |
| --- | --- | --- |
| Get Operation By ID | GET | `/environmentmanagement/operations/{operationId}?api-version=2024-10-01` |
| Get Environment Operation By ID | GET | `/environmentmanagement/environments/{targetEnvironmentId}/operations/{operationId}?api-version=2024-10-01` |
| Get Operations For Environment | GET | `/environmentmanagement/environments/{environmentId}/operations?api-version=2024-10-01` |

Sources:

* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-operation-by-id
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-environment-operation-by-id
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-operations-for-environment

`Get Environment Operation By ID` is described as: "Gets the status of an environment lifecycle operation scoped under a specific environment, **enabling environment-level authorization on the operation lookup**." Prefer this variant when the environment ID is known — it avoids requiring tenant-wide operation-read rights.

`Get Operations For Environment` extra query parameters:

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `limit` | query | no | string | The maximum number of records to return per request. Must be a positive integer; a server default applies if omitted. |
| `continuationToken` | query | no | string | An opaque token returned by a previous response, used to fetch the next page of results. Omit to retrieve the first page. |

Response `OperationExecutionResultPagedCollection`:

| JSON path | Type |
| --- | --- |
| `collection` | `OperationExecutionResult[]` |
| `continuationToken` | string |

### 6.2 `OperationExecutionResult` model

"The result of an environment lifecycle operation."

| JSON path | Type | Description |
| --- | --- | --- |
| `operationId` | string | The ID of the operation. |
| `name` | string | The name of the operation. |
| `status` | `OperationStatus` enum | The status of operation. |
| `startTime` | string (date-time) | The start time of the operation. |
| `endTime` | string (date-time) | The end time of the operation. |
| `requestedBy` | `UserIdentity` | Represents the identity of a user. |
| `requestedBy.userId` | string | The ID of the user. |
| `requestedBy.displayName` | string | The display name of the user. |
| `requestedBy.type` | string | The type of the user identity (for example, `User`). |
| `requestedBy.tenantId` | string | The tenant ID of the user. |
| `errorDetail` | `OperationErrorDetail` | Structured error detail for a failed request. |
| `errorDetail.code` | string | The error code. |
| `errorDetail.fieldErrors` | map<string, `FieldError`> | Per-field error detail, keyed by field name. |
| `errorDetail.fieldErrors.{field}.errorMessages` | string[] | The error messages describing what is wrong with the field. |
| `errorDetail.fieldErrors.{field}.suggestedValue` | string | A suggested or accepted value that would resolve the error. |
| `stageStatuses` | `StageStatus[]` | Per-stage progress of the operation. |
| `stageStatuses[].name` | string | The name of the stage. |
| `stageStatuses[].status` | `StepExecutionStatus` enum | The execution status of an operation stage. |
| `stageStatuses[].startTime` | string (date-time) | The start time of the stage. |
| `stageStatuses[].endTime` | string (date-time) | The end time of the stage. |
| `stageStatuses[].errorDetail` | `OperationErrorDetail` | |
| `updatedEnvironment` | `Environment` | Power Platform environment (minimal projection). |
| `updatedEnvironment.environmentId` | string | The environment ID. |
| `updatedEnvironment.displayName` | string | Display name of the environment. |
| `updatedEnvironment.dataverseOrganizationUrl` | string | Dataverse organization URL of the environment. |

> `updatedEnvironment.environmentId` is how a client learns the ID of a newly provisioned environment. This is the analogue of BAPI's `lifecycleOperations` response `properties.linkedEnvironmentMetadata.resourceId`.

### 6.3 Status enums

`OperationStatus` (7 values):

| Value | Terminal? |
| --- | --- |
| `Queued` | no |
| `InProgress` | no |
| `Succeeded` | **yes** (success) |
| `ValidationFailed` | **yes** (failure) |
| `Failed` | **yes** (failure) |
| `NoOperation` | **yes** (no-op) |
| `ValidationPassed` | **yes** for `ValidateOnly=true` calls |

`StepExecutionStatus` (6 values): `Succeeded`, `Failed`, `Skipped`, `Postponed`, `InProgress`, `NotStarted`.

Compare with BAPI `provisioningState`: `Succeeded`, `Failed`, `Deleting`, `Deleted`, `Creating`, `Updating`, `NotSpecified`, `Accepted`.

**Example polling response (reconstructed):**

```json
{
  "operationId": "b03e1e6d-73db-4367-90e1-2e378bf7e2fc",
  "name": "ProvisionEnvironment",
  "status": "InProgress",
  "startTime": "2026-08-18T10:00:00Z",
  "endTime": null,
  "requestedBy": {
    "userId": "0f747967-84c4-4f29-84c2-682fb00390c8",
    "displayName": "Jane Doe",
    "type": "User",
    "tenantId": "00000000-0000-0000-0000-000000000005"
  },
  "errorDetail": null,
  "stageStatuses": [
    { "name": "Validation", "status": "Succeeded", "startTime": "2026-08-18T10:00:00Z", "endTime": "2026-08-18T10:00:05Z" },
    { "name": "CreateEnvironment", "status": "Succeeded", "startTime": "2026-08-18T10:00:05Z", "endTime": "2026-08-18T10:01:00Z" },
    { "name": "ProvisionDataverse", "status": "InProgress", "startTime": "2026-08-18T10:01:00Z", "endTime": null }
  ],
  "updatedEnvironment": {
    "environmentId": "00000000-0000-0000-0000-000000000001",
    "displayName": "Contoso Dev",
    "dataverseOrganizationUrl": null
  }
}
```

**Example failure response (reconstructed):**

```json
{
  "operationId": "b03e1e6d-73db-4367-90e1-2e378bf7e2fc",
  "name": "ProvisionEnvironment",
  "status": "ValidationFailed",
  "errorDetail": {
    "code": "InvalidDomainName",
    "fieldErrors": {
      "linkedEnvironmentMetadata.domainName": {
        "errorMessages": ["The domain name 'contoso-dev' is already in use."],
        "suggestedValue": "contoso-dev1"
      }
    }
  }
}
```

> The `fieldErrors[].suggestedValue` mechanism is a notable improvement over BAPI — it returns an available alternative domain name directly, replacing the provider's current BAPI `POST /providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails?api-version=2021-04-01` round-trip.

### 6.4 `ValidationResponse` (400/409 body)

| JSON path | Type |
| --- | --- |
| `errorDetail` | `OperationErrorDetail` |
| `errorDetail.code` | string |
| `errorDetail.fieldErrors` | map<string, `FieldError`> |

Environment Groups operations use a different error shape, RFC 7807 `ProblemDetails`: `type`, `title`, `status` (int32), `detail`, `instance`, `extensions`.

### 6.5 Async response headers

> **UNVERIFIED — CRITICAL**. None of the `environmentmanagement` Learn pages document the response headers on `202 Accepted`. Unknown: whether the polling URL arrives via `Location`, `Operation-Location`, or an `operationId` in the body; and whether `Retry-After` is emitted.
>
> Evidence for the two candidates in this repo:
>
> * BAPI uses `Location` → `https://{geo}.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/{id}?api-version=2023-06-01` (see internal/services/environment/resource_environment_test.go).
> * `api.powerplatform.com/appmanagement` uses `Operation-Location` → `https://api.powerplatform.com/appmanagement/environments/{envId}/operations/{opId}?api-version=1` (see internal/services/application/resource_environment_application_package_install_test.go).
>
> Given `environmentmanagement` follows the same host and versioning conventions as `appmanagement`, `Operation-Location` pointing at `/environmentmanagement/environments/{envId}/operations/{opId}?api-version=2024-10-01` is the most likely shape — **but must be captured live** (see devdocs/adr/mitmproxy.md and devdocs/adr/httpmocks.md) before mock responders are written.

---

## 7. Environment Management Settings

Sources:

* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/list-environment-management-settings
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/update-environment-management-settings

```http
GET   https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/settings?api-version=2024-10-01
GET   https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/settings?$top={$top}&$select={$select}&api-version=2024-10-01
PATCH https://api.powerplatform.com/environmentmanagement/environments/{environmentId}/settings?api-version=2024-10-01
```

`$top` defaults to 500 if not set.

**Response envelope** `GetEnvironmentManagementSettingResponse`:

| JSON path | Type | Description |
| --- | --- | --- |
| `objectResult` | `EnvironmentManagementSetting[]` | Gets or sets the fields for the entities being queried. |
| `nextLink` | string (uri) | Gets or sets the next link if there are more records to be returned. |
| `responseMessage` | string | Gets or sets the error message. |
| `errors` | `EnvironmentServiceErrorResponse` | |
| `errors.code` | string | |
| `errors.message` | string | |
| `errors.details` | `ErrorDetail[]` | |
| `errors.details[].code` / `.message` / `.target` / `.value` | string | |

`EnvironmentManagementSetting` properties (complete list as published):

`id`, `tenantId` (uuid), `allowedIpRangeForStorageAccessSignatures`, `enableIpBasedStorageAccessSignatureRule`, `ipBasedStorageAccessSignatureMode` (int32), `loggingEnabledForIpBasedStorageAccessSignature`, `copilotStudio_CodeInterpreter`, `copilotStudio_ComputerUseAppAllowlist`, `copilotStudio_ComputerUseCredentialsAllowed`, `copilotStudio_ComputerUseSharedMachines`, `copilotStudio_ComputerUseWebAllowlist`, `copilotStudio_ConnectedAgents`, `copilotStudio_ConversationAuditLoggingEnabled`, `d365CustomerService_AIAgents`, `d365CustomerService_Copilot`, `powerApps_AllowCodeApps`, `powerApps_ChartVisualization`, `powerApps_CopilotChat`, `powerApps_EnableFormInsights`, `powerApps_FormPredictAutomatic`, `powerApps_FormPredictSmartPaste`, `powerApps_NLSearch`, `powerPages_AllowIntelligentFormsCopilotForSites`, `powerPages_AllowListSummaryCopilotForSites`, `powerPages_AllowMakerCopilotsForExistingSites`, `powerPages_AllowMakerCopilotsForNewSites`, `powerPages_AllowNonProdPublicSites`, `powerPages_AllowNonProdPublicSites_Exemptions`, `powerPages_AllowProDevCopilotsForEnvironment`, `powerPages_AllowProDevCopilotsForSites`, `powerPages_AllowSearchSummaryCopilotForSites`, `powerPages_AllowSiteCopilotForSites`, `powerPages_AllowSummarizationAPICopilotForSites`

`PATCH` response is `OperationResponse` (`objectResult`, `nextLink`, `responseMessage`, `errors`, `debugErrors`).

> **Important**: this is **NOT** the same surface as the provider's `powerplatform_environment_settings` data source/resource (which reads Dataverse `organizations` via the environment's Web API), nor is it the managed-environment `extendedSettings`. These are AI/Copilot governance toggles. It does NOT cover `allowBingSearch` / `bingChatEnabled`, `allowMovingDataAcrossRegions` / `crossGeoCopilotDataMovementEnabled`, or `m365Enabled` — those remain BAPI-only (`properties.copilotPolicies`, `properties.bingChatEnabled`, `properties.m365Enabled`, and the `generativeAiFeatures` sub-resource).

Tutorial: https://learn.microsoft.com/en-us/power-platform/admin/programmability-tutorial-environmentmanagement-settings

---

## 8. Environment Groups

Source: https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-groups

| Operation | Method | Path |
| --- | --- | --- |
| Get Environment Group | GET | `/environmentmanagement/environmentGroups/{groupId}?api-version=2024-10-01` |
| Add Environment To Group | POST | `/environmentmanagement/environmentGroups/{groupId}/addEnvironment/{environmentId}?api-version=2024-10-01` |
| Get Environment Group Operation | GET | `/environmentmanagement/environmentGroupOperations/{operationId}?api-version=2024-10-01` |
| Create / List / Update / Delete / Remove Environment From Group | — | *(paths UNVERIFIED — not fetched)* |

`EnvironmentGroup` model:

| JSON path | Type |
| --- | --- |
| `id` | string (uuid) |
| `displayName` | string |
| `description` | string |
| `parentGroupId` | string (uuid) |
| `childrenGroupIds` | string[] (uuid) |
| `createdBy` | `Principal` |
| `createdTime` | string (date-time) |
| `lastModifiedBy` | `Principal` |
| `lastModifiedTime` | string (date-time) |

`Principal`: `id`, `displayName`, `email`, `userPrincipalName`, `type`, `tenantId`.

`Add Environment To Group` responses: `202 Accepted`, `204 No Content`, `400 ProblemDetails`. `Get Environment Group Operation` responses: `200 OK` (no typed body documented), `204 No Content`, `400 ProblemDetails`.

> Note the **separate operation-status endpoint** for group operations (`environmentGroupOperations`) — group operations do NOT share the `environmentmanagement/operations/{operationId}` polling path.

> The provider currently manages environment-group rule sets via the **tenant-scoped host** `https://{tenantId-hex}.{NN}.tenant.api.powerplatform.com/governance/environmentGroups/{groupId}/ruleSets?api-version=2021-10-01-preview` (see internal/services/environment_group_rule_set/resource_environment_group_rule_set_test.go). That is a different namespace (`governance`) and host pattern, out of scope for this migration.

---

## 9. Legacy BAPI baseline (for the mapping table)

### 9.1 BAPI endpoints the provider uses today

From internal/services/environment/api_environment.go and internal/services/environment/resource_environment_test.go:

| Purpose | Method | BAPI URL |
| --- | --- | --- |
| List environments | GET | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01&$expand=properties/billingPolicy,properties/copilotPolicies` |
| Get environment | GET | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{id}?api-version=2023-06-01&$expand=permissions,properties.capacity,properties/billingPolicy,properties/copilotPolicies` |
| Create environment | POST | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01` |
| Update environment | PATCH | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{id}?api-version=2021-04-01` |
| Delete environment | DELETE | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{id}` |
| Validate delete | POST | `.../scopes/admin/environments/{id}/validateDelete` |
| Add Dataverse DB | POST | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/{id}/provisionInstance?api-version=2021-04-01` |
| Validate name/domain | POST | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/validateEnvironmentDetails?api-version=2021-04-01` |
| Poll lifecycle op | GET | `https://{geo}.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/{id}?api-version=2023-06-01` (from the `Location` header) |
| List locations | GET | `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations` |
| List currencies | GET | `.../locations/{location}/environmentCurrencies` |
| List languages | GET | `.../locations/{location}/environmentLanguages` |
| List templates | GET | `.../locations/{location}/environmentTemplates` |

### 9.2 BAPI `EnvironmentDto` (provider's current model)

From internal/services/environment/dto.go. Top level: `id`, `type`, `location`, `name`, `properties{...}`.

`properties.*`: `azureRegion`, `databaseType`, `displayName`, `environmentSku`, `linkedAppMetadata{id,type,url}`, `runtimeEndpoints{microsoft.BusinessAppPlatform, microsoft.CommonDataModel, microsoft.PowerApps, microsoft.PowerAppsAdvisor, microsoft.PowerVirtualAgents, microsoft.ApiManagement, microsoft.Flow}`, `linkedEnvironmentMetadata{backgroundOperationsState, domainName, instanceUrl, baseLanguage, securityGroupId, resourceId, version, template[], templateMetadata, uniqueName}`, `states{management{id}, runtime{id}, disasterRecovery{id}}`, `tenantId`, `governanceConfiguration{protectionLevel, settings.extendedSettings{...}}`, `billingPolicy{id}`, `provisioningState`, `description`, `updateCadence{id}`, `parentEnvironmentGroup{id}`, `enterprisePolicies{vnets, customerManagedKeys, identity}` each `{policyId, location, id, systemId, linkStatus}`, `cluster{category}`, `usedBy{id,type,tenantId}`, `bingChatEnabled`, `m365Enabled`, `copilotPolicies{crossGeoCopilotDataMovementEnabled, crossBoundaryCopilotDataMovementEnabled}`.

Plus (from the published BAPI example, https://learn.microsoft.com/en-us/power-platform/admin/list-environments): `createdTime`, `createdBy{id,displayName,type,tenantId}`, `lastModifiedTime`, `creationType`, `isDefault`, `capacity[]{capacityType, actualConsumption, ratedConsumption, capacityUnit, updatedOn}`, `addons[]{addonType, allocated, addonUnit}`, `clientUris{admin, maker}`, `notificationMetadata{state, branding}`, `retentionPeriod`, `retentionDetails{retentionPeriod, backupsAvailableFromDateTime}`, `protectionStatus{keyManagedBy}`, `cluster{number}`, `connectedGroups[]`, and inside `linkedEnvironmentMetadata`: `friendlyName`, `uniqueName`, `createdTime`, `scaleGroup`, `platformSku`, `instanceApiUrl`, `instanceState`.

### 9.3 BAPI `LocationDto` (provider's current model)

`id`, `type`, `name`, `properties{displayName, code, isDisabled, isDefault, canProvisionDatabase, canProvisionCustomerEngagementDatabase, azureRegions[]}`.

---

## 10. Mapping table — BAPI to Power Platform API

### 10.1 Endpoint mapping

| Purpose | BAPI | Power Platform API | Status |
| --- | --- | --- | --- |
| List environments | `GET .../scopes/admin/environments?$expand=properties` | `GET /environmentmanagement/environments` | available (flat shape; `$filter`/`$select`/`$top`/`$skip`/`$orderby`/`ids`) |
| Get environment | `GET .../scopes/admin/environments/{id}?$expand=...` | `GET /environmentmanagement/environments/{id}` | available (`$select` only) |
| Create environment | `POST .../environments` | `POST /environmentmanagement/provisioning/environments` | available |
| **Update environment** | `PATCH .../scopes/admin/environments/{id}` | **none** | **BLOCKING GAP** |
| Delete environment | `DELETE .../scopes/admin/environments/{id}` | `DELETE /environmentmanagement/environments/{id}` | available |
| Validate delete | `POST .../environments/{id}/validateDelete` | `DELETE /environmentmanagement/environments/{id}?ValidateOnly=true` | available (better pattern) |
| Add Dataverse DB | `POST .../environments/{id}/provisionInstance` | `PATCH /environmentmanagement/provisioning/environments/{id}/link` | available |
| Validate name/domain | `POST .../validateEnvironmentDetails` | Provision `ValidateOnly` support UNVERIFIED; `errorDetail.fieldErrors[].suggestedValue` on 400 covers the use case | partial |
| Recover soft-deleted | BAPI recover | `POST /environmentmanagement/environments/{id}/recover` | available |
| Poll lifecycle op | `GET https://{geo}.api.bap.microsoft.com/.../lifecycleOperations/{id}` | `GET /environmentmanagement/operations/{opId}` or `GET /environmentmanagement/environments/{envId}/operations/{opId}` | available (single global host — no geo-specific host) |
| List locations | `GET .../locations` | `GET /environmentmanagement/provisioning/locations` | available, but property gaps |
| List currencies | `GET .../locations/{loc}/environmentCurrencies` | `GET /environmentmanagement/provisioning/locations/{loc}/currencies` | available, fewer properties |
| List languages | `GET .../locations/{loc}/environmentLanguages` | `GET /environmentmanagement/provisioning/locations/{loc}/languages` | available, fewer properties |
| List templates | `GET .../locations/{loc}/environmentTemplates` | `GET /environmentmanagement/provisioning/locations/{loc}/templates` | available (richer) |
| Managed env on/off | BAPI `governanceConfiguration.protectionLevel` PATCH | `POST .../governanceSetting/enableManaged` / `disableManaged` | toggle only; no extendedSettings |
| Admin mode on/off | BAPI `properties.states.runtime.id` PATCH | `environment-state/enable-environment` / `disable-environment` | **paths UNVERIFIED (docs 404)** |

### 10.2 Property mapping — environment read model

| BAPI JSON path | New API JSON path | Notes |
| --- | --- | --- |
| `name` | `id` | New API drops the ARM-style `id` and uses the bare GUID. |
| `id` (ARM resource path) | — | Removed. |
| `type` (`Microsoft.BusinessAppPlatform/scopes/environments`) | — | Removed. New `type` means SKU. |
| `location` | `geo` | Renamed. |
| `properties.azureRegion` | `azureRegion` | Flattened. |
| `properties.displayName` | `displayName` | Flattened. |
| `properties.description` | — | **Removed.** |
| `properties.environmentSku` | `type` | Renamed and flattened. |
| `properties.tenantId` | `tenantId` | Flattened. |
| `properties.databaseType` | — | **Removed** (inferable from `dataverseId != null`). |
| `properties.provisioningState` | — | **Removed** (moved to the operations API). |
| `properties.createdTime` | `createdDateTime` | Renamed. |
| `properties.createdBy{id,displayName,type,tenantId}` | `createdBy{id,type}` | **Loses `displayName` and `tenantId`.** |
| — | `createdFor{id,type}` | **NEW.** |
| `properties.lastModifiedTime` | — | **Removed.** |
| `properties.creationType` | — | **Removed.** |
| `properties.isDefault` | — | **Removed.** |
| `properties.states.management.id` | `state` | Flattened/renamed. |
| `properties.states.runtime.id` | `adminMode` | Approximate mapping. |
| `properties.states.disasterRecovery.id` | — | **Removed** (moved to the failover API). |
| `properties.linkedEnvironmentMetadata.instanceUrl` | `url` | Flattened. |
| `properties.linkedEnvironmentMetadata.instanceApiUrl` | — | **Removed.** |
| `properties.linkedEnvironmentMetadata.domainName` | `domainName` | Flattened. |
| `properties.linkedEnvironmentMetadata.version` | `version` | Flattened. |
| `properties.linkedEnvironmentMetadata.resourceId` | `dataverseId` | Renamed and flattened. |
| `properties.linkedEnvironmentMetadata.securityGroupId` | `securityGroupId` | Flattened. |
| `properties.linkedEnvironmentMetadata.baseLanguage` | — | **Removed** (write-only via provisioning/reset). |
| `properties.linkedEnvironmentMetadata.uniqueName` | — | **Removed.** |
| `properties.linkedEnvironmentMetadata.friendlyName` | — | **Removed.** |
| `properties.linkedEnvironmentMetadata.template[]` | — | **Removed.** |
| `properties.linkedEnvironmentMetadata.templateMetadata` | — | **Removed.** |
| `properties.linkedEnvironmentMetadata.backgroundOperationsState` | `backgroundOperationsState` | Flattened. |
| `properties.linkedEnvironmentMetadata.instanceState` | — | **Removed** (approx. covered by `state`). |
| `properties.linkedEnvironmentMetadata.scaleGroup` | — | **Removed.** |
| `properties.linkedEnvironmentMetadata.platformSku` | — | **Removed.** |
| `properties.linkedEnvironmentMetadata.createdTime` | — | **Removed.** |
| `properties.governanceConfiguration.protectionLevel` | `protectionLevel` | Flattened. |
| `properties.governanceConfiguration.settings.extendedSettings.*` (14+ keys) | — | **Removed. MAJOR GAP.** |
| `properties.billingPolicy.id` | — | **Removed** (use `licensing/billingPolicies/{id}/environments`). |
| `properties.updateCadence.id` | — | **Removed.** |
| `properties.parentEnvironmentGroup.id` | `environmentGroupId` | Renamed and flattened. |
| `properties.cluster.category` | `clusterCategory` | Flattened. |
| `properties.cluster.number` | — | **Removed.** |
| `properties.usedBy{id,type,tenantId}` | — | **Removed** from read model (still in create request). |
| `properties.bingChatEnabled` | — | **Removed.** |
| `properties.m365Enabled` | — | **Removed.** |
| `properties.copilotPolicies.crossGeoCopilotDataMovementEnabled` | — | **Removed.** |
| `properties.copilotPolicies.crossBoundaryCopilotDataMovementEnabled` | — | **Removed.** |
| `properties.enterprisePolicies.vnets` | `enterprisePolicies.networkInjection` | Renamed. |
| `properties.enterprisePolicies.customerManagedKeys` | `enterprisePolicies.encryption` | Renamed. |
| `properties.enterprisePolicies.identity` | `enterprisePolicies.identity` | Same. |
| — | `enterprisePolicies.privateEndpoint` | **NEW.** |
| `properties.enterprisePolicies.*.policyId` | `enterprisePolicies.*.id` | Renamed. |
| `properties.enterprisePolicies.*.systemId` | `enterprisePolicies.*.resourceId` | Approximate. |
| `properties.enterprisePolicies.*.linkStatus` | `enterprisePolicies.*.status` | Renamed; enum now published (6 values). |
| `properties.enterprisePolicies.*.location` | — | **Removed.** |
| — | `enterprisePolicies.*.error` | **NEW.** |
| `properties.retentionPeriod` | `retentionDetails.retentionPeriod` | Consolidated. |
| `properties.retentionDetails.backupsAvailableFromDateTime` | `retentionDetails.availableFromDateTime` | Renamed. |
| `properties.protectionStatus.keyManagedBy` | — | **Removed.** |
| `properties.capacity[]` | — | **Removed** (use `licensing/*`). |
| `properties.addons[]` | — | **Removed.** |
| `properties.clientUris{admin,maker}` | — | **Removed.** |
| `properties.runtimeEndpoints{...}` (7 keys) | — | **Removed.** |
| `properties.notificationMetadata{state,branding}` | — | **Removed.** |
| `properties.connectedGroups[]` | `connectedGroupId` (singular string) | Shape change. |
| `properties.linkedAppMetadata{id,type,url}` | `finOpsMetadata{id,type,url}` | Renamed (FinOps = Dynamics 365 F&O linked app). |
| — | `deletedDateTime` | **NEW.** |
| — | `scenarioName` | **NEW.** |

### 10.3 Property mapping — locations

| BAPI | New API | Notes |
| --- | --- | --- |
| `name` | `collection[].name` | |
| `properties.displayName` | `collection[].displayName` | |
| `properties.code` | `collection[].code` | |
| `properties.isDefault` | `collection[].isDefault` | |
| `properties.isDisabled` | `collection[].isDisabled` | |
| `properties.canProvisionDatabase` | `collection[].canProvisionDatabase` | |
| `properties.canProvisionCustomerEngagementDatabase` | — | **Removed** (breaking for `powerplatform_locations`). |
| `properties.azureRegions[]` | — | **Removed** (breaking). |
| `id`, `type` | — | Removed. |
| — | `collection[].hasFirstReleaseIslandAvailableForProvisioning` | **NEW.** |
| — | `locationSelectionMode` (`Region` or `MacroRegion`) | **NEW.** |
| — | `macroRegions[]{macroRegionId, displayName, dataResidencyNote}` | **NEW.** |

### 10.4 Summary — new API only (not in BAPI)

* `createdFor{id,type}` — who the environment was created *for* (vs. *by*).
* `deletedDateTime` — soft-delete timestamp on the environment itself.
* `scenarioName` — "singleton scenario type".
* `enterprisePolicies.privateEndpoint` — fourth enterprise-policy slot.
* `enterprisePolicies.*.error` — failure detail on a policy link.
* `EnterprisePolicyLinkStatus` published as a real enum (`Linking`, `Unlinking`, `Linked`, `Failed`, `LinkingOnline`, `UnlinkingOnline`).
* Locations: `hasFirstReleaseIslandAvailableForProvisioning`, `locationSelectionMode`, `macroRegions[]` (with `dataResidencyNote`).
* Provisioning request: `macroRegion` (mutually exclusive with `location`), `connectedGroupIdForTeamsEnvironment`, `finOpsMetadata`, `currency.name`/`.symbol`/`.precision` (BAPI create only accepted `currency.code`).
* Templates: per-SKU `availability[]` with `disabledReason{code,message}`, `isSupportedForResetOperation`, `isCustomerEngagement`.
* `ValidateOnly` / `ValidateProperties` on every write operation.
* `errorDetail.fieldErrors[].suggestedValue` — server-proposed valid alternative.
* `stageStatuses[]` — per-stage operation progress (BAPI `lifecycleOperations` only exposes a single `provisioningState`).
* `EnvironmentSku` expanded to 12 documented values.
* Standard OData paging on environments (`$top`/`$skip`/`$orderby`/`@odata.nextlink`) and `$filter` on `dataverseId`, `type`, `geo`, `state`, `environmentGroupId`, `domainName`.
* `Copy Environment`, `Reset Environment`, `Restore Environment`, environment backup, failover/DR — all first-class endpoints with no BAPI equivalent surfaced in this provider.
* Single global polling host (no `https://{geo}.api.bap.microsoft.com` geo-routing).

### 10.5 Summary — BAPI only (removed in the new API)

Blocking or high-impact:

1. **`PATCH` update of an environment** — no equivalent at all.
2. **`governanceConfiguration.settings.extendedSettings`** (all managed-environment settings).
3. `description`, `lastModifiedTime`, `creationType`, `isDefault`, `provisioningState`, `databaseType`.
4. `updateCadence.id` (Frequent/Moderate release cadence).
5. `billingPolicy.id` on the environment object.
6. `bingChatEnabled`, `m365Enabled`, `copilotPolicies.*`.
7. `linkedEnvironmentMetadata.uniqueName`, `.friendlyName`, `.baseLanguage`, `.instanceApiUrl`, `.template[]`, `.templateMetadata`, `.instanceState`, `.scaleGroup`, `.platformSku`.
8. `capacity[]`, `addons[]`, `clientUris`, `runtimeEndpoints`, `notificationMetadata`, `protectionStatus.keyManagedBy`, `cluster.number`, `states.disasterRecovery`, `usedBy` (read side), `enterprisePolicies.*.location`.
9. Locations: `canProvisionCustomerEngagementDatabase`, `azureRegions[]`.
10. Currencies: `name`, `isDisabled`.
11. Languages: `displayName`, `isDisabled`.

---

## 11. Gaps and UNVERIFIED items

| # | Item | Impact | How to verify |
| --- | --- | --- | --- |
| 1 | 202 response headers for provisioning/delete/recover/reset/copy (`Location` vs `Operation-Location`), and `Retry-After`. | **Blocking** — cannot write the polling loop or httpmock responders. | mitmproxy capture (devdocs/adr/mitmproxy.md) against a real tenant, or inspect `Microsoft.PowerPlatform.Management`. |
| 2 | No Update Environment endpoint. | **Blocking** for full migration. | Confirm with the Power Platform API team / release plans. |
| 3 | Managed-environment `extendedSettings` have no new-API surface. | **Blocking** for managed environment support. | Check the `governance` namespace rule sets and the Power Platform for Admins V2 connector. |
| 4 | `environment-state/enable-environment` and `disable-environment` doc pages 404. | High — admin mode toggling. | Try the `Microsoft.PowerPlatform.Management` SDK; retry Learn later. |
| 5 | Enum value sets for `state`, `adminMode`, `backgroundOperationsState`, `protectionLevel`, `clusterCategory`, `type` on `EnvironmentResponse` are undocumented (bare `string`). | Medium — validators/plan modifiers. | Live capture across environment types. |
| 6 | `templateMetadata` shape (typed as bare `object`). | Medium — post-provisioning packages. | Live capture; compare with BAPI `createTemplateMetadataDto`. |
| 7 | Whether `Provision New Environment` supports `ValidateOnly`/`ValidateProperties` (not listed in its URI-parameter table, unlike every other write op). | Medium — replaces `validateEnvironmentDetails`. | Live probe. |
| 8 | Query-string casing of `ValidateOnly` / `ValidateProperties`. | Low. | Live probe. |
| 9 | Granular Entra permission (or admin role) required for create/delete/recover/reset/copy. | **High** — service-principal auth for the provider. | Test with an SP holding only `EnvironmentManagement.Environments.Read`. |
| 10 | Whether `location` accepts BAPI-style names (`unitedstates`, `europe`) — i.e. whether `Location.name` from Get Supported Locations is what `CreateEnvironmentRequest.location` expects. | Medium. | Live probe. |
| 11 | Exact paths for Restore / Backup / Copy-candidates / Failover operation groups. | Low (out of current provider scope). | Fetch remaining Learn pages. |
| 12 | Environment Groups Create/List/Update/Delete/Remove exact paths and bodies. | Low (separate service in the provider). | Fetch remaining Learn pages. |
| 13 | No published example request/response payloads on ANY `environmentmanagement` Learn page — every JSON block in this document is reconstructed from the schema tables, not copied verbatim. | Medium — risk of casing/shape mistakes. | Live capture. |

---

## Source Links

Power Platform API reference:

* https://learn.microsoft.com/en-us/rest/api/power-platform/
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/provision-new-environment
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/link-dataverse
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-supported-locations
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-currencies
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-languages
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-provisioning/get-provisioning-templates
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/list-environments-for-user
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environments/get-environment-by-id-for-user
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-delete/delete-environment-by-id
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-recover/recover-environment
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-reset/reset-environment
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-copy/copy-environment
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-managed-governance/enable-managed-environment
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-managed-governance/disable-managed-environment
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-operation-by-id
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-environment-operation-by-id
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/operation/get-operations-for-environment
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/list-environment-management-settings
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-management-settings/update-environment-management-settings
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-groups
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-groups/get-environment-group
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-groups/add-environment-to-group
* https://learn.microsoft.com/en-us/rest/api/power-platform/environmentmanagement/environment-groups/get-environment-group-operation

Conceptual and policy:

* https://learn.microsoft.com/en-us/power-platform/admin/programmability-versioning-support
* https://learn.microsoft.com/en-us/power-platform/admin/programmability-permission-reference
* https://learn.microsoft.com/en-us/power-platform/admin/programmability-whats-new-changed
* https://learn.microsoft.com/en-us/power-platform/admin/programmability-authentication-v2
* https://learn.microsoft.com/en-us/power-platform/admin/programmability-tutorial-environmentmanagement-settings

Legacy BAPI:

* https://learn.microsoft.com/en-us/power-platform/admin/list-environments

SDKs:

* https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management?view=power-platform-latest
* https://learn.microsoft.com/en-us/dotnet/api/microsoft.powerplatform.management.environmentmanagement?view=power-platform-latest
* https://www.nuget.org/packages/Microsoft.PowerPlatform.Management
* https://pypi.org/project/powerplatform-management/
* https://learn.microsoft.com/en-us/connectors/powerplatformadminv2/

Workspace files consulted:

* internal/constants/constants.go
* internal/constants/constants_test.go
* internal/services/environment/dto.go
* internal/services/environment/api_environment.go
* internal/services/environment/resource_environment_test.go
* internal/services/environment/datasource_environments_test.go
* internal/services/application/resource_environment_application_package_install_test.go
* .copilot-tracking/research/2026-08-18/powerplatform-api-migration-research.md

## Recommended Next Research

* [ ] Live mitmproxy capture of `POST /environmentmanagement/provisioning/environments` to confirm 202 headers, polling URL shape, and `Retry-After`.
* [ ] Live capture of `GET /environmentmanagement/environments` to confirm exact JSON casing, null handling, and the enum value sets for `state`/`adminMode`/`backgroundOperationsState`/`protectionLevel`/`clusterCategory`.
* [ ] Inspect the `Microsoft.PowerPlatform.Management` NuGet package to recover the missing `environment-state` enable/disable paths and any undocumented models.
* [ ] Fetch the remaining Learn pages: `environment-restore/*`, `environment-backup/*`, `environment-copy/get-environment-copy-candidates`, `failover/*`, `environment-groups` create/list/update/delete/remove, `environment-management-settings/create-environment-management-settings`.
* [ ] Determine the required Entra permission or admin role for provisioning-side operations (test with a least-privilege service principal).
* [ ] Investigate whether managed-environment `extendedSettings` are reachable via the `governance` namespace (rule sets / rule-based policies) rather than `environmentmanagement`.
* [ ] Confirm whether an "Update Environment" endpoint is on the published Power Platform release plan.
* [ ] Verify the `templateMetadata` / `PostProvisioningPackages` wire shape on the new API.
* [ ] Assess non-public cloud coverage (GCC/GCC High/DoD/China) — does `api.powerplatform.com` have sovereign-cloud equivalents for `environmentmanagement`? Compare with the cloud tables in internal/constants/constants.go.
