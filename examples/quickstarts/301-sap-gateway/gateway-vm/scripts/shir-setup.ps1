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
$installerPath = "C:\shir\IntegrationRuntime.msi"

# Download the Intergration Runtime MSI installer.
Invoke-WebRequest -Uri $downloadUrl -OutFile $installerPath

# Define the installation arguments for a silent install
$installArguments = '/qn /norestart'

# Install Intergration Runtime silently.
Start-Process -Wait -FilePath $installerPath -ArgumentList $installArguments

# Define the path to the Integration Runtime command line tool.
$irPath = "C:\Program Files\Microsoft Integration Runtime\5.0\Shared\dmgcmd.exe"

#Retrieve the Integration Runtime key from Key Vault
$Response = Invoke-RestMethod -Uri 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fvault.azure.net' -Method GET -Headers @{Metadata="true"}
$KeyVaultToken = $Response.access_token
$Uri = "$keyVaultUri" + "secrets/" + $irKeyVaultName + "?api-version=2016-10-01"
$irKey = Invoke-RestMethod -Uri $Uri -Method GET -Headers @{Authorization="Bearer $KeyVaultToken"}
    
# Define the installation arguments for registering the Integration Runtime.
$irArguments = '-k ' + $irKey.value

# Register the Integration Runtime.
Start-Process -Wait -FilePath $irPath -ArgumentList $irArguments
#Start-Process -Wait -FilePath "irKeyVaultNamedmgcmd.exe" -ArgumentList "-k IR@c372143d-35ce-4985-9028-9d726c8f30ff@s4hana-poc@we@uQOLE/v1RUJ/CKx6xrkVLb4f7WfJ0lywxds3nrB25Ww="

# Define the restart arguments
$irArguments = '-r'

# Restart the Intergration Runtime
Start-Process -Wait -FilePath $irPath -ArgumentList $irArguments
