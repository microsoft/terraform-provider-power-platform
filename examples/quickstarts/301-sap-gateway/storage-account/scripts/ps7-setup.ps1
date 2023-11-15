# Define the download URL for the MSI installer.
$downloadUrl = "https://github.com/PowerShell/PowerShell/releases/download/v7.3.9/PowerShell-7.3.9-win-x64.msi"

# Define the path where the installer will be downloaded.
$installerPath = "$env:TEMP\PowerShell.msi"

# Download the PowerShell 7 MSI installer.
Invoke-WebRequest -Uri $downloadUrl -OutFile $installerPath

# Define the installation arguments for a silent install
$installArguments = '/qn /norestart INSTALLFOLDER="C:\powershell7"'

# Install PowerShell 7 silently.
Start-Process -Wait -FilePath $installerPath -ArgumentList $installArguments

# Clean up the downloaded installer.
Remove-Item -Path $installerPath
