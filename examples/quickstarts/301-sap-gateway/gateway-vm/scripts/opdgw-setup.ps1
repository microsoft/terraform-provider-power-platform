param (
    [Parameter(Mandatory=$true)]
    [string]$keyVaultUri,
    [Parameter(Mandatory=$true)]
    [string]$secretPPName,
    [Parameter(Mandatory=$true)]
    [string]$userAdmin
)

$Psversion = (Get-Host).Version

if($Psversion.Major -ge 7)
{

    #Retrieve the secret from Key Vault
    Write-Host "Retrieve the secrete from Key Vault"
    $Response = Invoke-RestMethod -Uri 'http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fvault.azure.net' -Method GET -Headers @{Metadata="true"}
    $KeyVaultToken = $Response.access_token

    $Uri = "$keyVaultUri" + "secrets/" + $secretPPName + "?api-version=2016-10-01"
    $SecretPP = Invoke-RestMethod -Uri $Uri -Method GET -Headers @{Authorization="Bearer $KeyVaultToken"}

    $securePassword = $SecretPP.value | ConvertTo-SecureString -AsPlainText -Force;
    $ApplicationId ="2d0b62aa-765d-4e0f-b7f2-61debc6611d7";
    $Tenant = "0d7fbacd-d6d8-4652-9f58-ae0f94edde5c";
    $GatewayName = "OPDGW-SAPAzureIntegration";
    $RecoverKey = "recover01" | ConvertTo-SecureString -AsPlainText -Force;
    $userIDToAddasAdmin = $userAdmin


    #Gateway Login
    Write-Host "Gateway Login"
    Connect-DataGatewayServiceAccount -ApplicationId $ApplicationId -ClientSecret $securePassword -Tenant $Tenant


    #Installing Gateway
    Write-Host "Installing Gateway"
    Install-DataGateway -AcceptConditions 


    #Configuring Gateway
    $GatewayObjectId = (Get-DataGatewayCluster | Where-Object {$_.Name -eq "OPDGW-SAPAzureIntegration"}).Id

    if (![string]::IsNullOrEmpty($GatewayObjectId)) {
        Write-Host "Remove Cluster"
        Remove-DataGatewayCluster -GatewayClusterId $GatewayObjectId
    }
    
    Write-Host "Add Cluster"
    $GatewayDetails = Add-DataGatewayCluster -Name $GatewayName -RecoveryKey  $RecoverKey -RegionKey westus3 -OverwriteExistingGateway
    $GatewayObjectId = $GatewayDetails.GatewayObjectId

    Write-Host "$GatewayName ID: $GatewayObjectId"
    #Add User as Admin
    Write-Host "Add User as Admin"
    Add-DataGatewayClusterUser -GatewayClusterId $GatewayObjectId -PrincipalObjectId $userIDToAddasAdmin -AllowedDataSourceTypes $null -Role Admin -RegionKey westus3

}
else{
    Write-Host "PowerShell version 7 or higher is required to run this script."
    exit 1
}
