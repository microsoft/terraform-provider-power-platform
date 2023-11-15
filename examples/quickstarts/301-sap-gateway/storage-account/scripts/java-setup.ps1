Write-Output "Checking if Java is installed..."
$java = Get-WmiObject -Class win32_product | where {$_.Name -like "*Java*"}
if ($java) {
    Write-Output "Java is already installed."
} else {
    $url = "https://javadl.oracle.com/webapps/download/AutoDL?BundleId=248774_8c876547113c4e4aab3c868e9e0ec572"
    $file = "$env:TEMP\jre.exe"
    Invoke-WebRequest -Uri $url -OutFile $file
    Start-Process -FilePath $file -ArgumentList "/s" -Wait
    Write-Output "Java has been installed."
}
