<!-- markdownlint-disable-file -->
# Subagent Research: Live verification mechanics — probing `api.powerplatform.com` and `api.bap.microsoft.com`

Research topic: determine exactly HOW a developer on this repository can execute live/authenticated probe calls against `https://api.powerplatform.com` and `https://api.bap.microsoft.com` to verify open questions about the Environment Management API, and how to capture the raw responses for later reuse as `httpmock` fixtures.

Status: **COMPLETE** for all 7 questions. Two blockers identified (no `az login` session in this container; no documented test tenant). Everything else is verified against the repository and against the live devcontainer toolchain.

Research date: 2026-08-19.

Prior research consumed:

* .copilot-tracking/research/subagents/2026-08-18/get-environment-response-completeness.md (the open questions this probe procedure is meant to resolve)
* .copilot-tracking/research/subagents/2026-08-18/powerplatform-api-environmentmanagement.md (target API surface, `api-version=2024-10-01`)
* .copilot-tracking/research/subagents/2026-08-18/current-environment-implementation.md (as-is BAPI implementation)

---

## 1. mitmproxy + httpmock workflow (question 1)

### 1.1 What the ADRs actually say

* devdocs/adr/mitmproxy.md — **Status: Proposed.** It is a decision record only. It contains **no commands, no env vars, no capture paths**. Key statements: mitmproxy is integrated into the devcontainer; the recorded traffic "has the potential to serve as a foundation for generating mock responses"; and explicitly "No changes have been made to the provider code to support mitmproxy, ensuring that there is no test-only code path."
  * **Conclusion: there is no automated capture-to-fixture pipeline.** Conversion of captures into JSON fixtures is a manual copy/paste today. The ADR lists that automation as future work.
* devdocs/adr/httpmocks.md — **Status: Accepted.** Also decision-only. It mandates `github.com/jarcoal/httpmock`, activate at test start / deactivate at test end, and "Create Acceptance Tests that mirror unit tests to detect changes in external APIs." No capture workflow.

### 1.2 The actual commands (assembled from Makefile, DEVELOPER.md, and the devcontainer feature)

Three mitmproxy binaries are installed at `/usr/local/bin/` (verified present in this container): `mitmproxy` (TUI), `mitmweb` (web UI), `mitmdump` (headless).

Installed by .devcontainer/features/local_provider_dev/install.sh lines 29-42:

* MITM_VERSION = `11.1.3` (line 30) — verified: `Mitmproxy: 11.1.3 binary`, Python 3.13.2.
* Line 38: `timeout 5s mitmdump -p 8080 || true` — generates the CA on image build.
* Line 40: `install -D ~/.mitmproxy/mitmproxy-ca-cert.pem /usr/local/share/ca-certificates/mitmproxy-ca.crt`
* Line 42: `cat ~/.mitmproxy/mitmproxy-ca-cert.pem >> /etc/ssl/certs/ca-certificates.crt`

  **This is the critical bit**: the mitmproxy CA is appended to the system CA bundle, so Go's `crypto/x509` system pool trusts it. TLS interception of the provider's own HTTPS calls works with **no** `InsecureSkipVerify` and no provider code change.

Proxy port is **8080** everywhere.

| Purpose | Exact command | Source |
| --- | --- | --- |
| Interactive TUI capture | `mitmproxy` (defaults to port 8080) | DEVELOPER.md "Debugging network calls" |
| Web UI capture | `mitmweb` | DEVELOPER.md |
| Headless capture to a flow file | `make netdump` → `mitmdump -p 8080 -w /tmp/mitmproxy.dump` | Makefile `netdump` target |
| Run acceptance tests through the proxy | `make acctest TEST=<prefix> USE_PROXY=1` | Makefile `acctest` target |
| Run a `terraform apply` example through the proxy | `HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 terraform apply` | DEVELOPER.md |

Makefile `acctest` target, `USE_PROXY=1` branch (exact line):

```makefile
@HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
```

`USE_PROXY` is **only** consumed by the Makefile; it is not read anywhere in Go code. `HTTP_PROXY` / `HTTPS_PROXY` are honoured because internal/api/request.go line 43 uses `http.DefaultClient`, whose transport is `http.DefaultTransport` with `Proxy: http.ProxyFromEnvironment`.

### 1.3 Where captures land

* `make netdump` writes the mitmproxy binary flow file to **`/tmp/mitmproxy.dump`**. Nothing else in the repo reads or post-processes it.
* There is **no** `.gitignore`d capture directory and **no** conversion script anywhere (scripts/ contains only `user_story_prompt.sh`).

### 1.4 Turning a capture into a fixture (manual, but with useful mitmdump flags)

Verified available options on mitmproxy 11.1.3 (`mitmdump --options`):

* `-w PATH` / `--save-stream-file PATH` — stream flows to a binary flow file as they arrive (prefix path with `+` to append).
* `-r PATH` / `--rfile PATH` — read flows back from a file.
* `--set hardump=<file.har>` — **"Save a HAR file with all flows on exit."** This is the cleanest machine-readable export; a HAR entry's `response.content.text` is the raw JSON body, extractable with `jq`.
* `--set flow_detail=<0..4>` — verbosity; `3` prints full request/response bodies to stdout.
* `--set save_stream_filter=<filter>` — restrict what is written.
* Positional `filter_args` — a mitmproxy filter expression, e.g. `"~u environmentmanagement"` or `"~d api.powerplatform.com"`.

Practical replay/export commands:

