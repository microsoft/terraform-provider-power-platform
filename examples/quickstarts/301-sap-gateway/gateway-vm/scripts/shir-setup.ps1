param (
    [Parameter(Mandatory=$true)]
    [string]$keyVaultUri,
    [Parameter(Mandatory=$true)]
    [string]$irKeyVaultName
)

# Define the download URL for the MSI installer.
# Please, update the version number to the latest version.
$downloadUrl = "https://download.microsoft.com/download/E/4/7/E4771905-1079-445B-8BF9-8A1A075D8A10/IntegrationRuntime_5.32.8600.2.msi"

# Define the path where the installer will be downloaded.
$installerPath = "C:\sapint\IntegrationRuntime.msi"

# Download the Integration Runtime MSI installer.
Write-Output "Downloading Integration Runtime installer..."
Invoke-WebRequest -Uri $downloadUrl -OutFile $installerPath

# Define the installation arguments for a silent install
$installArguments = '/qn /norestart'

# Install Integration Runtime silently.
Write-Output "Installing Integration Runtime..."
Start-Process -Wait -FilePath $installerPath -ArgumentList $installArguments

# Define the path to the Integration Runtime command line tool.
$irPath = "C:\Program Files\Microsoft Integration Runtime\5.0\Shared\dmgcmd.exe"

#Retrieve the Integration Runtime key from Key Vault
Write-Output "Retrieve the Integration Runtime key from Key Vault"
$Response = Invoke-RestMethod -Uri 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fvault.azure.net' -Method GET -Headers @{Metadata="true"}
$KeyVaultToken = $Response.access_token
$Uri = "$keyVaultUri" + "secrets/" + $irKeyVaultName + "?api-version=2016-10-01"
$irKey = Invoke-RestMethod -Uri $Uri -Method GET -Headers @{Authorization="Bearer $KeyVaultToken"}
    
# Define the installation arguments for registering the Integration Runtime.
$irArguments = '-k ' + $irKey.value

# Register the Integration Runtime.
Write-Output "Registering the Integration Runtime..."
Start-Process -Wait -FilePath $irPath -ArgumentList $irArguments

# Define the restart arguments
$irArguments = '-r'

# Restart the Intergration Runtime
Write-Output "Restarting the Integration Runtime..."
Start-Process -Wait -FilePath $irPath -ArgumentList $irArguments
