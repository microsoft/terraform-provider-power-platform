clear
 
# URL Parameter
$WebURL = "https://sdlc-esd.oracle.com/ESD6/JSCDL/jdk/8u381-b09/8c876547113c4e4aab3c868e9e0ec572/jre-8u381-windows-x64.exe?GroupName=JSC&FilePath=/ESD6/JSCDL/jdk/8u381-b09/8c876547113c4e4aab3c868e9e0ec572/jre-8u381-windows-x64.exe&BHost=javadl.sun.com&File=jre-8u381-windows-x64.exe&AuthParam=1695229341_ac11b257b629a57b5efebd434db8987a&ext=.exe"
  
# Directory Parameter
$FileDirectory = "$($env:USERPROFILE)$("\downloads\")"
 
#Write-Output $FileDirectory
 
# If directory doesn't exist create the directory
if((Test-Path $FileDirectory) -eq 0)
    {
        mkdir $FileDirectory;
    }
 
# We assume the file you download is named what you want it to be on your computer
$FileName = [System.IO.Path]::GetFileName($WebURL)
 
# Concatenate the two values to prepare the download
$FullFilePath = "$($FileDirectory)$($FileName)"
 
#Write-Output $FullFilePath
 
function Get-FileDownload([String] $WebURL, [String] $FullFilePath)
{
        # Give a basic message to the user to let them know what we are doing
        Write-Output "Downloading '$WebURL' to '$FullFilePath'"
 
        $uri = New-Object "System.Uri" "$WebURL"
        $request = [System.Net.HttpWebRequest]::Create($uri) 
        $request.set_Timeout(30000) #15 second timeout 
        $response = $request.GetResponse() 
        $totalLength = [System.Math]::Floor($response.get_ContentLength()/1024) 
        $responseStream = $response.GetResponseStream() 
        $targetStream = New-Object -TypeName System.IO.FileStream -ArgumentList $FullFilePath, Create 
        $buffer = new-object byte[] 10KB 
        $count = $responseStream.Read($buffer,0,$buffer.length) 
        $downloadedBytes = $count
        while ($count -gt 0) 
            { 
                [System.Console]::Write("`r`nDownloaded {0}K of {1}K", [System.Math]::Floor($downloadedBytes/1024), $totalLength) 
                $targetStream.Write($buffer, 0, $count) 
                $count = $responseStream.Read($buffer,0,$buffer.length) 
                $downloadedBytes = $downloadedBytes + $count
            } 
         
        $targetStream.Flush()
        $targetStream.Close() 
        $targetStream.Dispose() 
        $responseStream.Dispose() 
         
        # Give a basic message to the user to let them know we are done
        Write-Output "`r`nDownload complete"
    }
 
function AddSystemPaths([array] $PathsToAdd) {
#http://blogs.technet.com/b/sqlthoughts/archive/2008/12/12/powershell-function-to-add-system-path.aspx
 
  
  $VerifiedPathsToAdd = ""
  
  foreach ($Path in $PathsToAdd) {
    if ($Env:Path -like "*$Path*") {
      echo "  Path to $Path already added"
    }
    else {
      $VerifiedPathsToAdd += ";$Path";echo "  Path to $Path needs to be added"
    }
  }
  
  if ($VerifiedPathsToAdd -ne "") {
    echo "Adding paths: $VerifiedPathsToAdd"
    [System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + "$VerifiedPathsToAdd","Machine")
    echo "Note: The new path does NOT take immediately in running processes. Only new processes will see new path."
  }
}
 
 
Get-FileDownload $WebURL $FullFilePath
 
cd $FileDirectory
 
"INSTALL_SILENT=Enable" | Set-Content "$FileDirectory/JavaInstallConfig.txt"
"INSTALLDIR=C:\java" | Add-Content "$FileDirectory/JavaInstallConfig.txt"
"AUTO_UPDATE=Enable" | Add-Content "$FileDirectory/JavaInstallConfig.txt"
"WEB_JAVA_SECURITY_LEVEL=VH" | Add-Content "$FileDirectory/JavaInstallConfig.txt"
 
start-process $FullFilePath INSTALLCFG=$FileDirectory/JavaInstallConfig.txt -Wait
 
AddSystemPaths ("C:\java\bin")
