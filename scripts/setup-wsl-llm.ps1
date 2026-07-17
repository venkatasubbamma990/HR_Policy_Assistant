#Requires -RunAsAdministrator
<#
  Forwards Windows port 12345 -> LM Studio on 127.0.0.1:1234 so WSL can reach it.
  Run once from PowerShell (as Administrator) while LM Studio is NOT bound to 12345.

  Usage:
    powershell -ExecutionPolicy Bypass -File scripts/setup-wsl-llm.ps1
#>

$ListenPort = if ($env:HR_WSL_LLM_PORT) { [int]$env:HR_WSL_LLM_PORT } else { 12345 }
$TargetPort = 1234

Write-Host "Setting up WSL -> LM Studio port forward on port $ListenPort ..."

netsh interface portproxy delete v4tov4 listenaddress=0.0.0.0 listenport=$ListenPort 2>$null | Out-Null
netsh interface portproxy add v4tov4 listenaddress=0.0.0.0 listenport=$ListenPort connectaddress=127.0.0.1 connectport=$TargetPort | Out-Null

$ruleName = "HR Policy LM Studio WSL ($ListenPort)"
netsh advfirewall firewall delete rule name="$ruleName" 2>$null | Out-Null
netsh advfirewall firewall add rule name="$ruleName" dir=in action=allow protocol=TCP localport=$ListenPort | Out-Null

Write-Host ""
Write-Host "Done. Port forward active:"
Write-Host "  0.0.0.0:$ListenPort -> 127.0.0.1:$TargetPort"
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Start LM Studio server on port $TargetPort (Developer -> Start Server)"
Write-Host "  2. From WSL, run: make check-llm"
Write-Host "  3. Then: make ask Q=`"Can I take PL during notice?`""
Write-Host ""
Write-Host "The Go app auto-rewrites OPENAI_BASE_URL to http://<wsl-gateway>:$ListenPort/v1 when running in WSL."