```bash
# capture only Power Platform API + BAPI traffic straight to a HAR
mitmdump -p 8080 -w /tmp/ppapi.dump --set hardump=/tmp/ppapi.har "~d api.powerplatform.com | ~d api.bap.microsoft.com"

# offline: re-read a flow file and print full bodies
mitmdump -n -r /tmp/mitmproxy.dump --set flow_detail=3 "~u environmentmanagement"

# offline: convert an existing flow file to HAR without a proxy listening
mitmdump -n -r /tmp/mitmproxy.dump --set hardump=/tmp/ppapi.har

# pull a single response body out of the HAR
jq -r '.log.entries[] | select(.request.url | test("environmentmanagement/environments/")) | .response.content.text' /tmp/ppapi.har | jq . > get_environment.json
```

### 1.5 Fixture placement rules (from .github/copilot-instructions.md and devdocs/testing_guidelines.md)

* Path: `internal/services/<service>/tests/resource/<Test_Scenario>/` or `.../tests/datasource/<Test_Scenario>/`.
  * Note the divergence: devdocs/testing_guidelines.md says `internal/services/<service>/test/...` (singular), but the **actual on-disk convention is `tests/`** (plural) — e.g. internal/services/environment/tests/datasource/Validate_Read/. Follow the on-disk convention.
* Scenario folder name = the test function name with the `TestUnit<Something>` prefix stripped (e.g. `TestUnitEnvironmentsDataSource_Validate_Read` → `Validate_Read`).
* File name = `<method>_<object>.json`, lowercase, no spaces, optional numeric suffix when a scenario issues the same call more than once (`get_environment_1.json`, `get_environment_2.json`, ...).
* Fixtures are loaded in tests with `httpmock.File("tests/resource/<Scenario>/<file>.json").String()` (path is relative to the service package directory).
* **Anonymize everything**: real GUIDs → `00000000-0000-0000-0000-00000000000N`, no tenant IDs, no emails, no PII. Existing fixtures follow this strictly.

---

## 2. Acceptance tests, `TF_ACC`, `TF_LOG`, and auth env vars (question 2)

### 2.1 Makefile targets (verbatim)

```makefile
unittest: clean install
	@TF_ACC=0 go test -p 16 -timeout 10m -v -cover ./... -run "^TestUnit$(TEST)"

acctest: clean install
ifeq ($(USE_PROXY),1)
	@HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
else
	@TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
endif

test: clean install
	@TF_ACC=1 go test -p 10 -timeout 300m -v -cover ./...

netdump:
	@mitmdump -p 8080 -w /tmp/mitmproxy.dump
```

* `TF_ACC=1` is set by the Makefile — do **not** set it manually.
* `TEST=<prefix>` is interpolated into the regex `^TestAcc<prefix>` (and `^TestUnit<prefix>` for unit tests). It is a **prefix**, not a substring, so `TEST=EnvironmentsResource_Validate_Create` matches `TestAccEnvironmentsResource_Validate_Create` and everything starting with it.
* `acctest` has a 300 minute timeout and `-p 10` package parallelism.
* Both targets run `clean install` first — `go clean -testcache`, `rm -rf ./bin`, `rm -rf /go/bin/terraform-provider-power-platform`, `go mod tidy`, `go fmt ./...`, `go install`.

### 2.2 `TF_LOG` / `TF_LOG_PROVIDER`

* devdocs/observability_guidelines.md line 5: "The provider is configured so that when a user sets `TF_LOG` or `TF_LOG_PROVIDER` environment variables, the logs from `tflog` calls will appear."
* The devcontainer **pre-sets `TF_LOG=ERROR`** (.devcontainer/devcontainer.json line 42) — verified live in this container's environment. **You must override it** (`TF_LOG=DEBUG` or `TF_LOG_PROVIDER=DEBUG`) to see any provider `tflog.Debug` output.
* `TF_ACC` and `TF_LOG` are independent: `TF_ACC` gates `resource.Test` execution (skips when unset/0 for non-`IsUnitTest` cases); `TF_LOG` only controls log verbosity.

### 2.3 The complete list of env var names the provider reads

All defined in internal/constants/constants.go lines 160-187 and consumed in internal/provider/provider.go lines 219-259.

| Constant (constants.go line) | Env var name | Provider attribute | Consumed at provider.go |
| --- | --- | --- | --- |
| 160 | `POWER_PLATFORM_CLOUD` | `cloud` (default `public`) | 219 |
| 161 | `POWER_PLATFORM_TENANT_ID` | `tenant_id` | 220 |
| 162 | `POWER_PLATFORM_CLIENT_ID` | `client_id` | 222 |
| 163 | `POWER_PLATFORM_AUXILIARY_TENANT_IDS` | `auxiliary_tenant_ids` | 221 |
| 164 | `POWER_PLATFORM_CLIENT_SECRET` | `client_secret` | 223 |
| 165 | `POWER_PLATFORM_USE_OIDC` | `use_oidc` | 224 |
| 166 | `POWER_PLATFORM_USE_CLI` | `use_cli` | 225 |
| 167 | `POWER_PLATFORM_USE_DEV_CLI` | `use_dev_cli` (azd) | 226 |
| 168 | `POWER_PLATFORM_USE_MSI` | `use_msi` | 230 |
| 169 | `POWER_PLATFORM_CLIENT_CERTIFICATE` | `client_certificate` (base64 PFX) | 227 |
| 170 | `POWER_PLATFORM_CLIENT_CERTIFICATE_FILE_PATH` | `client_certificate_file_path` | 228 |
| 171 | `POWER_PLATFORM_CLIENT_CERTIFICATE_PASSWORD` | `client_certificate_password` | 229 |
| 172 | `POWER_PLATFORM_TELEMETRY_OPTOUT` | `telemetry_optout` | 240 |
| 173 | `POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID` | `azdo_service_connection_id` | 231 |
| 174 | `POWER_PLATFORM_ENABLE_CAE` | `enable_continuous_access_evaluation` | 259 |
| 176 | `ARM_OIDC_REQUEST_URL` | `oidc_request_url` (first of pair) | 234 |
| 177 | `ACTIONS_ID_TOKEN_REQUEST_URL` | `oidc_request_url` (fallback) | 234 |
| 178 | `ARM_OIDC_REQUEST_TOKEN` | `oidc_request_token` (first of pair) | 235 |
| 179 | `ACTIONS_ID_TOKEN_REQUEST_TOKEN` | `oidc_request_token` (fallback) | 235 |
| 180 | `ARM_OIDC_TOKEN` | `oidc_token` | 236 |
| 181 | `ARM_OIDC_TOKEN_FILE_PATH` | `oidc_token_file_path` | 237 |
| 182 | `ARM_AUXILIARY_TENANT_IDS` | `auxiliary_tenant_ids` (fallback) | 221 |
| 184 | `POWER_PLATFORM_PARTNER_ID` | `partner_id` | 242 |
| 185 | `ARM_PARTNER_ID` | `partner_id` (fallback) | 242 |
| 186 | `POWER_PLATFORM_DISABLE_TERRAFORM_PARTNER_ID` | `disable_terraform_partner_id` | 253 |
| 187 | `ARM_DISABLE_TERRAFORM_PARTNER_ID` | fallback for the above | 255 |

