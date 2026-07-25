param(
    [switch]$Force
)

$ErrorActionPreference = "Stop"

$binaryPath = Join-Path $PSScriptRoot "yt-dlp.exe"
$downloadUrl = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"

function Get-LocalVersion {
    if (-not (Test-Path $binaryPath)) {
        return $null
    }
    try {
        return (& $binaryPath --version) 2>$null
    } catch {
        return $null
    }
}

function Install-YtDlp {
    Write-Host "Telechargement de yt-dlp.exe..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $downloadUrl -OutFile $binaryPath
    Write-Host "yt-dlp.exe installe dans $binaryPath" -ForegroundColor Green
}

function Update-YtDlp {
    Write-Host "Verification des mises a jour..." -ForegroundColor Cyan
    & $binaryPath -U
}

function Test-DenoInstalled {
    $deno = Get-Command deno -ErrorAction SilentlyContinue
    return $null -ne $deno
}

function Install-Deno {
    Write-Host "`nDeno non trouve. Installation via winget..." -ForegroundColor Cyan
    $winget = Get-Command winget -ErrorAction SilentlyContinue
    if (-not $winget) {
        Write-Warning "winget introuvable. Installe Deno manuellement: https://deno.com/manual/getting_started/installation"
        return
    }
    winget install DenoLand.Deno --accept-source-agreements --accept-package-agreements
    Write-Host "Deno installe. Redemarre ton terminal pour que le PATH soit pris en compte." -ForegroundColor Green
}

$currentVersion = Get-LocalVersion

if (-not $currentVersion -or $Force) {
    Install-YtDlp
    $currentVersion = Get-LocalVersion
} else {
    Write-Host "yt-dlp.exe deja present (version $currentVersion)" -ForegroundColor Yellow
    Update-YtDlp
}

if (-not (Test-Path $binaryPath)) {
    Write-Error "Echec de l'installation de yt-dlp.exe"
    exit 1
}

if (Test-DenoInstalled) {
    Write-Host "Deno deja installe (runtime JS pour yt-dlp OK)" -ForegroundColor Yellow
} else {
    Install-Deno
}

$finalVersion = Get-LocalVersion
Write-Host "`nyt-dlp pret. Version: $finalVersion" -ForegroundColor Green