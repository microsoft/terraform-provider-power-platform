$Psversion = (Get-Host).Version

if($Psversion.Major -ge 7)
{

if (!(Get-Module "DataGateway")) {
Install-Module -Name DataGateway 
}

#Gateway Login

Connect-DataGatewayServiceAccount -ApplicationId $ApplicationId -ClientSecret $securePassword  -Tenant $Tenant


#Installing Gateway

Install-DataGateway -AcceptConditions 


#Configuring Gateway

$clusters = Get-DataGatewayCluster | Select -Property Id,Name
#loop clusters to check if gateway is already configured
foreach($cluster in $clusters){
    if($cluster.Name -eq $GatewayName){
        Write-Host "Gateway is already configured"
        $clusterId = $cluster.Id
    }
}

if ($clusterId -eq $null){
    $clusterId = Add-DataGatewayCluster -Name $GatewayName -RecoveryKey  $RecoverKey -OverwriteExistingGateway
}


#Add User as Admin
Write-Host "Adding user as Admin"
Add-DataGatewayClusterUser -GatewayClusterId $clusterId -PrincipalObjectId $userIDToAddasAdmin -AllowedDataSourceTypes $null -Role Admin

}
else{
    Write-Host "Please install Powershell 7 or above"
exit 1
}
