# Resource Environment Issues - Merged Issues

## ISSUE 1

# Function Naming Inconsistency: dvExits

##

/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install.go

## Problem

The variable name `dvExits` is likely a typo and may have been intended to be `dvExists`. Naming is crucial for maintainability and readability. Typos can cause confusion for future readers and contributors.

## Impact

Severity: **Low**

This is a minor naming issue; however, typos in variable names reduce code readability and can lead to mistakes if the code is copied or modified.

## Location

Lines in the `Create` method:

```go
dvExits, err := r.ApplicationClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
    return
}

if !dvExits {
    resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
    return
}
```

## Code Issue

```go
dvExits, err := r.ApplicationClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
    return
}

if !dvExits {
    resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
    return
}
```

## Fix

Rename `dvExits` to `dvExists` in all locations within the method.

```go
dvExists, err := r.ApplicationClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
    return
}

if !dvExists {
    resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
    return
}
```
---

This output will be saved as:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_environment_application_package_install.go_dvExits_naming-low.md`


---

## ISSUE 2

# Misspelled Test Function Name: `TestUnitEnvirionmentGroupResource_Validate_Create`

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group_test.go

## Problem

The test function `TestUnitEnvirionmentGroupResource_Validate_Create` contains a typo in its name: "Envirionment" instead of "Environment". This can hinder searchability and consistency in test naming conventions.

## Impact

This issue impacts the readability and maintainability of the codebase. It can make it difficult for contributors to find this test via standard search (e.g., `TestUnitEnvironmentGroupResource_Validate_Create`). Severity: **low**.

## Location

Line: Function declaration of `TestUnitEnvirionmentGroupResource_Validate_Create`

## Code Issue

```go
func TestUnitEnvirionmentGroupResource_Validate_Create(t *testing.T) {
```

## Fix

Rename the function to correct the typo in "Envirionment":

```go
func TestUnitEnvironmentGroupResource_Validate_Create(t *testing.T) {
```


---

## ISSUE 3

# Title

Misspelled Method Name: `aiGenerativeFeaturesValidaor` Should be `aiGenerativeFeaturesValidator`

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go

## Problem

The method name `aiGenerativeFeaturesValidaor` contains a typo and does not follow the expected English spelling convention ("Validator" instead of "Validaor"). This is inconsistent with general Go naming conventions, making the code less readable and potentially confusing for future maintainers.

## Impact

- **Readability**: Reduces code readability and professionalism.
- **Maintainability**: Makes searching and reusing code elements less straightforward.
- **Severity**: Low (no functional problem, but should be corrected for maintainability and standardization).

## Location

Lines with the following code:

```go
func (r *Resource) aiGenerativeFeaturesValidaor(plan *SourceModel) error {
```

and all calls to this method (e.g., lines in `Create` and `Update` methods).

## Code Issue

```go
func (r *Resource) aiGenerativeFeaturesValidaor(plan *SourceModel) error {
    // implementation ...
}
...
err = r.aiGenerativeFeaturesValidaor(plan)
```

## Fix

Update the method name and all references to it to use the correct spelling, `aiGenerativeFeaturesValidator`.

```go
// Function definition
func (r *Resource) aiGenerativeFeaturesValidator(plan *SourceModel) error {
    if r.EnvironmentClient.Api.Config.CloudType != config.CloudTypePublic {
        return errors.New("moving data across regions is not supported in non public clouds")
    }
    if plan.Location.ValueString() == "unitedstates" && plan.AllowMovingDataAcrossRegions.ValueBool() {
        return errors.New("moving data across regions is not supported in the unitedstates location")
    }
    if plan.Location.ValueString() != "unitedstates" && plan.AllowBingSearch.ValueBool() && !plan.AllowMovingDataAcrossRegions.ValueBool() {
        return errors.New("to enable ai generative features, moving data across regions must be enabled")
    }
    return nil
}

// Update all usages as well
err = r.aiGenerativeFeaturesValidator(plan)
```

---

This change makes the codebase more consistent and maintainable.

---


---

## ISSUE 4

# Inconsistent Attribute Naming: Dataverse vs. Environment

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

Throughout the test cases, the term `dataverse` is used for a nested block and its attributes are frequently mixed with other resource attributes (e.g., `dataverse.language_code`, `dataverse.organization_id`). However, other attributes are named at the top level (like `location`, `display_name`).

While this may reflect the schema of the actual resource, within Go code and tests, inconsistent scoping and naming for subresource attributes can reduce clarity. In some config snippets, sometimes the attribute nesting/scope is ambiguous ("dataverse" vs "environment"), which could lead to confusion in implementation, documentation, or onboarding.

## Impact

- **Severity: Low**
- This is a minor maintainability and readability nuisance, and can lead to confusion for contributors/readers unfamiliar with schema conventions.
- Not technically incorrect, but expert guidance suggests naming should be predictable and consistently scoped.

## Location

Example (from many spots in file):

```go
resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
```

## Code Issue

```go
resource "powerplatform_environment" "development" {
    display_name = "..."
    location = "..."
    dataverse = {
        language_code = "1033"
        currency_code = "PLN"
        // ...
    }
}
```

## Fix

Consider making a clarification in both test helper naming and documentation as to why "dataverse" is at this level. Optionally, apply a clear prefix everywhere or encapsulate subresource checks and config generation in helper functions to minimize ambiguity.

```go
// Helper/clarification
func checkDataverseAttr(attr, expected string) resource.TestCheckFunc {
    return resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse."+attr, expected)
}

// Usage
Check: resource.ComposeTestCheckFunc(
    resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
    checkDataverseAttr("currency_code", "PLN"),
)
```


Also consider documenting this distinction at the top of the test file, so that reviewers and new maintainers can tell top-level vs. nested resource properties at a glance.


---

## ISSUE 5

# Title

Naming: Constant Naming Does Not Follow Go Conventions

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resources_environment_settings.go

## Problem

The constant `SERVICE_TAGS_NAMES` is in all uppercase with underscores, which doesn't follow the Go naming convention where constants should use `CamelCase` (e.g., `ServiceTagsNames`). This makes the code inconsistent with other Go code and can hinder readability.

## Impact

Code inconsistency with Go's idiomatic style may make the code harder to read and maintain for Go developers. Severity: low.

## Location

```go
const SERVICE_TAGS_NAMES = "ApiManagement,AppConfiguration,AppService,ActionGroup,AppServiceManagement,ApplicationInsightsAvailability,AutonomousDevelopmentPlatform,AzureActiveDirectory,AzureAdvancedThreatProtection,AzureArcInfrastructure,AzureAttestation,AzureBackup,AzureBotService,AzureCognitiveSearch,AzureConnectors,AzureContainerRegistry,AzureCosmosDB,AzureDataExplorerManagement,AzureDataLake,AzureDatabricks,AzureDevOps,AzureDevSpaces,AzureDeviceUpdate,AzureDigitalTwins,AzureEventGrid,AzureHealthcareAPIs,AzureInformationProtection,AzureIoTHub,AzureKeyVault,AzureLoadTestingInstanceManagement,AzureMachineLearning,AzureMachineLearningInference,AzureManagedGrafana,AzureMonitorForSAP,AzureMonitor,AzureOpenDatasets,AzurePortal,AzureRemoteRendering,AzureResourceManager,AzureSecurityCenter,AzureSentinel,AzureSignalR,AzureSiteRecovery,AzureSphere,AzureSpringCloud,AzureStack,AzureTrafficManager,AzureUpdateDelivery,AzureWebPubSub,BatchNodeManagement,ChaosStudio,CognitiveServicesFrontend,CognitiveServicesManagement,ContainerAppsManagement,DataFactory,Dynamics365BusinessCentral,Dynamics365ForMarketingEmail,Dynamics365FraudProtection,EOPExternalPublishedIPs,EventHub,GatewayManager,Grafana,GuestAndHybridManagement,HDInsight,KustoAnalytics,LogicApps,M365ManagementActivityApi,M365ManagementActivityApiWebhook,Marketplace,MicrosoftAzureFluidRelay,MicrosoftCloudAppSecurity,MicrosoftContainerRegistry,MicrosoftDefenderForEndpoint,MicrosoftPurviewPolicyDistribution,OneDsCollector,PowerBI,PowerPlatformPlex,PowerQueryOnline,SCCservice,Scuba,SecurityCopilot,SerialConsole,ServiceBus,ServiceFabric,Sql,SqlManagement,Storage,StorageMover,StorageSyncService,VideoIndexer,WindowsAdminCenter,WindowsVirtualDesktop"
```

## Code Issue

```go
const SERVICE_TAGS_NAMES = "ApiManagement,AppConfiguration,AppService,ActionGroup,AppServiceManagement,ApplicationInsightsAvailability,..."
```

## Fix

Rename the constant using Go's CamelCase for increased readability:

```go
const ServiceTagsNames = "ApiManagement,AppConfiguration,AppService,ActionGroup,AppServiceManagement,ApplicationInsightsAvailability,..."
```
Apply this change wherever `SERVICE_TAGS_NAMES` is referenced.


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
