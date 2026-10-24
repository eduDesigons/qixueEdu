# Publish admin / web / api. Static config: qixue-api/.env STATIC_PUBLISH.
# Encoding: UTF-8 with BOM (required for Windows PowerShell 5).

$ErrorActionPreference = 'Stop'
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

$ApiGoRoot = $PSScriptRoot
Set-Location $ApiGoRoot

function Fail([string]$msg) {
    Write-Host "[错误] $msg" -ForegroundColor Red
    Read-Host '按回车退出'
    exit 1
}

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Fail '找不到 go，请把 Go 加入 PATH'
}
if (-not (Get-Command npm -ErrorAction SilentlyContinue)) {
    Fail '找不到 npm，请先安装 Node.js'
}
if (-not (Test-Path (Join-Path $ApiGoRoot '.env'))) {
    Fail '缺少 qixue-api\.env'
}

Write-Host ''
Write-Host '  Publish (static + api)'
Write-Host '  --------------------------------'
Write-Host '  [1] admin'
Write-Host '  [2] web'
Write-Host '  [3] admin + web'
Write-Host '  [4] api（Go 交叉编译 → deploy/release）'
Write-Host '  [5] 全部（admin + web + api）'
Write-Host '  [0] 取消'
Write-Host ''
$choice = Read-Host '请选择'

function Publish([string[]]$targets) {
    $goArgs = @('run', './cmd/publish') + $targets
    & go @goArgs
    if ($LASTEXITCODE -ne 0) { throw "发布失败，退出码 $LASTEXITCODE" }
}

try {
    switch ($choice) {
        '1' { Publish @('admin') }
        '2' { Publish @('web') }
        '3' { Publish @('admin', 'web') }
        '4' { Publish @('api') }
        '5' { Publish @('admin', 'web', 'api') }
        Default {
            Write-Host '已取消'
            Read-Host '按回车退出'
            exit 0
        }
    }
    Write-Host ''
    Write-Host '完成' -ForegroundColor Green
} catch {
    Write-Host ''
    Write-Host "[失败] $($_.Exception.Message)" -ForegroundColor Red
    Read-Host '按回车退出'
    exit 1
}

Read-Host '按回车退出'
