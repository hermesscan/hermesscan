<#
.SYNOPSIS
Updates the local Scoop manifest for a HermesScan release.

.DESCRIPTION
Downloads the release checksum manifest, extracts Windows amd64 and arm64 hashes,
and updates packaging/scoop/hermesscan.json with release URLs and hashes.

.PARAMETER Version
HermesScan release version without the leading v.

.PARAMETER ManifestPath
Path to the Scoop manifest to update.

.PARAMETER Repository
GitHub repository that hosts HermesScan releases.

.EXAMPLE
.\scripts\Update-ScoopManifest.ps1 -Version 0.8.0

.NOTES
Compatible with Windows PowerShell 5.1.
#>
[CmdletBinding(SupportsShouldProcess = $true)]
param(
    [Parameter(Mandatory = $true)]
    [ValidatePattern('^\d+\.\d+\.\d+$')]
    [string] $Version,

    [Parameter()]
    [string] $ManifestPath = '.\packaging\scoop\hermesscan.json',

    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string] $Repository = 'hermesscan/hermesscan'
)

Set-StrictMode -Version 2.0

function Get-ChecksumEntry {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)]
        [string[]] $Lines,

        [Parameter(Mandatory = $true)]
        [string] $Asset
    )

    $escapedAsset = [regex]::Escape($Asset)
    $line = $Lines | Where-Object { $_ -match "\s+$escapedAsset$" } | Select-Object -First 1
    if (-not $line) {
        throw "No checksum entry found for $Asset."
    }

    $parts = $line.Trim() -split '\s+', 2
    if ($parts.Count -ne 2) {
        throw "Invalid checksum entry for $Asset."
    }

    return $parts[0].ToLowerInvariant()
}

try {
    $resolvedManifestPath = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($ManifestPath)
    $manifestDirectory = Split-Path -Parent $resolvedManifestPath
    if (-not (Test-Path -LiteralPath $manifestDirectory)) {
        throw "Manifest directory does not exist: $manifestDirectory"
    }

    $baseUri = "https://github.com/$Repository/releases/download/v$Version"
    $checksumsUri = "$baseUri/checksums.txt"
    $temporaryChecksums = Join-Path ([System.IO.Path]::GetTempPath()) "hermesscan-checksums-$Version.txt"

    try {
        Invoke-WebRequest -Uri $checksumsUri -OutFile $temporaryChecksums -ErrorAction Stop
        $checksumLines = Get-Content -LiteralPath $temporaryChecksums -ErrorAction Stop
    }
    finally {
        Remove-Item -LiteralPath $temporaryChecksums -Force -ErrorAction SilentlyContinue
    }

    $amd64Asset = 'hermesscan-windows-amd64.exe'
    $arm64Asset = 'hermesscan-windows-arm64.exe'
    $amd64Hash = Get-ChecksumEntry -Lines $checksumLines -Asset $amd64Asset
    $arm64Hash = Get-ChecksumEntry -Lines $checksumLines -Asset $arm64Asset

    $manifest = [ordered]@{
        version      = $Version
        description  = 'Static analyzer for build scripts, CI scripts, and pipeline definitions.'
        homepage     = "https://github.com/$Repository"
        license      = 'MIT'
        architecture = [ordered]@{
            '64bit' = [ordered]@{
                url  = "$baseUri/$amd64Asset#/hermesscan.exe"
                hash = $amd64Hash
            }
            arm64   = [ordered]@{
                url  = "$baseUri/$arm64Asset#/hermesscan.exe"
                hash = $arm64Hash
            }
        }
        bin          = 'hermesscan.exe'
        checkver     = [ordered]@{
            github = "https://github.com/$Repository"
        }
        autoupdate   = [ordered]@{
            architecture = [ordered]@{
                '64bit' = [ordered]@{
                    url = "https://github.com/$Repository/releases/download/v`$version/$amd64Asset#/hermesscan.exe"
                }
                arm64   = [ordered]@{
                    url = "https://github.com/$Repository/releases/download/v`$version/$arm64Asset#/hermesscan.exe"
                }
            }
        }
    }

    $json = $manifest | ConvertTo-Json -Depth 10
    if ($PSCmdlet.ShouldProcess($resolvedManifestPath, "Update Scoop manifest for HermesScan $Version")) {
        Set-Content -LiteralPath $resolvedManifestPath -Value $json -Encoding ASCII
    }

    [pscustomobject]@{
        Manifest = $resolvedManifestPath
        Version  = $Version
        Amd64    = $amd64Hash
        Arm64    = $arm64Hash
    }
}
catch {
    Write-Error -ErrorRecord $_
}
