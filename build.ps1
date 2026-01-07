Write-Host "[air] running go generate..."
go generate ./internal/web/handlers
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host "[air] building binary..."
go build -o ./tmp/kea-web.exe ./cmd/kea-web
if ($LASTEXITCODE -ne 0) { exit 1 }