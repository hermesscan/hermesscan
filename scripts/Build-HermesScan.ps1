<#
.SYNOPSIS
Builds HermesScan native binaries.

.DESCRIPTION
Builds the HermesScan Go CLI for the current platform by default. When -AllTargets
is specified, builds Windows, Linux, and macOS amd64/arm64 binaries into the dist folder.
The script also writes SHA256 checksum files for release artifacts.

.PARAMETER OutputPath
Directory where built binaries are written.

.PARAMETER AllTargets
Build Windows, Linux, and macOS binaries.

.PARAMETER Version
Version string injected into the HermesScan binary.

.PARAMETER SkipChecksum
Do not create SHA256 checksum files.

.EXAMPLE
.\scripts\Build-HermesScan.ps1 -WhatIf

Shows what would be built.

.EXAMPLE
.\scripts\Build-HermesScan.ps1 -AllTargets -Version 0.7.0

Builds cross-platform binaries into .\dist and writes checksums.

.INPUTS
None.

.OUTPUTS
System.IO.FileInfo

.NOTES
Requires Go to be installed and available on PATH.
Compatible with Windows PowerShell 5.1.
#>
[CmdletBinding(SupportsShouldProcess = $true)]
param(
    [Parameter()]
    [string]$OutputPath = '.\dist',

    [Parameter()]
    [switch]$AllTargets,

    [Parameter()]
    [string]$Version = '0.7.0',

    [Parameter()]
    [switch]$SkipChecksum
)

Set-StrictMode -Version 2.0

function New-HermesScanChecksum {
    [CmdletBinding(SupportsShouldProcess = $true)]
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path
    )

    $checksumPath = "$($Path).sha256"
    if ($PSCmdlet.ShouldProcess($checksumPath, 'Write SHA256 checksum')) {
        $hash = Get-FileHash -LiteralPath $Path -Algorithm SHA256
        $line = '{0}  {1}' -f $hash.Hash.ToLowerInvariant(), (Split-Path -Leaf $Path)
        Set-Content -LiteralPath $checksumPath -Value $line -Encoding ASCII
        Get-Item -LiteralPath $checksumPath
    }
}

function Invoke-GoBuildTarget {
    [CmdletBinding(SupportsShouldProcess = $true)]
    param(
        [Parameter(Mandatory = $true)]
        [string]$GoOs,

        [Parameter(Mandatory = $true)]
        [string]$GoArch,

        [Parameter(Mandatory = $true)]
        [string]$OutputFile,

        [Parameter(Mandatory = $true)]
        [string]$BuildVersion
    )

    if ($PSCmdlet.ShouldProcess($OutputFile, 'Build HermesScan binary')) {
        $oldGoOs = $env:GOOS
        $oldGoArch = $env:GOARCH
        try {
            $env:GOOS = $GoOs
            $env:GOARCH = $GoArch
            & go build -ldflags "-X main.version=$BuildVersion" -o $OutputFile .\cmd\hermesscan
            if ($LASTEXITCODE -ne 0) {
                throw "go build failed for $($GoOs)/$($GoArch)."
            }
            Get-Item -LiteralPath $OutputFile
            if (-not $SkipChecksum) {
                New-HermesScanChecksum -Path $OutputFile
            }
        }
        finally {
            if ($null -eq $oldGoOs) {
                Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
            }
            else {
                $env:GOOS = $oldGoOs
            }

            if ($null -eq $oldGoArch) {
                Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
            }
            else {
                $env:GOARCH = $oldGoArch
            }
        }
    }
}

try {
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        throw 'Go was not found on PATH.'
    }

    if ($PSCmdlet.ShouldProcess($OutputPath, 'Create output directory')) {
        New-Item -Path $OutputPath -ItemType Directory -Force | Out-Null
    }

    if ($AllTargets) {
        Invoke-GoBuildTarget -GoOs windows -GoArch amd64 -OutputFile (Join-Path $OutputPath 'hermesscan-windows-amd64.exe') -BuildVersion $Version
        Invoke-GoBuildTarget -GoOs windows -GoArch arm64 -OutputFile (Join-Path $OutputPath 'hermesscan-windows-arm64.exe') -BuildVersion $Version
        Invoke-GoBuildTarget -GoOs linux -GoArch amd64 -OutputFile (Join-Path $OutputPath 'hermesscan-linux-amd64') -BuildVersion $Version
        Invoke-GoBuildTarget -GoOs linux -GoArch arm64 -OutputFile (Join-Path $OutputPath 'hermesscan-linux-arm64') -BuildVersion $Version
        Invoke-GoBuildTarget -GoOs darwin -GoArch amd64 -OutputFile (Join-Path $OutputPath 'hermesscan-darwin-amd64') -BuildVersion $Version
        Invoke-GoBuildTarget -GoOs darwin -GoArch arm64 -OutputFile (Join-Path $OutputPath 'hermesscan-darwin-arm64') -BuildVersion $Version
    }
    else {
        $name = 'hermesscan'
        if ($env:OS -eq 'Windows_NT') {
            $name = 'hermesscan.exe'
        }
        $outputFile = Join-Path $OutputPath $name
        Invoke-GoBuildTarget -GoOs (go env GOOS) -GoArch (go env GOARCH) -OutputFile $outputFile -BuildVersion $Version
    }
}
catch {
    Write-Error -ErrorRecord $_
}
