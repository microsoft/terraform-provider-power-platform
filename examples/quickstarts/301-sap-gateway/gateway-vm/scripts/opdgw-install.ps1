$Psversion = (Get-Host).Version

if($Psversion.Major -ge 7)
{

    Write-Output "Installing DataGateway Module"
    Install-Module -Name DataGateway -Force

}
else{
    Write-Output "PowerShell version 7 or higher is required to run this script."
    exit 1
}
