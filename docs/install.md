# Installing HermesScan

HermesScan is distributed as a native binary for Windows, Linux, and macOS.

## Windows PowerShell

Download the Windows release binary from GitHub Releases, then place it somewhere on your user PATH.

```powershell
$Version = '0.9.0'
$Uri = "https://github.com/hermesscan/hermesscan/releases/download/v$($Version)/hermesscan-windows-amd64.exe"
New-Item -Path "$env:USERPROFILE\bin" -ItemType Directory -Force | Out-Null
Invoke-WebRequest -Uri $Uri -OutFile "$env:USERPROFILE\bin\hermesscan.exe"
```

If you already downloaded the binary manually:

```powershell
New-Item -Path "$env:USERPROFILE\bin" -ItemType Directory -Force | Out-Null
Copy-Item -Path .\hermesscan-windows-amd64.exe -Destination "$env:USERPROFILE\bin\hermesscan.exe" -Force
```

Or use the bundled helper from a source checkout:

```powershell
.\scripts\Install-HermesScan.ps1 `
    -SourcePath .\hermesscan.exe `
    -DestinationDirectory "$env:USERPROFILE\bin" `
    -AddToUserPath
```

Restart PowerShell after updating PATH.

```powershell
hermesscan version
hermesscan scan . --summary --no-fail
```

When running from the current directory without PATH installation, PowerShell requires `./` or `.\`:

```powershell
.\hermesscan.exe scan . --summary --no-fail
```

## Scoop

Install from the raw Scoop manifest:

```powershell
scoop install https://raw.githubusercontent.com/hermesscan/hermesscan/main/packaging/scoop/hermesscan.json
hermesscan version
```

The raw manifest is a prototype and currently points at the latest published release assets. See [Scoop packaging plan](scoop-packaging.md) for validation and publishing notes.

## Linux

```bash
VERSION=0.9.0
curl -fsSL -o hermesscan-linux-amd64 \
  "https://github.com/hermesscan/hermesscan/releases/download/v${VERSION}/hermesscan-linux-amd64"
chmod +x ./hermesscan-linux-amd64
sudo install -m 0755 ./hermesscan-linux-amd64 /usr/local/bin/hermesscan
hermesscan version
```

## macOS

Apple Silicon:

```bash
VERSION=0.9.0
curl -fsSL -o hermesscan-darwin-arm64 \
  "https://github.com/hermesscan/hermesscan/releases/download/v${VERSION}/hermesscan-darwin-arm64"
chmod +x ./hermesscan-darwin-arm64
sudo install -m 0755 ./hermesscan-darwin-arm64 /usr/local/bin/hermesscan
hermesscan version
```

Use `hermesscan-darwin-amd64` for Intel Macs and `hermesscan-darwin-arm64` for Apple Silicon.

## Checksums

Release artifacts include `checksums.txt`. After downloading a binary, verify it against the published checksum.

Linux:

```bash
VERSION=0.9.0
ASSET=hermesscan-linux-amd64
BASE_URL="https://github.com/hermesscan/hermesscan/releases/download/v${VERSION}"

curl -fsSLO "${BASE_URL}/${ASSET}"
curl -fsSLO "${BASE_URL}/checksums.txt"
grep "  ${ASSET}$" checksums.txt | sha256sum -c -
chmod +x "./${ASSET}"
sudo install -m 0755 "./${ASSET}" /usr/local/bin/hermesscan
hermesscan version
```

macOS:

```bash
VERSION=0.9.0
ASSET=hermesscan-darwin-arm64
BASE_URL="https://github.com/hermesscan/hermesscan/releases/download/v${VERSION}"

curl -fsSLO "${BASE_URL}/${ASSET}"
curl -fsSLO "${BASE_URL}/checksums.txt"
EXPECTED="$(grep "  ${ASSET}$" checksums.txt | awk '{print $1}')"
ACTUAL="$(shasum -a 256 "./${ASSET}" | awk '{print $1}')"
test "${ACTUAL}" = "${EXPECTED}"
chmod +x "./${ASSET}"
sudo install -m 0755 "./${ASSET}" /usr/local/bin/hermesscan
hermesscan version
```

Windows PowerShell:

```powershell
$Version = '0.9.0'
$Asset = 'hermesscan-windows-amd64.exe'
$BaseUri = "https://github.com/hermesscan/hermesscan/releases/download/v$($Version)"

Invoke-WebRequest -Uri "$BaseUri/$Asset" -OutFile ".\$Asset"
Invoke-WebRequest -Uri "$BaseUri/checksums.txt" -OutFile .\checksums.txt

$line = Get-Content .\checksums.txt |
    Where-Object { $_ -match "\s+$([regex]::Escape($Asset))$" } |
    Select-Object -First 1
if (-not $line) {
    throw "No checksum entry found for $Asset."
}

$expected = ($line -split '\s+')[0].ToLowerInvariant()
$actual = (Get-FileHash -LiteralPath ".\$Asset" -Algorithm SHA256).Hash.ToLowerInvariant()
if ($actual -ne $expected) {
    throw "Checksum mismatch for $Asset."
}

New-Item -Path "$env:USERPROFILE\bin" -ItemType Directory -Force | Out-Null
Copy-Item -Path ".\$Asset" -Destination "$env:USERPROFILE\bin\hermesscan.exe" -Force
& "$env:USERPROFILE\bin\hermesscan.exe" version
```

Use `hermesscan-windows-arm64.exe`, `hermesscan-linux-arm64`, or `hermesscan-darwin-amd64` when that asset matches your platform.

## Build from source

```bash
git clone https://github.com/hermesscan/hermesscan.git
cd hermesscan
go test ./...
go build -ldflags "-X main.version=0.10.0" -o hermesscan ./cmd/hermesscan
./hermesscan version
```

Windows PowerShell:

```powershell
go test .\...
go build -ldflags "-X main.version=0.10.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
```
