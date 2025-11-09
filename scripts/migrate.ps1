param(
    [ValidateSet('up','down','force','version','drop','goto')]
    [string]$Action = 'up',
    [int]$Steps = 0,
    [int]$Version,
    [string]$DatabaseUrl = $env:PG_URL,
    [string]$MigrationsPath = "migrations"
)

function Get-MigratePath {
    $cmd = Get-Command migrate -ErrorAction SilentlyContinue
    if ($cmd) {
        return $cmd.Path
    }

    $go = Get-Command go -ErrorAction SilentlyContinue
    if ($go) {
        $gopath = & go env GOPATH
        if (-not $gopath) { $gopath = "$env:USERPROFILE\go" }
        $candidate = Join-Path (Join-Path $gopath "bin") "migrate.exe"
        if (Test-Path $candidate) {
            return $candidate
        }

        Write-Host "Installing golang-migrate CLI via 'go install'..."
        & go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        if ($LASTEXITCODE -eq 0 -and (Test-Path $candidate)) {
            return $candidate
        }
    }

    Write-Error "migrate CLI not found. Install with: go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
}

if (-not $DatabaseUrl) {
    Write-Error "Database URL is not set. Provide -DatabaseUrl or set PG_URL environment variable."
    exit 1
}

$migrateExe = Get-MigratePath
$commonArgs = @('-database', $DatabaseUrl, '-path', $MigrationsPath)

switch ($Action) {
    'up' {
        if ($Steps -gt 0) { $args = $commonArgs + @('up', $Steps) } else { $args = $commonArgs + @('up') }
    }
    'down' {
        if ($Steps -gt 0) { $args = $commonArgs + @('down', $Steps) } else { $args = $commonArgs + @('down') }
    }
    'version' {
        $args = $commonArgs + @('version')
    }
    'force' {
        if (-not $PSBoundParameters.ContainsKey('Version')) {
            Write-Error "Specify -Version for 'force' action."
            exit 1
        }
        $args = $commonArgs + @('force', $Version)
    }
    'drop' {
        $args = $commonArgs + @('drop', '-f')
    }
    'goto' {
        if (-not $PSBoundParameters.ContainsKey('Version')) {
            Write-Error "Specify -Version for 'goto' action."
            exit 1
        }
        $args = $commonArgs + @('goto', $Version)
    }
    default {
        Write-Error "Unknown action: $Action"
        exit 1
    }
}

& $migrateExe @args
exit $LASTEXITCODE


