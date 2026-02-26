$ErrorActionPreference = 'Stop'
$packageName = 'overthink'
$toolsDir    = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

Remove-Item "$toolsDir\overthink.exe" -Force -ErrorAction SilentlyContinue
