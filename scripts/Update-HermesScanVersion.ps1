<#
.SYNOPSIS
Updates HermesScan version references across repository text files.

.DESCRIPTION
Finds and replaces HermesScan version references such as 0.6.1 and v0.6.1
with a new version, for example 0.6.1 and v0.6.1.

The script is Windows PowerShell 5.1-compatible and avoids modifying .git,
dist, reports, binaries, archives, and other generated files.

.PARAMETER OldVersion
The old semantic version without the leading "v". Example: 0.6.1

.PARAMETER NewVersion
The new semantic version without the leading "v". Example: 0.6.1

.PARAMETER Path
Repository root path. Defaults to the current directory.

.EXAMPLE
.\scripts\Update-HermesScanVersion.ps1 -OldVersion 0.6.1 -NewVersion 0.6.1 -WhatIf

Shows which files would be updated.

.EXAMPLE
.\scripts\Update-HermesScanVersion.ps1 -OldVersion 0.6.1 -NewVersion 0.6.1

Updates all matching text files.

.INPUTS
None.

.OUTPUTS
System.Management.Automation.PSCustomObject

.NOTES
Compatible with Windows PowerShell 5.1.
#>

[CmdletBinding(SupportsShouldProcess = $true)]
param (
    [Parameter(Mandatory = $true)]
    [ValidatePattern('^\d+\.\d+\.\d+$')]
    [string] $OldVersion,

    [Parameter(Mandatory = $true)]
    [ValidatePattern('^\d+\.\d+\.\d+$')]
    [string] $NewVersion,

    [Parameter()]
    [string] $Path = '.'
)

begin {
    try {
        $root = (Resolve-Path -Path $Path -ErrorAction Stop).ProviderPath

        $oldTag = "v$OldVersion"
        $newTag = "v$NewVersion"

        $excludedDirectoryNames = @(
            '.git',
            'dist',
            'reports',
            'coverage',
            'tmp',
            'node_modules',
            'vendor'
        )

        $excludedExtensions = @(
            '.exe',
            '.dll',
            '.so',
            '.dylib',
            '.zip',
            '.tar',
            '.gz',
            '.7z',
            '.png',
            '.jpg',
            '.jpeg',
            '.gif',
            '.ico',
            '.pdf',
            '.sarif'
        )

        $includedExtensions = @(
            '.md',
            '.txt',
            '.go',
            '.mod',
            '.sum',
            '.yml',
            '.yaml',
            '.json',
            '.ps1',
            '.psm1',
            '.psd1',
            '.sh',
            '.gitignore'
        )

        Write-Verbose "Repository root: $root"
        Write-Verbose "Replacing $OldVersion -> $NewVersion"
        Write-Verbose "Replacing $oldTag -> $newTag"
    }
    catch {
        Write-Error -ErrorRecord $_
        return
    }
}

process {
    try {
        $files = Get-ChildItem -Path $root -Recurse -File -ErrorAction Stop | Where-Object {
            $file = $_

            foreach ($excludedDirectoryName in $excludedDirectoryNames) {
                $pattern = '\\' + [regex]::Escape($excludedDirectoryName) + '(\\|$)'
                if ($file.FullName -match $pattern) {
                    return $false
                }
            }

            if ($excludedExtensions -contains $file.Extension.ToLowerInvariant()) {
                return $false
            }

            if ($includedExtensions -contains $file.Extension.ToLowerInvariant()) {
                return $true
            }

            if ($file.Name -eq 'action.yml' -or $file.Name -eq 'action.yaml') {
                return $true
            }

            return $false
        }

        foreach ($file in $files) {
            $content = Get-Content -LiteralPath $file.FullName -Raw -ErrorAction Stop

            if ($content -notmatch [regex]::Escape($OldVersion) -and $content -notmatch [regex]::Escape($oldTag)) {
                continue
            }

            $updated = $content
            $updated = $updated.Replace($oldTag, $newTag)
            $updated = $updated.Replace($OldVersion, $NewVersion)

            if ($updated -eq $content) {
                continue
            }

            $relativePath = $file.FullName.Substring($root.Length).TrimStart('\', '/')

            if ($PSCmdlet.ShouldProcess($relativePath, "Update version references from $OldVersion to $NewVersion")) {
                Set-Content -LiteralPath $file.FullName -Value $updated -NoNewline -Encoding UTF8 -ErrorAction Stop
            }

            [pscustomobject]@{
                Path       = $relativePath
                OldVersion = $OldVersion
                NewVersion = $NewVersion
                Updated    = $true
            }
        }
    }
    catch {
        Write-Error -ErrorRecord $_
    }
}