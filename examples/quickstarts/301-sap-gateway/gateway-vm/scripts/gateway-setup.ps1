param (
    [Parameter(Mandatory=$true)]
    [string]$secretPP,
    [Parameter(Mandatory=$true)]
    [string]$userAdmin
)

$Psversion = (Get-Host).Version

if($Psversion.Major -ge 7)
{

if (!(Get-Module "DataGateway")) {
Install-Module -Name DataGateway -Force
}

$securePassword = $secretPP | ConvertTo-SecureString -AsPlainText -Force;
$ApplicationId ="2d0b62aa-765d-4e0f-b7f2-61debc6611d7";
$Tenant = "0d7fbacd-d6d8-4652-9f58-ae0f94edde5c";
$GatewayName = "OPDGW-SAPAzureIntegration";
$RecoverKey = "recover01" | ConvertTo-SecureString -AsPlainText -Force;
$userIDToAddasAdmin = $userAdmin


#Gateway Login

Connect-DataGatewayServiceAccount -ApplicationId $ApplicationId -ClientSecret $securePassword -Tenant $Tenant


#Installing Gateway

Install-DataGateway -AcceptConditions 


#Configuring Gateway
$GatewayObjectId = (Get-DataGatewayCluster | Where-Object {$_.Name -eq "OPDGW-SAPAzureIntegration"}).Id

if([string]::IsNullOrEmpty($GatewayObjectId)) {
Write-Host "Add Cluster"
$GatewayDetails = Add-DataGatewayCluster -Name $GatewayName -RecoveryKey  $RecoverKey
$GatewayObjectId = $GatewayDetails.GatewayObjectId
}

Write-Host $GatewayObjectId
#Add User as Admin
Add-DataGatewayClusterUser -GatewayClusterId $GatewayObjectId -PrincipalObjectId $userIDToAddasAdmin -AllowedDataSourceTypes $null -Role Admin -RegionKey westus3

}
else{
exit 1
}
