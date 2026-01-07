Write-Host "[air] running go generate..."
go generate ./internal/web/handlers
if ($LASTEXITCODE -ne 0) { exit 1 }

# ---- WASM build ----
$dist = ".\internal\web\handlers\assets_dist"
$jsOut = Join-Path $dist "js"
New-Item -ItemType Directory -Force -Path $jsOut | Out-Null

Write-Host "[air] building wasm netutil..."
$env:GOOS="js"
$env:GOARCH="wasm"

go build -o (Join-Path $jsOut "netutil.wasm") .\cmd\wasm
if ($LASTEXITCODE -ne 0) { exit 1 }

Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

Write-Host "[air] copying wasm_exec.js..."
$goroot = (go env GOROOT)
$wasmExec1 = Join-Path $goroot "lib\wasm\wasm_exec.js"
$wasmExec2 = Join-Path $goroot "misc\wasm\wasm_exec.js"
if (Test-Path $wasmExec1) {
  Copy-Item $wasmExec1 (Join-Path $jsOut "wasm_exec.js") -Force
} else {
  Copy-Item $wasmExec2 (Join-Path $jsOut "wasm_exec.js") -Force
}
# --------------------

Write-Host "[air] building binary..."
go build -o ./tmp/kea-web.exe ./cmd/kea-web
if ($LASTEXITCODE -ne 0) { exit 1 }