Test-only env var (not provider config): `ACCEPTANCE_TESTS_LICENSING_GROUP_NAME` — internal/mocks/mocks.go line 31, defaults to `"pptestusers"`.

Minimum viable auth sets (dispatch order in internal/api/auth.go lines 505-522, `GetTokenForScopes`):

1. **Client secret** — `POWER_PLATFORM_TENANT_ID` + `POWER_PLATFORM_CLIENT_ID` + `POWER_PLATFORM_CLIENT_SECRET`
2. **Azure CLI** — `POWER_PLATFORM_USE_CLI=true` (+ optional `POWER_PLATFORM_TENANT_ID` to pin the tenant) and a prior `az login`
3. **Azure Developer CLI** — `POWER_PLATFORM_USE_DEV_CLI=true`
4. **AzDO workload identity federation** — `POWER_PLATFORM_USE_OIDC=true` + `POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID`
5. **OIDC (GitHub Actions)** — `POWER_PLATFORM_USE_OIDC=true` + `ARM_OIDC_REQUEST_URL`/`ACTIONS_ID_TOKEN_REQUEST_URL` + `ARM_OIDC_REQUEST_TOKEN`/`ACTIONS_ID_TOKEN_REQUEST_TOKEN` (or `ARM_OIDC_TOKEN` / `ARM_OIDC_TOKEN_FILE_PATH`) + tenant/client id
6. **Client certificate** — tenant/client id + `POWER_PLATFORM_CLIENT_CERTIFICATE` or `..._FILE_PATH` (+ `..._PASSWORD`)
7. **MSI** — `POWER_PLATFORM_USE_MSI=true` (+ `POWER_PLATFORM_CLIENT_ID` for user-assigned)

If none match, `GetTokenForScopes` returns `errors.New("no credentials provided")` (internal/api/auth.go line 521).

### 2.4 `testAccPreCheck` does NOT exist

devdocs/testing_guidelines.md line 58 claims "The test suite includes **pre-checks** (`testAccPreCheck(t)`) to validate required environment variables before running tests." A repo-wide search for `PreCheck` finds **zero** `PreCheck:` fields on any `resource.TestCase` and no `testAccPreCheck` function. **The documentation is stale.** Acceptance tests fail at provider-configure time with a diagnostic (`validateProviderAttribute`, internal/provider/provider.go lines 482-494) rather than skipping. Plan for a hard failure, not a skip.

### 2.5 What the devcontainer already sets

.devcontainer/devcontainer.json `containerEnv` (verified live in this shell):

| Line | Variable | Value |
| --- | --- | --- |
| 42 | `TF_LOG` | `ERROR` |
| 44 | `POWER_PLATFORM_USE_CLI` | `true` |
| 45 | `ARM_USE_CLI` | `true` |
| — | `TF_CLI_CONFIG_FILE` | points at the dev-override tfrc |

So the container is **pre-wired for Azure CLI user auth**; the only missing piece is an `az login`.

---

## 3. Existing mechanism to dump raw HTTP bodies (question 3)

**There is none.** Verified by reading the whole request path.

