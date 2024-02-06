package settings

const (
	ADMIN_ANALYTICS_SCOPE    = "https://adminanalytics.powerplatform.microsoft.com/.default"
	SERVICE_POWERAPPS_SCOPE  = "https://service.powerapps.com/.default"
	POWER_PLATFORM_API_SCOPE = "https://api.powerplatform.com/.default"

	CLIENT_ID = "1950a258-227b-4e31-a9cf-717495945fc2"

	MSAL_CACHE_FILE_NAME = "terraform_power_platform_cache.dat"

	OAUTH_AUTHORITY_URL = "https://login.microsoftonline.com/"

	CACHE_SAVE_FOLDER_PATH_LINUX   = ".local/share/Microsoft/TerraformPowerPlatformProvider"
	CACHE_SAVE_FOLDER_PATH_WINDOWS = "Microsoft\\Terraform Power Platform Provider"
)

var REQUIRED_SCOPES = []string{
	ADMIN_ANALYTICS_SCOPE,
	SERVICE_POWERAPPS_SCOPE,
	POWER_PLATFORM_API_SCOPE,
}
