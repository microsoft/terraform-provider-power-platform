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
	MAKER_WELCOME_CONTENT      = "MakerOnboarding"
	MAKER_ONBOARDING_URL       = "makerOnboardingUrl"
	MAKER_ONBOARDING_MARKDOWN  = "makerOnboardingMarkdown"
	MAKER_ONBOARDING_TIMESTAMP = "makerOnboardingTimestamp"
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
