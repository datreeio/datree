#Requires -Version 5

$old_erroractionpreference = $erroractionpreference
$erroractionpreference = 'stop' # quit if anything goes wrong

if (($PSVersionTable.PSVersion.Major) -lt 5) {
    Write-Output "PowerShell 5 or later is required to run Datree."
    Write-Output "Upgrade PowerShell: https://docs.microsoft.com/en-us/powershell/scripting/setup/installing-windows-powershell"
    break
}

# show notification to change execution policy:
$allowedExecutionPolicy = @('Unrestricted', 'RemoteSigned', 'ByPass')
if ((Get-ExecutionPolicy).ToString() -notin $allowedExecutionPolicy) {
    Write-Output "PowerShell requires an execution policy in [$($allowedExecutionPolicy -join ", ")] to run Datree."
    Write-Output "For example, to set the execution policy to 'RemoteSigned' please run :"
    Write-Output "'Set-ExecutionPolicy RemoteSigned -scope CurrentUser'"
    break
}

$DOWNLOAD_URL = (Invoke-WebRequest -Uri 'https://api.github.com/repos/datreeio/datree/releases/latest' -UseBasicParsing | select-string -Pattern 'https://github.com/datreeio/datree/releases/download/\d+\.\d+\.\d+/datree-cli_\d+\.\d+\.\d+_windows_x86_64.zip').Matches.Value
$OUTPUT_BASENAME = "datree-latest"
$OUTPUT_BASENAME_WITH_POSTFIX = "$OUTPUT_BASENAME.zip"

Write-Host 'Installing Datree...'
Write-Host ''
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $OUTPUT_BASENAME_WITH_POSTFIX -UseBasicParsing
Write-Host "[V] Downloaded Datree" -ForegroundColor DarkGreen

Expand-Archive -Path $OUTPUT_BASENAME_WITH_POSTFIX -DestinationPath $OUTPUT_BASENAME -Force | Out-Null

$localAppDataPath = $env:LOCALAPPDATA
$datreePath = Join-Path "$localAppDataPath" 'datree'
New-Item -ItemType Directory -Force -Path $datreePath

Copy-Item -Path "$OUTPUT_BASENAME/*" -Destination "$datreePath" -PassThru -Force | Out-Null

Remove-Item -Recurse $OUTPUT_BASENAME
Remove-Item $OUTPUT_BASENAME_WITH_POSTFIX

$dotDatreePath = "$home/.datree"
mkdir -Force $dotDatreePath | Out-Null
$k8sDemoPath = Join-Path "$dotDatreePath" "k8s-demo.yaml"

Invoke-WebRequest -Uri "https://get.datree.io/k8s-demo.yaml" -OutFile $k8sDemoPath -UseBasicParsing

Write-Host "[V] Finished Installation" -ForegroundColor DarkGreen
Write-Host ""
Write-Host "To run datree globally, please run the following command as administrator:" -ForegroundColor Cyan
Write-Host ""
Write-Host "setx PATH `$env:path;$datreePath -m"
Write-Host ""
Write-Host "    Usage: datree test `$home/.datree/k8s-demo.yaml" -ForegroundColor DarkGreen
Write-Host ""
Write-Host "    Using Helm? => https://hub.datree.io/helm-plugin"
Write-Host ""
