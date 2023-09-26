# Download the PowerShell 7 installer
Invoke-WebRequest -Uri "https://aka.ms/install-powershell.ps1" -OutFile "install-powershell.ps1"

# Run the installer script silently
.\install-powershell.ps1 -Quiet

# Clean up the installer file
Remove-Item -Path "install-powershell.ps1"
