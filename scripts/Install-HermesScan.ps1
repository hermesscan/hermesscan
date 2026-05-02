<#
.SYNOPSIS
Installs a local HermesScan executable into a user-selected tools directory.

.DESCRIPTION
Copies a built HermesScan executable to a target directory and optionally prints
instructions for adding that directory to the user PATH. This script is intended
for Windows PowerShell 5.1 and newer.

.PARAMETER SourcePath
Path to the built hermesscan.exe file.

.PARAMETER DestinationDirectory
Directory where hermesscan.exe should be copied.

.PARAMETER AddToUserPath
When specified, adds DestinationDirectory to the current user's PATH if missing.

.EXAMPLE
.\scripts\Install-HermesScan.ps1 -SourcePath .\hermesscan.exe -DestinationDirectory "$env:USERPROFILE\bin" -WhatIf

.EXAMPLE
.\scripts\Install-HermesScan.ps1 -SourcePath .\hermesscan.exe -DestinationDirectory "$env:USERPROFILE\bin" -AddToUserPath

.INPUTS
None.

.OUTPUTS
System.IO.FileInfo

.NOTES
Requires Windows PowerShell 5.1 or newer. Restart your shell after changing PATH.
#>
[CmdletBinding(SupportsShouldProcess = $true)]
param(
    [Parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [string]$SourcePath = '.\hermesscan.exe',

    [Parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [string]$DestinationDirectory = (Join-Path -Path $env:USERPROFILE -ChildPath 'bin'),

    [Parameter(Mandatory = $false)]
    [switch]$AddToUserPath
)

begin {
    Set-StrictMode -Version 2.0
}

process {
    try {
        $resolvedSource = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($SourcePath)
        if (-not (Test-Path -LiteralPath $resolvedSource -PathType Leaf)) {
            throw "Source executable was not found: $resolvedSource"
        }

        $resolvedDestination = $ExecutionContext.SessionState.Path.GetUnresolvedProviderPathFromPSPath($DestinationDirectory)
        if ($PSCmdlet.ShouldProcess($resolvedDestination, 'Create destination directory')) {
            New-Item -ItemType Directory -Path $resolvedDestination -Force | Out-Null
        }

        $targetPath = Join-Path -Path $resolvedDestination -ChildPath 'hermesscan.exe'
        if ($PSCmdlet.ShouldProcess($targetPath, 'Copy HermesScan executable')) {
            Copy-Item -LiteralPath $resolvedSource -Destination $targetPath -Force
        }

        if ($AddToUserPath) {
            $currentPath = [Environment]::GetEnvironmentVariable('Path', 'User')
            $parts = @()
            if (-not [string]::IsNullOrWhiteSpace($currentPath)) {
                $parts = $currentPath -split ';' | Where-Object { -not [string]::IsNullOrWhiteSpace($_) }
            }

            $alreadyPresent = $false
            foreach ($part in $parts) {
                if ([string]::Equals($part.TrimEnd('\'), $resolvedDestination.TrimEnd('\'), [StringComparison]::OrdinalIgnoreCase)) {
                    $alreadyPresent = $true
                    break
                }
            }

            if (-not $alreadyPresent) {
                $newPath = if ([string]::IsNullOrWhiteSpace($currentPath)) { $resolvedDestination } else { "$currentPath;$resolvedDestination" }
                if ($PSCmdlet.ShouldProcess('User PATH', "Add $resolvedDestination")) {
                    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
                    Write-Verbose 'User PATH was updated. Restart the shell to pick up the change.'
                }
            }
        }

        Get-Item -LiteralPath $targetPath
    }
    catch {
        Write-Error -ErrorRecord $_
    }
}
