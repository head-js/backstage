# make-windows.ps1
$DIST = ".dist"
$VERSION = "0.0.1"
$LDFLAGS = "-X main.cliVersion=$VERSION"

$target = $args[0]

if (-not $target) {
    $target = "build"
}

switch ($target) {
    "build" {
        Write-Host "Building Windows version..."
        if (-not (Test-Path $DIST)) {
            New-Item -ItemType Directory -Path $DIST | Out-Null
        }
        go build -ldflags $LDFLAGS -o "$DIST/backstage-gitea.exe"
        Write-Host "Build complete: $DIST/backstage-gitea.exe"
    }
    "clean" {
        Write-Host "Cleaning build artifacts..."
        if (Test-Path $DIST) {
            Remove-Item -Recurse -Force $DIST
        }
        Write-Host "Clean complete."
    }
    "test" {
        Write-Host "Running tests..."
        go test ./...
    }
    "version" {
        Write-Host "Version: $VERSION"
    }
    default {
        Write-Host "Usage: make-windows.ps1 [build|clean|test|version]"
        Write-Host "  build   - Build Windows version"
        Write-Host "  clean   - Remove build artifacts"
        Write-Host "  test    - Run tests"
        Write-Host "  version - Print version info"
    }
}
