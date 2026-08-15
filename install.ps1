# UMD installer - Windows
# Usage:
#   irm https://vegidio.github.io/umd/install.ps1 | iex
#   $env:UMD_INSTALL='cli'; irm https://vegidio.github.io/umd/install.ps1 | iex
#   $env:UMD_VERSION='<tag>'; irm https://vegidio.github.io/umd/install.ps1 | iex
#
# UMD_VERSION defaults to 'latest', which is resolved dynamically from
# the GitHub API (releases/latest) at run time.

#Requires -Version 5.1
$ErrorActionPreference = 'Stop'

$Repo       = 'vegidio/umd'
$Mode       = if ($env:UMD_INSTALL)     { $env:UMD_INSTALL }     else { 'both' }
$Version    = if ($env:UMD_VERSION)     { $env:UMD_VERSION }     else { 'latest' }
$InstallDir = if ($env:UMD_INSTALL_DIR) { $env:UMD_INSTALL_DIR } else { Join-Path $env:LOCALAPPDATA 'Programs\umd' }

function Write-Info($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Warn($msg) { Write-Host "warn: $msg" -ForegroundColor Yellow }
function Write-Fail($msg) { Write-Host "error: $msg" -ForegroundColor Red; throw $msg }

if ($Mode -notin @('cli','gui','both')) {
    Write-Fail "invalid UMD_INSTALL=$Mode (expected: cli, gui, or both)"
}

# WOW64 reports x86 in PROCESSOR_ARCHITECTURE; ARCHITEW6432 holds the real arch when present.
$archEnv = if ($env:PROCESSOR_ARCHITEW6432) { $env:PROCESSOR_ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE }
switch ($archEnv) {
    'AMD64' { $Arch = 'amd64' }
    'ARM64' { $Arch = 'arm64' }
    default { Write-Fail "unsupported architecture: $archEnv" }
}

if ($Version -eq 'latest') {
    Write-Info 'resolving latest version...'
    try {
        $release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -UseBasicParsing
        $Tag = $release.tag_name
    } catch {
        Write-Fail "could not resolve latest version: $($_.Exception.Message)"
    }
    if (-not $Tag) { Write-Fail 'could not parse latest version from GitHub API' }
} else {
    $Tag = $Version
}

Write-Info "installing umd $Tag (windows/$Arch)"

$Tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("umd-install-" + [System.Guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $Tmp -Force | Out-Null

function Get-Asset($asset) {
    $url = "https://github.com/$Repo/releases/download/$Tag/$asset"
    $zip = Join-Path $Tmp $asset
    Write-Info "downloading $asset"
    $oldProgress = $ProgressPreference
    $ProgressPreference = 'SilentlyContinue'   # Invoke-WebRequest is glacial with progress on
    try {
        Invoke-WebRequest -Uri $url -OutFile $zip -UseBasicParsing
    } catch {
        Write-Fail "download failed: $url ($($_.Exception.Message))"
    } finally {
        $ProgressPreference = $oldProgress
    }
    $extractDir = Join-Path $Tmp ([System.IO.Path]::GetFileNameWithoutExtension($asset))
    Expand-Archive -Path $zip -DestinationPath $extractDir -Force
    return $extractDir
}

function Install-Binary($srcPath, $name) {
    if (-not (Test-Path $srcPath)) { Write-Fail "$name not found in archive" }
    $dst = Join-Path $InstallDir $name
    if (Test-Path $dst) { Remove-Item -Path $dst -Force }
    Move-Item -Path $srcPath -Destination $dst -Force
    # Strip Mark-of-the-Web so SmartScreen/Defender don't block the first launch
    Unblock-File -Path $dst -ErrorAction SilentlyContinue
    Write-Info "$name installed at $dst"
    return $dst
}

try {
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    if ($Mode -eq 'cli' -or $Mode -eq 'both') {
        $dir = Get-Asset "umd-cli_windows_${Arch}.zip"
        Install-Binary -srcPath (Join-Path $dir 'umd-dl.exe') -name 'umd-dl.exe' | Out-Null
    }

    if ($Mode -eq 'gui' -or $Mode -eq 'both') {
        $dir = Get-Asset "umd-gui_windows_${Arch}.zip"
        $guiPath = Install-Binary -srcPath (Join-Path $dir 'umd.exe') -name 'umd.exe'

        # Start Menu shortcut for the GUI
        $startMenu = Join-Path $env:APPDATA 'Microsoft\Windows\Start Menu\Programs'
        $lnk = Join-Path $startMenu 'UMD.lnk'
        try {
            $shell = New-Object -ComObject WScript.Shell
            $shortcut = $shell.CreateShortcut($lnk)
            $shortcut.TargetPath = $guiPath
            $shortcut.WorkingDirectory = $InstallDir
            $shortcut.Save()
            Write-Info "Start Menu shortcut created at $lnk"
        } catch {
            Write-Warn "could not create Start Menu shortcut: $($_.Exception.Message)"
        }
    }

    # Add install dir to user PATH (no admin required)
    $userPath = [System.Environment]::GetEnvironmentVariable('Path', 'User')
    if (-not $userPath) { $userPath = '' }
    $segments = $userPath -split ';' | Where-Object { $_ -ne '' }
    if ($segments -notcontains $InstallDir) {
        $newPath = if ($userPath) { "$userPath;$InstallDir" } else { $InstallDir }
        [System.Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
        Write-Info "added $InstallDir to user PATH (open a new terminal to pick it up)"
    }
} finally {
    Remove-Item -Recurse -Force -Path $Tmp -ErrorAction SilentlyContinue
}

Write-Host 'done.' -ForegroundColor Green
