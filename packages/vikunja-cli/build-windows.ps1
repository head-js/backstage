$dist = ".dist"
$target = $args[0]

if (-not $target) {
    $target = "build"
}

switch ($target) {
    "build" {
        if (-not (Test-Path $dist)) {
            New-Item -ItemType Directory -Path $dist | Out-Null
        }
        go build -o "$dist/backstage-vikunja.exe"
    }
    "clean" {
        if (Test-Path $dist) {
            Remove-Item -Recurse -Force $dist
        }
    }
    "test" {
        go test ./...
    }
    default {
        Write-Host "Usage: build-windows.ps1 [build|clean|test]"
    }
}
