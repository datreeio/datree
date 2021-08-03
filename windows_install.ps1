[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$DOWNLOAD_URL = (Invoke-WebRequest -Uri 'https://api.github.com/repos/datreeio/datree/releases/latest' | select-string -Pattern 'https://github.com/datreeio/datree/releases/download/\d+\.\d+\.\d+/datree-cli_\d+\.\d+\.\d+_windows_x86_64.zip').Matches.Value
$OUTPUT_BASENAME = "datree-latest"
$OUTPUT_BASENAME_WITH_POSTFIX = "$OUTPUT_BASENAME.zip"

Write-Host 'Installing Datree...'
Write-Host ''
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $OUTPUT_BASENAME_WITH_POSTFIX
Write-Host "[V] Downloaded Datree" -ForegroundColor DarkGreen

Expand-Archive -Path $OUTPUT_BASENAME_WITH_POSTFIX -DestinationPath $OUTPUT_BASENAME -Force | Out-Null

$localAppDataPath = $env:LOCALAPPDATA
$datreePath = Join-Path "$localAppDataPath" 'datree'

Copy-Item -Path "$OUTPUT_BASENAME/*" -Destination $datreePath -PassThru -Force | Out-Null

Remove-Item -Recurse $OUTPUT_BASENAME
Remove-Item $OUTPUT_BASENAME_WITH_POSTFIX

$dotDatreePath = "~/.datree"
mkdir -Force $dotDatreePath | Out-Null
$k8sDemoPath =  Join-Path "$dotDatreePath" "k8s-demo.yaml"

Invoke-WebRequest -Uri "https://get.datree.io/k8s-demo.yaml" -OutFile $k8sDemoPath

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
