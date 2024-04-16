$processes = Get-Process

$cpu_load = ($processes | Measure-Object -Property CPU -Sum).Sum / (Get-WmiObject Win32_ComputerSystem).NumberOfLogicalProcessors

$threshold = 90

if ($cpu_load -gt $threshold) {
    Write-Host "CPU load is high ($cpu_load%). Taking action to reduce load."

    $logDirectory = "C:\Logs"
    if (-not (Test-Path -Path $logDirectory)) {
        New-Item -Path $logDirectory -ItemType Directory | Out-Null
    }

    $logPath = Join-Path -Path $logDirectory -ChildPath "automated_response.log"
    Add-Content -Path $logPath -Value "$(Get-Date): CPU load is high ($cpu_load%). Action taken."
} else {
    Write-Host "CPU load is normal ($cpu_load%). No action required."
}
