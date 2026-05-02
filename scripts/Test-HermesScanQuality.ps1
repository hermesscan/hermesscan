<#
.SYNOPSIS
Runs HermesScan Go quality checks.

.DESCRIPTION
Runs gofmt validation, go vet, and go test for the HermesScan repository.

.PARAMETER SkipFormat
Skip gofmt validation.

.PARAMETER SkipVet
Skip go vet.

.PARAMETER SkipTests
Skip go test.

.EXAMPLE
.\scripts\Test-HermesScanQuality.ps1

Runs all checks.

.EXAMPLE
.\scripts\Test-HermesScanQuality.ps1 -SkipVet

Runs format and tests only.

.INPUTS
None.

.OUTPUTS
None.

.NOTES
Requires Go to be installed and available on PATH.
Compatible with Windows PowerShell 5.1.
#>
[CmdletBinding(SupportsShouldProcess = $true)]
param(
    [Parameter()]
    [switch]$SkipFormat,

    [Parameter()]
    [switch]$SkipVet,

    [Parameter()]
    [switch]$SkipTests
)

Set-StrictMode -Version 2.0

try {
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        throw 'Go was not found on PATH.'
    }

    if (-not $SkipFormat) {
        if ($PSCmdlet.ShouldProcess('Go source files', 'Validate gofmt')) {
            $formatOutput = & gofmt -l .
            if ($LASTEXITCODE -ne 0) {
                throw 'gofmt failed.'
            }
            if ($formatOutput) {
                $formatOutput | ForEach-Object { Write-Error "File is not gofmt formatted: $_" }
                throw 'One or more files are not gofmt formatted.'
            }
        }
    }

    if (-not $SkipVet) {
        if ($PSCmdlet.ShouldProcess('Go packages', 'Run go vet')) {
            & go vet .\...
            if ($LASTEXITCODE -ne 0) {
                throw 'go vet failed.'
            }
        }
    }

    if (-not $SkipTests) {
        if ($PSCmdlet.ShouldProcess('Go packages', 'Run go test')) {
            & go test .\...
            if ($LASTEXITCODE -ne 0) {
                throw 'go test failed.'
            }
        }
    }
}
catch {
    Write-Error -ErrorRecord $_
}
