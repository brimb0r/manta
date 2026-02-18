#
# Install a terraform-provider-manta release zip for local use.
#
# Usage:
#   .\install.ps1 terraform-provider-manta_0.0.1_windows_amd64.zip
#
# What it does:
#   1. Extracts the binary from the zip
#   2. Places it in the Terraform filesystem_mirror plugin directory
#      so that `terraform init` picks it up locally

param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$ZipFile
)

$ErrorActionPreference = "Stop"

$Hostname = "registry.terraform.io"
$Namespace = "gagno"
$ProviderType = "manta"

# Resolve zip path
if (-not (Test-Path $ZipFile)) {
    $ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
    $Candidate = Join-Path $ScriptDir $ZipFile
    if (Test-Path $Candidate) {
        $ZipFile = $Candidate
    } else {
        Write-Error "Zip file not found: $ZipFile"
    }
}

# Parse version, os, arch from filename: terraform-provider-manta_0.0.1_windows_amd64.zip
$BaseName = [System.IO.Path]::GetFileNameWithoutExtension($ZipFile)
if ($BaseName -match '^terraform-provider-manta_(\d+\.\d+\.\d+)_(\w+)_(\w+)$') {
    $Version = $Matches[1]
    $Os = $Matches[2]
    $Arch = $Matches[3]
} else {
    Write-Error "Could not parse version/os/arch from filename: $ZipFile"
}

# Build the filesystem_mirror plugin directory path
$PluginDir = Join-Path $env:APPDATA "terraform.d\plugins\$Hostname\$Namespace\$ProviderType\$Version\${Os}_${Arch}"

if (-not (Test-Path $PluginDir)) {
    New-Item -ItemType Directory -Path $PluginDir -Force | Out-Null
}

# Extract to a temp directory first, then copy the binary into the plugin dir
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) "terraform-provider-manta-install"
if (Test-Path $TempDir) {
    Remove-Item -Recurse -Force $TempDir
}
New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

Write-Host "Extracting $ZipFile ..."
Expand-Archive -Path $ZipFile -DestinationPath $TempDir -Force

# Copy the binary into the plugin directory
$Binary = Get-ChildItem -Path $TempDir -Filter "terraform-provider-manta*" | Select-Object -First 1
if (-not $Binary) {
    Write-Error "No terraform-provider-manta binary found in zip"
}
Copy-Item -Path $Binary.FullName -Destination $PluginDir -Force

# Clean up temp dir
Remove-Item -Recurse -Force $TempDir

Write-Host "Installed to: $PluginDir\$($Binary.Name)"
Write-Host ""
Write-Host "Done! v$Version is ready for local use."
Write-Host "Run 'terraform init' to pick up the local provider."