| Fact | Location |
| --- | --- |
| Every `tflog` call in `internal/api` — 8 in auth.go (106, 112, 117, 184, 210, 423, 493, 528), 1 in client.go (198), 4 in lifecycle.go (58, 78, 81, 94), 3 in request.go (99, 130, 152) | internal/api/*.go |
| **None** of them logs a request or response body | — |
| The response body is read into memory here and never logged | internal/api/request.go line 84 (`resp.BodyAsBytes = body`) |
| Token is explicitly masked in logs: `"Token acquired (expire: %s): **********"` | internal/api/auth.go line 528 |
| Only place a body reaches the user is inside an **error**: `customerrors.NewUnexpectedHttpStatusCodeError(..., resp.BodyAsBytes)` — i.e. only on an unacceptable status code | internal/api/client.go line 193 |
| No `logging.NewLoggingHTTPTransport`, no `httputil.DumpRequest`/`DumpResponse`, no custom `http.RoundTripper` anywhere | repo-wide search: 0 hits |
| The HTTP client is the bare `http.DefaultClient` (no wrapped transport) | internal/api/request.go line 43 |
| `TF_LOG_PROVIDER` is mentioned only in prose | devdocs/observability_guidelines.md line 5 |
| ADR explicitly states no test-only code path was added for capture | devdocs/adr/mitmproxy.md, Consequences |

**Consequence: `TF_LOG=TRACE` will NOT give you response bodies.** The only way to see raw bodies from provider-originated traffic is mitmproxy (which works because `http.DefaultTransport` honours `HTTPS_PROXY` and the mitm CA is in the system trust store).

**One exception worth knowing**: a non-2xx response *is* surfaced verbatim in the Terraform error message via `NewUnexpectedHttpStatusCodeError`. A deliberately-wrong `acceptableStatusCodes` list in a throwaway probe would echo a 200 body into the error — hacky, but it works without mitmproxy.

---

## 4. Existing standalone probe tool, and the minimal idiomatic way to build one (question 4)

### 4.1 Nothing exists

| Searched | Result |
| --- | --- |
| scripts/ | only `user_story_prompt.sh` (a `gh`/prompt helper — unrelated) |
| tools.go | build-tag `tools` only; single blank import of `tfplugindocs`. No probe tooling. |
| internal/api | `Client.Execute` / `ExecuteWithoutRetry` only; no CLI entrypoint |
| main.go | the provider `plugin.Serve` entrypoint only |
| Makefile | no probe/curl/token target |

There is **no** existing way to fire a one-off authenticated request at an arbitrary Power Platform URL.

### 4.2 Minimal idiomatic Go probe using the existing packages

The building blocks (all exported):

| Symbol | File:line | Signature |
| --- | --- | --- |
| `config.ProviderConfig` | internal/config/config.go:31-63 | struct; set `UseCli`, `TenantId`, `Urls`, `Cloud`, `TestMode: false` |
| `config.ProviderConfigUrls` | internal/config/config.go:65-76 | `BapiUrl`, `PowerPlatformUrl`, `PowerAppsScope`, `PowerPlatformScope`, ... |
| `api.NewAuthBase` | internal/api/auth.go:87 | `func NewAuthBase(*config.ProviderConfig) *Auth` |
| `api.NewApiClientBase` | internal/api/client.go:45 | `func NewApiClientBase(*config.ProviderConfig, *Auth) *Client` |
| `(*api.Client).Execute` | internal/api/client.go:114 | `(ctx, scopes []string, method, url string, headers http.Header, body any, acceptableStatusCodes []int, responseObj any) (*api.Response, error)` |
| `(*api.Auth).GetTokenForScopes` | internal/api/auth.go:492 | `(ctx, scopes []string) (*string, error)` |
| `tryGetScopeFromURL` | internal/api/client.go:253-270 | **unexported** — cannot be called from outside package `api` |

`tryGetScopeFromURL` behaviour (internal/api/client.go lines 253-270), which is what makes `scopes = nil` work:

```go
case url contains cloudConfig.BapiUrl || cloudConfig.PowerAppsUrl -> cloudConfig.PowerAppsScope
case url contains cloudConfig.PowerPlatformUrl               -> cloudConfig.PowerPlatformScope
case url contains cloudConfig.PowerAppsAdvisor              -> cloudConfig.PowerAppsAdvisorScope
case url contains cloudConfig.AdminPowerPlatformUrl         -> constants.PPAC_SCOPE
case url contains "csanalytics"                             -> cloudConfig.AnalyticsScope
default                                                     -> scheme://host/.default
```

Because `PUBLIC_POWERPLATFORM_API_DOMAIN = "api.powerplatform.com"` (constants.go:31) is a substring of the target URL, passing `nil` scopes auto-selects `PUBLIC_POWERPLATFORM_API_SCOPE = "https://api.powerplatform.com/.default"` (constants.go:32). Likewise `PUBLIC_BAPI_DOMAIN = "api.bap.microsoft.com"` (constants.go:28) auto-selects `PUBLIC_POWERAPPS_SCOPE = "https://service.powerapps.com/.default"` (constants.go:30) — **note the BAPI host uses the PowerApps scope, not a `bap` scope.**

Blocker for an out-of-package probe: `getCloudPublicUrls()` (internal/provider/provider.go:509-522) is **unexported**, so the URL/scope table must either be rebuilt by hand from `internal/constants` or the probe must live inside package `provider`.

**Recommended shape — a throwaway `TestAcc` in the `environment` package** (no new build targets, picked up by `make acctest`, deleted before commit):

```go
// internal/services/environment/probe_live_test.go   (temporary, do not commit)
package environment_test

func TestAccProbeEnvironmentApis(t *testing.T) {
    envID := os.Getenv("PROBE_ENVIRONMENT_ID")
    ctx := context.Background()

    cfg := &config.ProviderConfig{
        UseCli:   true,
        TenantId: os.Getenv("POWER_PLATFORM_TENANT_ID"),
        Cloud:    cloud.AzurePublic,
        Urls: config.ProviderConfigUrls{
            BapiUrl:            constants.PUBLIC_BAPI_DOMAIN,
            PowerAppsUrl:       constants.PUBLIC_POWERAPPS_API_DOMAIN,
            PowerAppsScope:     constants.PUBLIC_POWERAPPS_SCOPE,
            PowerPlatformUrl:   constants.PUBLIC_POWERPLATFORM_API_DOMAIN,
            PowerPlatformScope: constants.PUBLIC_POWERPLATFORM_API_SCOPE,
        },
    }
    client := api.NewApiClientBase(cfg, api.NewAuthBase(cfg))

    // new Power Platform API — nil scopes resolve to PUBLIC_POWERPLATFORM_API_SCOPE
    newURL := "https://api.powerplatform.com/environmentmanagement/environments/" + envID + "?api-version=2024-10-01"
    respNew, err := client.Execute(ctx, nil, "GET", newURL, nil, nil, []int{http.StatusOK}, nil)
    require.NoError(t, err)
    os.WriteFile("/tmp/probe_new_get_environment.json", respNew.BodyAsBytes, 0o600)

    // legacy BAPI — mirrors internal/services/environment/api_environment.go:258-263
    bapiURL := "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/" + envID +
        "?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01"
    respOld, err := client.Execute(ctx, nil, "GET", bapiURL, nil, nil, []int{http.StatusOK}, nil)
    require.NoError(t, err)
    os.WriteFile("/tmp/probe_bapi_get_environment.json", respOld.BodyAsBytes, 0o600)
}
```

Run with `make acctest TEST=ProbeEnvironmentApis` (the Makefile supplies `TF_ACC=1`; `resource.Test` is not used here so `TF_ACC` is irrelevant, but the target regex `^TestAcc` requires the `TestAcc` prefix).

Notes:

* `resp.BodyAsBytes` (internal/api/request.go:105-108) is the raw, unparsed body — exactly what you want for a fixture.
* Pass `responseObj = nil` so nothing is unmarshalled and no field is silently dropped.
* Keep `TestMode: false` — `TestMode: true` short-circuits `GetTokenForScopes` to a mock token (internal/api/auth.go:495-498).
* Also probe the list sibling `GET https://api.powerplatform.com/environmentmanagement/environments?api-version=2024-10-01`, and `GET https://api.powerplatform.com/licensing/environments/{id}/billingPolicy?api-version=2024-10-01` for the `billing_policy_id` question.

### 4.3 Even smaller: token-only probe

```go
auth := api.NewAuthBase(cfg)
tok, err := auth.GetTokenForScopes(ctx, []string{constants.PUBLIC_POWERPLATFORM_API_SCOPE})
```

then hand the token to `curl`. This avoids the whole `Client` plumbing but still exercises the provider's exact credential chain.

---

## 5. Getting a token outside the provider (question 5)

### 5.1 Tool availability in this devcontainer (verified live)

| Tool | Path | Status |
| --- | --- | --- |
| `az` | `/usr/local/bin/az` | **present** (installed via pipx, .devcontainer/features/local_provider_dev/install.sh line 60) |
| `pac` (Power Platform CLI) | — | **MISSING** |
| `mitmproxy` / `mitmdump` / `mitmweb` | `/usr/local/bin/` | present (11.1.3) |
| `terraform` | `/home/runtimeuser/tfenv/bin/terraform` | present |
| `curl` | `/usr/bin/curl` | present |
| `jq` | `/usr/bin/jq` | present |
| `changie` | `/home/runtimeuser/go/bin/changie` | present |

`az account show` → **`ERROR: Please run 'az login' to setup account.`** — **no active Azure CLI session. This is blocker #1.**

### 5.2 Yes, a token can be obtained outside the provider

```bash
az login --allow-no-subscriptions                     # docs/guides/azure_cli.md
# optionally: az login --service-principal --username <CLIENT_ID> --password <CLIENT_SECRET> --tenant <TENANT_ID> --allow-no-subscriptions

TOKEN=$(az account get-access-token \
  --resource https://api.powerplatform.com \
  --query accessToken -o tsv)

curl -sS -H "Authorization: Bearer $TOKEN" \
  "https://api.powerplatform.com/environmentmanagement/environments/<ENV_ID>?api-version=2024-10-01" | jq .
```

Resource values line up exactly with the provider's scope constants (drop the `/.default` suffix for `az`'s `--resource`):

| Host | Provider scope constant (constants.go line) | `az --resource` value |
| --- | --- | --- |
| `api.powerplatform.com` | `PUBLIC_POWERPLATFORM_API_SCOPE` (32) = `https://api.powerplatform.com/.default` | `https://api.powerplatform.com` |
| `api.bap.microsoft.com` | `PUBLIC_POWERAPPS_SCOPE` (30) = `https://service.powerapps.com/.default` | `https://service.powerapps.com` |

```bash
BAPI_TOKEN=$(az account get-access-token --resource https://service.powerapps.com --query accessToken -o tsv)
curl -sS -H "Authorization: Bearer $BAPI_TOKEN" \
  'https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/<ENV_ID>?$expand=permissions,properties.capacity,properties/billingPolicy,properties/copilotPolicies&api-version=2023-06-01' | jq .
```

### 5.3 What the repo documents about this

* docs/guides/azure_cli.md — `az login --allow-no-subscriptions`, then `use_cli = true`. Covers both user and service-principal login. Does **not** mention `az account get-access-token`.
* docs/guides/app_registration.md line 49 — "You will need to preauthorize Azure CLI to access your API by adding client application `04b07795-8ddb-461a-bbee-02f9e1bf7b46`". **This is required for `az account get-access-token --resource https://api.powerplatform.com` to succeed.**
* docs/guides/app_registration.md lines 39-45 — if the Power Platform API doesn't appear in the tenant, create its service principal: app id **`8578e004-a5c6-46e7-913e-12f58912df43`** ("Power Platform API"). Without this SP the token request fails with an AAD "resource principal not found" error.
* docs/guides/app_registration.md lines 22-32 — required delegated permissions, including `EnvironmentManagement.Environments.Read`.
* DEVELOPER.md "Power Platform Prerequisites" — points at the external [power-platform-terraform-quickstarts bootstrap](https://github.com/microsoft/power-platform-terraform-quickstarts/blob/main/bootstrap/README.md) for tenant setup. **No test tenant id, no shared credentials are documented in this repo. This is blocker #2.**
* `pac` CLI is **not** installed and **not** referenced anywhere in the repo — not a viable path here.

### 5.4 curl through mitmproxy

`curl` does not trust the mitm CA via `/etc/ssl/certs/ca-certificates.crt`? It does — the feature script appends the CA to that bundle, which is curl's default on this distro. So `curl -x http://127.0.0.1:8080 ...` also works and lets you capture the probe itself into the same flow file as any provider traffic.

---

## 6. Existing `powerplatform_environment` acceptance tests that provision real environments (question 7 in the brief, item 6)

File: internal/services/environment/resource_environment_test.go. All use `ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories` (the **real** provider, no `TestMode`), and none has a `PreCheck` or `CheckDestroy`.

| Line | Test function | What it provisions live |
| --- | --- | --- |
| 424 | `TestAccEnvironmentsResource_Validate_Update_Name_Field` | Sandbox + Dataverse (`unitedstates`, 1033/USD), then renames `aaa` → `aaa1`. Fixed names — **collision-prone.** |
| 463 | `TestAccEnvironmentsResource_Validate_CreateGenerativeAiFeatures_Non_US_Region_Update` | Sandbox + Dataverse, randomized name, toggles generative-AI flags |
| 546 | `..._CreateGenerativeAiFeatures_US_Region_Update` | Same, US region |
| 610 | `..._CreateGenerativeAiFeatures_US_Region_Expect_Fail` | Negative — expects an error, may not create |
| 632 | `..._CreateGenerativeAiFeatures_Non_US_Region_Expect_Fail` | Negative |
| 654 | `..._Non_US_Region_Microsoft_365_Services_Expect_Fail` | Negative |
| 677 | `..._Non_US_Region_Flex_Routing_Expect_Fail` | Negative |
| 700 | `TestAccEnvironmentsResource_Validate_CreateDeveloperEnvironment` | Developer-type environment |
| 966 | `TestAccEnvironmentsResource_Validate_Create_Early_Release_Cycle` | Sandbox with `release_cycle = "Early"` |
| 1050 | `TestAccEnvironmentsResource_Validate_Update` | Create then update |
| 1120 | `TestAccEnvironmentsResource_Validate_Domain_Uniqueness_On_Update` | Two environments, domain collision on update |
| **1176** | **`TestAccEnvironmentsResource_Validate_Create`** | **Single-step: Sandbox, `unitedstates`, Dataverse 1033/USD, randomized `domain = orgtestNNNNN`, name from `mocks.TestName()`. Asserts `dataverse.url`, `organization_id`, `unique_name`, `version`, `billing_policy_id`, `release_cycle`. — the richest single-step, most attribute-complete live create. BEST PIGGYBACK TARGET.** |
| 2633 | `TestAccEnvironmentsResource_Validate_Create_No_Dataverse` | Sandbox in `europe`, **no** Dataverse. Useful for the null-Dataverse shape. |
| 2663 | `TestAccEnvironmentsResource_Validate_Create_Them_Try_Remove_Dataverse` | Create with Dataverse, then attempt removal |
| 2914 | `TestAccEnvironmentsResource_Validate_Create_Environment_And_Dataverse` | Step 1 no Dataverse → step 2 adds Dataverse. **Captures both shapes for the same environment id — excellent for a before/after diff.** |
| 2969 | `TestAccEnvironmentsResource_Validate_Locations_And_Azure_Regions` | Multiple locations / azure regions |
| 3086 | `TestAccEnvironmentsResource_Validate_Enable_Admin_Mode` | Exercises `adminMode` / `backgroundOperationsState` — **directly relevant to the unresolved enum questions** |
| 3281 | `TestAccEnvironmentsResource_Create_Environment_With_Env_Group` | Environment + environment group |
| 3529 | `TestAccEnvironmentsResource_Create_Environment_And_Add_Env_Group` | Create then join a group |
| 3620 | `TestAccEnvironmentsResource_Create_Environment_No_Dataverse_Add_Dataverse_And_Add_Env_Group` | Three-way |
| 3705 | `TestAccEnvironmentsResource_Create_Environment_No_Dataverse_Add_Env_Group` | — |
| 3835 | `TestAccEnvironmentsResource_Validate_Update_Environment_Type` | Sandbox ↔ Production SKU change |

Also: internal/services/environment/datasource_environments_test.go line 20 — `TestAccEnvironmentsDataSource_Basic` (read-only list; the cheapest live call, provisions nothing).

Naming helper: `mocks.TestName()` (internal/mocks/mocks.go lines 23-28) returns the calling test's function name via `runtime.Caller`, which is what most tests use for `display_name`.

`make acctest TEST=EnvironmentsResource_Validate_Create USE_PROXY=1` will match `TestAccEnvironmentsResource_Validate_Create` **plus** every test whose name starts with that string (`_Environment_And_Dataverse`, `_No_Dataverse`, `_Them_Try_Remove_Dataverse`, `_Early_Release_Cycle`, ...). Use the longest unique prefix to isolate one test.

---

## 7. Most complete existing BAPI environment fixtures (question 7)

Under internal/services/environment/tests/** — 156 JSON files total. Largest / most attribute-complete (bytes):

| Bytes | Path | Notes |
| --- | --- | --- |
| 12152 | internal/services/environment/tests/datasource/Validate_Read/get_environments.json | **List** response, multiple environments, includes `cluster`, `states.runtime`, `updateCadence`. Best BAPI *list* baseline. |
| 10452 | internal/services/environment/tests/resource/Validate_Create_Environment_And_Dataverse/get_lifecycle_new_dataverse.json | Lifecycle operation payload with the freshly-linked Dataverse block |
| 9608 | internal/services/environment/tests/resource/Validate_Create_With_D365_Template/get_environments.json | Has `templates` / `templateMetadata` / `linkedAppMetadata` |
| **7905** | **internal/services/environment/tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json** | **Single-environment BAPI GET, richest single-object fixture. Primary diff target for the new-API capture.** |
| 7771 | internal/services/environment/tests/resource/Validate_Create_With_Billing_Policy/get_environment_00000000-0000-0000-0000-000000000001.json | Includes `properties/billingPolicy` expansion — relevant to the `billing_policy_id` question |
| 7485 | internal/services/environment/tests/resource/Create_Environment_And_Add_Env_Group/get_environments_1.json | Environment-group membership |
| 7434 | internal/services/environment/tests/resource/Validate_Update_Security_Group_Id/get_environments_0.json and get_environments_1.json | Before/after `securityGroupId` |
| 7357 | internal/services/environment/tests/resource/Validate_Update_With_Billing_Policy/get_environments_1.json .. _3.json | Billing-policy transitions |
| 7357 | internal/services/environment/tests/resource/Validate_Create_And_Update/get_environments_1.json | — |
| 7354 | internal/services/environment/tests/resource/Validate_Update_Environment_Type/get_environments_1.json | SKU change |
| 7120 | internal/services/environment/tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000002.json | Second single-object read fixture |
| 6701 | internal/services/environment/tests/resource/Validate_Create_Environment_And_Dataverse/get_environment_6.json .. _8.json | Post-Dataverse-link state |
| 6218 | internal/services/environment/tests/resource/Validate_Create_Dev_Env/get_environment_00000000-0000-0000-0000-000000000001.json | Developer SKU |
| 6173 | internal/services/environment/tests/resource/Create_Environment_With_Env_Group/get_environment_00000000-0000-0000-0000-000000000001.json | Group-joined single object |
| 6141 | internal/services/environment/tests/resource/Validate_Create_And_Force_Recreate/get_environment_00000000-0000-0000-0000-000000000002.json | — |

Enum values already observable across these fixtures (useful as the "old" side of the diff): `environmentSku` ∈ {`Sandbox`, `Production`, `Developer`}; `azureRegion` ∈ {`westeurope`, `northeurope`, `switzerlandnorth`}; `location` ∈ {`europe`, `switzerland`, `unitedstates`}; `backgroundOperationsState` = `Enabled`; `states.runtime.id` = `Enabled`; `updateCadence.id` ∈ {`Frequent`, `Moderate`}; `baseLanguage` ∈ {1033, 1031}.

Shared cross-service mocks used by every environment unit test: `mocks.ActivateEnvironmentHttpMocks()` (internal/mocks/mocks.go lines 45-112) — registers BAPI `locations`, `environmentLanguages`, `environmentCurrencies`, `validateEnvironmentDetails`, `tenant`, and Dataverse `WhoAmI` / `organizations` / `transactioncurrencies`. **Any new-API fixtures must be accompanied by matching responders, or `RegisterNoResponder` (line 46) will fail the test with "no responder found".**

The exact BAPI URL these fixtures correspond to, from internal/services/environment/api_environment.go lines 249-266:

```text
GET https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{environmentId}
    ?$expand=permissions,properties.capacity,properties/billingPolicy,properties/copilotPolicies
    &api-version=2023-06-01
```

(`constants.BAP_API_VERSION = "2023-06-01"`, constants.go line 199.)

---

## 8. Recommended probe procedure

Goal: capture a raw `GET https://api.powerplatform.com/environmentmanagement/environments/{id}?api-version=2024-10-01` response **and** the equivalent BAPI response for the **same** environment.

1. **Prepare the tenant (one time).** In the target tenant, ensure the Power Platform API service principal exists (`8578e004-a5c6-46e7-913e-12f58912df43`, docs/guides/app_registration.md line 44) and that Azure CLI (`04b07795-8ddb-461a-bbee-02f9e1bf7b46`) is pre-authorized on the app registration (docs/guides/app_registration.md line 49) with at least `EnvironmentManagement.Environments.Read`.
2. **Authenticate.** `az login --allow-no-subscriptions` (or `az login --service-principal --username <CLIENT_ID> --password <CLIENT_SECRET> --tenant <TENANT_ID> --allow-no-subscriptions`). The devcontainer already sets `POWER_PLATFORM_USE_CLI=true` (.devcontainer/devcontainer.json line 44), so the provider will pick this session up automatically. Verify with `az account show`.
3. **Start the capture.** In terminal A: `mitmdump -p 8080 -w /tmp/ppapi.dump --set hardump=/tmp/ppapi.har "~d api.powerplatform.com | ~d api.bap.microsoft.com"`.
4. **Get an environment id.** Either reuse an existing one, or provision one live in terminal B: `make acctest TEST=EnvironmentsResource_Validate_Create USE_PROXY=1 2>&1 | tee /tmp/acctest.log` (target at internal/services/environment/resource_environment_test.go line 1176). Grab the id from the log / from the captured BAPI create flow. **Leave the environment alive** if you want to re-probe; the framework destroys it at the end of the test, so capture inside the window or provision manually with `terraform apply` in `examples/resources/powerplatform_environment/`.
5. **Capture the legacy BAPI response** (through the proxy so it lands in the same flow file):

   ```bash
   BAPI_TOKEN=$(az account get-access-token --resource https://service.powerapps.com --query accessToken -o tsv)
   curl -sS -x http://127.0.0.1:8080 -H "Authorization: Bearer $BAPI_TOKEN" \
     'https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/<ENV_ID>?$expand=permissions,properties.capacity,properties/billingPolicy,properties/copilotPolicies&api-version=2023-06-01' \
     | jq . > /tmp/bapi_get_environment.json
   ```

6. **Capture the new Power Platform API response** for the same id:

   ```bash
   PPAPI_TOKEN=$(az account get-access-token --resource https://api.powerplatform.com --query accessToken -o tsv)
   curl -sS -x http://127.0.0.1:8080 -H "Authorization: Bearer $PPAPI_TOKEN" \
     "https://api.powerplatform.com/environmentmanagement/environments/<ENV_ID>?api-version=2024-10-01" \
     | jq . > /tmp/ppapi_get_environment.json
   # plus the list sibling and the billing-policy sidecar
   curl -sS -x http://127.0.0.1:8080 -H "Authorization: Bearer $PPAPI_TOKEN" \
     "https://api.powerplatform.com/environmentmanagement/environments?api-version=2024-10-01" | jq . > /tmp/ppapi_list_environments.json
   curl -sS -x http://127.0.0.1:8080 -H "Authorization: Bearer $PPAPI_TOKEN" \
     "https://api.powerplatform.com/licensing/environments/<ENV_ID>/billingPolicy?api-version=2024-10-01" | jq . > /tmp/ppapi_billing_policy.json
   ```

7. **Verify the token audience** if a call 401s: `echo "$PPAPI_TOKEN" | cut -d. -f2 | base64 -d 2>/dev/null | jq '{aud, tid, roles, scp}'`. A missing resource principal or a non-pre-authorized CLI shows up here.
8. **Stop the capture and export.** Ctrl-C terminal A (HAR is written on exit), then `jq -r '.log.entries[] | {url: .request.url, status: .response.status}' /tmp/ppapi.har` to index the flows, and extract any body with `jq -r '.log.entries[] | select(.request.url | test("...")) | .response.content.text'`.
9. **Diff old vs new.** `diff <(jq -S 'paths(scalars) | join(".")' /tmp/bapi_get_environment.json) <(jq -S 'paths(scalars) | join(".")' /tmp/ppapi_get_environment.json)` and compare against internal/services/environment/tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json. Focus on the open items: `enterprisePolicies.*.resourceId` vs BAPI `systemId`, `createdFor` vs `usedBy`, and the bare-string enums `state` / `adminMode` / `backgroundOperationsState` / `protectionLevel` / `clusterCategory`.
10. **Anonymize and land the fixtures.** Replace every real GUID with `00000000-0000-0000-0000-00000000000N`, strip tenant ids, org urls, user names, emails, and any token/`Authorization` remnant. Save as `internal/services/environment/tests/<resource|datasource>/<Scenario>/get_environment.json` (and `get_environments.json` for the list). Register matching `httpmock` responders — remember `mocks.ActivateEnvironmentHttpMocks()` installs a `RegisterNoResponder` that hard-fails on any unmocked call (internal/mocks/mocks.go line 46).
11. **Clean up.** Delete the provisioned environment (or let `terraform destroy` / the test teardown run), delete `/tmp/ppapi.dump`, `/tmp/ppapi.har`, and any `*_TOKEN` shell history. **Never commit a capture containing a bearer token.**

---

## 9. Blockers and clarifying questions

* **BLOCKER 1 — no Azure CLI session.** `az` is installed at `/usr/local/bin/az`, but `az account show` returns `ERROR: Please run 'az login' to setup account.` A human must complete an interactive `az login --allow-no-subscriptions` (device code) or supply service-principal credentials. Neither can be done by an agent.
* **BLOCKER 2 — no test tenant documented.** No tenant id, client id, or credential is recorded anywhere in this repository. DEVELOPER.md defers to the external power-platform-terraform-quickstarts bootstrap. The following env vars are **not** set in this container: `POWER_PLATFORM_TENANT_ID`, `POWER_PLATFORM_CLIENT_ID`, `POWER_PLATFORM_CLIENT_SECRET`. Only `POWER_PLATFORM_USE_CLI=true`, `ARM_USE_CLI=true`, `TF_LOG=ERROR`, `TF_CLI_CONFIG_FILE` are present.
* **BLOCKER 3 — `pac` CLI not installed** and not referenced in the repo. Not a fallback.
* **CAVEAT — no automated capture→fixture tooling.** devdocs/adr/mitmproxy.md explicitly lists this as future work. Step 10 above is manual.
* **CAVEAT — `TF_LOG=DEBUG` will not reveal bodies.** Confirmed: no body logging anywhere in `internal/api`. mitmproxy or an explicit probe is mandatory.
* **CAVEAT — `TF_LOG=ERROR` is preset by the devcontainer** and will suppress provider debug output unless overridden per-command.
* **CAVEAT — stale doc.** devdocs/testing_guidelines.md line 58 references a `testAccPreCheck(t)` that does not exist; acceptance tests fail rather than skip when credentials are absent.
* **CAVEAT — `TEST=` is a prefix match.** `make acctest TEST=EnvironmentsResource_Validate_Create` runs 5+ tests, each provisioning real environments.
* **Clarifying question:** which tenant / credential should be used, and is there an existing long-lived environment that can be probed read-only instead of provisioning a new one (which costs ~20-40 min per Dataverse-backed create and consumes capacity)?
* **Clarifying question:** should the `environmentmanagement` probe use delegated (user via `az login`) or application (service-principal) auth? Prior research (.copilot-tracking/research/subagents/2026-08-18/powerplatform-api-environmentmanagement.md section on permissions) flags that only `EnvironmentManagement.Environments.Read` is published — no `.Create`/`.Delete` — so the two contexts may return different payloads and must both be captured.

---

## 10. Recommended next research (not completed here)

* [ ] Execute the probe once a tenant + `az login` is available; record the actual `GET /environmentmanagement/environments/{id}` body.
* [ ] Confirm whether `az account get-access-token --resource https://api.powerplatform.com` succeeds without a bespoke app registration (i.e. whether Azure CLI is pre-authorized on the first-party Power Platform API by default) — untestable here without a tenant.
* [ ] Determine the 202 response header shape (`Location` vs `Operation-Location`, `Retry-After`) for `POST /environmentmanagement/provisioning/environments` — requires a live create.
* [ ] Compare delegated vs application-token responses for the same environment id.
* [ ] Evaluate adding a small committed probe helper (e.g. `make probe URL=...`) plus a mitm-flow→fixture converter script, closing the gap devdocs/adr/mitmproxy.md leaves open.
* [ ] Check whether `powerplatform_environment` acceptance tests need `CheckDestroy` added — none currently have it, so a failed probe run can leak environments.
