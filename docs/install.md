# Installing HermesScan

HermesScan is distributed as a native binary for Windows, Linux, and macOS.

## Windows PowerShell

Download the Windows release binary from GitHub Releases, then place it somewhere on your user PATH.

```powershell
$Version = '0.7.0'
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

## Linux

```bash
VERSION=0.7.0
curl -fsSL -o hermesscan-linux-amd64 \
  "https://github.com/hermesscan/hermesscan/releases/download/v${VERSION}/hermesscan-linux-amd64"
chmod +x ./hermesscan-linux-amd64
sudo install -m 0755 ./hermesscan-linux-amd64 /usr/local/bin/hermesscan
hermesscan version
```

## macOS

Apple Silicon:

```bash
VERSION=0.7.0
curl -fsSL -o hermesscan-darwin-arm64 \
  "https://github.com/hermesscan/hermesscan/releases/download/v${VERSION}/hermesscan-darwin-arm64"
chmod +x ./hermesscan-darwin-arm64
sudo install -m 0755 ./hermesscan-darwin-arm64 /usr/local/bin/hermesscan
hermesscan version
```

Use `hermesscan-darwin-amd64` for Intel Macs and `hermesscan-darwin-arm64` for Apple Silicon.

## Checksums

Release artifacts include `checksums.txt`. After downloading a binary, verify it against the published checksum.

Linux/macOS:

```bash
VERSION=0.7.0
curl -fsSL -O "https://github.com/hermesscan/hermesscan/releases/download/v${VERSION}/checksums.txt"
sha256sum -c checksums.txt --ignore-missing
```

Windows PowerShell:

```powershell
$Version = '0.7.0'
Invoke-WebRequest -Uri "https://github.com/hermesscan/hermesscan/releases/download/v$($Version)/checksums.txt" -OutFile .\checksums.txt
Get-FileHash .\hermesscan.exe -Algorithm SHA256
Get-Content .\checksums.txt
```

Compare the `Get-FileHash` value with the matching line in `checksums.txt`.

## Build from source

```bash
git clone https://github.com/hermesscan/hermesscan.git
cd hermesscan
go test ./...
go build -ldflags "-X main.version=0.7.0" -o hermesscan ./cmd/hermesscan
./hermesscan version
```

Windows PowerShell:

```powershell
go test .\...
go build -ldflags "-X main.version=0.7.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
```
