# Running with elevated privileges is required to install PowerShell 7
#Start-Process powershell.exe -Verb RunAs -ArgumentList ('-noprofile -noexit -file "{0}" -elevated' -f ($myinvocation.MyCommand.Definition))

# Download the PowerShell 7 installer
Invoke-WebRequest -Uri "https://aka.ms/install-powershell.ps1" -OutFile "install-powershell.ps1" -UseBasicParsing

# Run the installer script silently
.\install-powershell.ps1 -Quiet

# Clean up the installer file
Remove-Item -Path "install-powershell.ps1"
