#
# Install a terraform-provider-manta release zip for local dev use.
#
# Usage:
#   .\install.ps1 terraform-provider-manta_0.0.1_windows_amd64.zip
#
# What it does:
#   1. Extracts the binary from the zip
#   2. Puts it in a dev_overrides directory
#   3. Ensures your Terraform CLI config has a dev_overrides block
#      pointing to that directory so you can skip `terraform init`

param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$ZipFile
)

$ErrorActionPreference = "Stop"

$ProviderSource = "registry.terraform.io/gagno/manta"

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

# Parse version from filename: terraform-provider-manta_0.0.1_windows_amd64.zip
$BaseName = [System.IO.Path]::GetFileNameWithoutExtension($ZipFile)
if ($BaseName -match '^terraform-provider-manta_(\d+\.\d+\.\d+)') {
    $Version = $Matches[1]
} else {
    Write-Error "Could not parse version from filename: $ZipFile"
}

$DevDir = Join-Path $env:APPDATA "terraform-provider-manta-dev"
$RcFile = Join-Path $env:APPDATA "terraform.rc"

if (-not (Test-Path $DevDir)) {
    New-Item -ItemType Directory -Path $DevDir -Force | Out-Null
}

# Extract
Write-Host "Extracting $ZipFile ..."
Expand-Archive -Path $ZipFile -DestinationPath $DevDir -Force

Write-Host "Binary installed to: $DevDir"

# Configure dev_overrides in Terraform CLI config
$DevDirHcl = $DevDir -replace '\\', '/'

if ((Test-Path $RcFile) -and (Select-String -Path $RcFile -Pattern ([regex]::Escape($ProviderSource)) -Quiet)) {
    Write-Host "dev_overrides already configured in $RcFile"
} elseif ((Test-Path $RcFile) -and (Select-String -Path $RcFile -Pattern "dev_overrides" -Quiet)) {
    Write-Host ""
    Write-Host "WARNING: $RcFile already has a dev_overrides block but does not include"
    Write-Host "  $ProviderSource"
    Write-Host ""
    Write-Host "Add this line inside the dev_overrides block manually:"
    Write-Host "    `"$ProviderSource`" = `"$DevDirHcl`""
} else {
    $Block = @"

provider_installation {
  dev_overrides {
    "$ProviderSource" = "$DevDirHcl"
  }
  direct {}
}
"@
    Add-Content -Path $RcFile -Value $Block
    Write-Host "Wrote dev_overrides to $RcFile"
}

Write-Host ""
Write-Host "Done! v$Version is ready for local use."
Write-Host "Terraform will use the local binary - no 'terraform init' required."
