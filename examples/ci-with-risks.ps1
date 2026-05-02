[CmdletBinding()]
param()

Start-Job -ScriptBlock { npm install }
Start-Process -FilePath 'docker' -ArgumentList 'compose up -d'
Start-Sleep -Seconds 30
$env:PORT = '5432'
docker build .
