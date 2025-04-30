if (Test-Path "zincsearch.exe") { Remove-Item -Force .\zincsearch.exe }

Set-Location .\web
try {
    $output = npm run build 2>&1
    if ($LASTEXITCODE -ne 0) { throw $output }
}
catch {
    Write-Host $_ -ForegroundColor Red
    exit $_
}
finally { cd.. }

$version = git describe --tag --always 2>$null
if (-not $version) { $version = "unknown" }
Write-Host "Version: $version"
$buildDate = Get-Date -Format "yyyy-MM-dd_HH:mm:ss:UTCK"
Write-Host "BuildDate: $buildDate"
$commitHash = git rev-parse HEAD 2>$null
Write-Host "CommitHash: $commitHash"
if (-not $commitHash) { $commitHash = "unknown" }
$ldFlags = "-w -s -X 'github.com/zincsearch/zincsearch/pkg/meta.Version=$version'" +
"-X 'github.com/zincsearch/zincsearch/pkg/meta.BuildDate=$buildDate'" +
"-X 'github.com/zincsearch/zincsearch/pkg/meta.CommitHash=$commitHash'"

$output = go build -ldflags $ldFlags -o zincsearch.exe cmd/zincsearch/main.go 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host $output -ForegroundColor Red
    exit $LASTEXITCODE
}