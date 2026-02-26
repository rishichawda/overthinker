$ErrorActionPreference = 'Stop'

$packageName = 'overthink'
$version     = '<VERSION>'
$url64       = "https://github.com/rishichawda/overthinker/releases/download/v$version/overthink_Windows_x86_64.zip"
$checksum64  = '<WINDOWS_AMD64_SHA256>'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
  url64bit      = $url64
  checksum64    = $checksum64
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs
