Remove-Item 'C:\go' -Recurse -Force

$goDistPath = "$env:TEMP\go$GO_VERSION.windows-amd64.zip"
(New-Object Net.WebClient).DownloadFile("https://dl.google.com/go/go$GO_VERSION.windows-amd64.zip", $goDistPath)

Write-Host "Unpacking..."
7z x $goDistPath -o"$env:SystemDrive\" | Out-Null

go version