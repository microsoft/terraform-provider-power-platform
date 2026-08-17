// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

const (
	NOT_SPECIFIED = "NotSpecified"
	APP           = "App"
)

const (
	SHARING                            = "Sharing"
	CAN_SHARE_WITH_SECURITY_GROUPS     = "CanShareWithSecurityGroups"
	IS_GROUP_SHARING_DISABLED          = "IsGroupSharingDisabled"
	MAXIMUM_SHARE_LIMIT                = "MaximumShareLimit"
	NO_LIMIT                           = "noLimit"
	EXCLUDE_SHARING_TO_SECURITY_GROUPS = "excludeSharingToSecurityGroups"
)

const (
	USAGE_INSIGHTS                    = "AdminDigest"
	INCLUDE_ON_HOME_PAGE_INSIGHTS     = "IncludeOnHomePageInsights"
	EXCLUDE_ENVIRONMENT_FROM_ANALYSIS = "ExcludeEnvironmentFromAnalysis"
)

const (
	MAKER_WELCOME_CONTENT             = "MakerOnboarding"
	MAKER_ONBOARDING_URL              = "makerOnboardingUrl"
	MAKER_ONBOARDING_MARKDOWN         = "makerOnboardingMarkdown"
	MAKER_ONBOARDING_TIMESTAMP        = "makerOnboardingTimestamp"
	MAKER_ONBOARDING_CONTENT_ID       = "MakerOnboardingContent"
	MAKER_ONBOARDING_CONTENT_VERSION  = "1.0"
	MAKER_ONBOARDING_PORTALS          = "makerOnboardingPortals"
	MAKER_ONBOARDING_CONSENT_REQUIRED = "makerOnboardingConsentRequired"
)

const (
	SOLUTION_CHECKER_ENFORCEMENT    = "SolutionChecker"
	SOLUTION_CHECKER_MODE           = "solutionCheckerMode"
	SUPPRESS_VALIDATION_EMAILS      = "suppressValidationEmails"
	SOLUTION_CHECKER_RULE_OVERRIDES = "solutionCheckerRuleOverrides"
)

const (
	BACKUP_RETENTION = "Lifecycle"
	RETENTION_PERIOD = "RetentionPeriod"
)

const (
	AI_GENERATED_DESC                 = "Copilot"
	DISABLE_AI_GENERATED_DESCRIPTIONS = "DisableAiGeneratedDescriptions"
)

const (
	AI_GENERATIVE_SETTINGS                  = "GenerativeAISettings"
	CROSS_GEO_COPILOT_DATA_MOVEMENT_ENABLED = "crossGeoCopilotDataMovementEnabled"
	BING_CHAT_ENABLED                       = "bingChatEnabled"
)

// Rule sets below are served by the rule-based policies API rather than the legacy rule sets API.
const (
	POLICY_API_VERSION = "2024-10-01"

	ADVANCED_CONNECTOR_POLICIES_ONLY_ID    = "AdvancedConnectorPoliciesOnly"
	ADVANCED_CONNECTOR_POLICIES_VERSION    = "1.0"
	ENABLE_ADVANCED_CONNECTOR_POLICIES_KEY = "EnableAdvancedConnectorPoliciesOnly"

	CONTENT_SECURITY_POLICY_ID      = "PowerAppsContentSecurityPolicy"
	CONTENT_SECURITY_POLICY_VERSION = "1.0"

	CSP_IS_ENABLED_KEY                  = "IsContentSecurityPolicyEnabled"
	CSP_IS_ENABLED_FOR_CANVAS_KEY       = "IsContentSecurityPolicyEnabledForCanvas"
	CSP_IS_ENABLED_FOR_CODE_APPS_KEY    = "IsContentSecurityPolicyEnabledForCodeApps"
	CSP_OPTIONS_KEY                     = "ContentSecurityPolicyOptions"
	CSP_REPORT_URI_KEY                  = "ContentSecurityPolicyReportUri"
	CSP_REPORTING_ENDPOINT_KEY          = "ContentSecurityPolicyReportingEndpoint"
	CSP_CONFIGURATION_KEY               = "ContentSecurityPolicyConfiguration"
	CSP_CONFIGURATION_FOR_CANVAS_KEY    = "ContentSecurityPolicyConfigurationForCanvas"
	CSP_CONFIGURATION_FOR_CODE_APPS_KEY = "ContentSecurityPolicyConfigurationForCodeApps"

	CONNECTOR_MANAGEMENT_ID               = "ConnectorManagement"
	CONNECTOR_MANAGEMENT_VERSION          = "1.0"
	CONNECTOR_MANAGEMENT_ALLOWED_LIST_KEY = "AllowedConnectorList"
	CONNECTOR_API_PREFIX                  = "/providers/Microsoft.PowerApps/apis/"
)
