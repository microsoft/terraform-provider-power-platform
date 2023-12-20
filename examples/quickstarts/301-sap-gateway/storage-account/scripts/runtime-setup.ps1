param (
    [Parameter(Mandatory=$true)]
    [string]$keyVaultUri,
    [Parameter(Mandatory=$true)]
    [string]$secretNamePP,
    [Parameter(Mandatory=$true)]
    [string]$userAdmin,
    [Parameter(Mandatory=$true)]
    [string]$secretNameIRKey,
    [Parameter(Mandatory=$true)]
    [string]$ApplicationId,
    [Parameter(Mandatory=$true)]
    [string]$TenantId,
    [Parameter(Mandatory=$true)]
    [string]$GatewayName,
    [Parameter(Mandatory=$true)]
    [string]$SecretNameRecoverKey
)

$Psversion = (Get-Host).Version

if($Psversion.Major -ge 7)
{
    
    #Retrieve the secret from Key Vault
    Write-Output "Retrieve the secrete from Key Vault"
    $Response = Invoke-RestMethod -Uri 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fvault.azure.net' -Method GET -Headers @{Metadata="true"}
    $KeyVaultToken = $Response.access_token

    $Uri = "$keyVaultUri" + "secrets/" + $secretNamePP + "?api-version=2016-10-01"
    $SecretPP = Invoke-RestMethod -Uri $Uri -Method GET -Headers @{Authorization="Bearer $KeyVaultToken"}
    $securePassword = $SecretPP.value | ConvertTo-SecureString -AsPlainText -Force;

    $Uri = "$keyVaultUri" + "secrets/" + $SecretNameRecoverKey + "?api-version=2016-10-01"
    $RecoverKey = Invoke-RestMethod -Uri $Uri -Method GET -Headers @{Authorization="Bearer $KeyVaultToken"}
    $RecoverKey = $RecoverKey.value | ConvertTo-SecureString -AsPlainText -Force;
    $userIDToAddasAdmin = $userAdmin

    #Gateway Login
    Write-Output "Gateway Login"
    Connect-DataGatewayServiceAccount -ApplicationId $ApplicationId -ClientSecret $securePassword -Tenant $TenantId

    #Installing Gateway
    Write-Output "Installing Gateway"
    Install-DataGateway -AcceptConditions 

    #Configuring Gateway
    $GatewayObjectId = (Get-DataGatewayCluster | Where-Object {$_.Name -eq "OPDGW-SAPAzureIntegration"}).Id

    if (![string]::IsNullOrEmpty($GatewayObjectId)) {
        Write-Output "Remove Cluster"
        Remove-DataGatewayCluster -GatewayClusterId $GatewayObjectId
    }
    
    Write-Output "Add Cluster"
    $GatewayDetails = Add-DataGatewayCluster -Name $GatewayName -RecoveryKey  $RecoverKey -RegionKey westus3 -OverwriteExistingGateway
    $GatewayObjectId = $GatewayDetails.GatewayObjectId

    Write-Output "$GatewayName ID: $GatewayObjectId"
    #Add User as Admin
    Write-Output "Add User as Admin"
    Add-DataGatewayClusterUser -GatewayClusterId $GatewayObjectId -PrincipalObjectId $userIDToAddasAdmin -AllowedDataSourceTypes $null -Role Admin -RegionKey westus3

#####################################################################################
    Write-Output "Installing SHIR - Self-hosted Integration Runtime."

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
    #$Response = Invoke-RestMethod -Uri 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fvault.azure.net' -Method GET -Headers @{Metadata="true"}
    #$KeyVaultToken = $Response.access_token
    $Uri = "$keyVaultUri" + "secrets/" + $secretNameIRKey + "?api-version=2016-10-01"
    Write-Output $Uri
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

}
else{
    Write-Output "PowerShell version 7 or higher is required to run this script."
    exit 1
}
