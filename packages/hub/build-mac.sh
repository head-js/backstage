#!/bin/sh
set -e
cd "$(dirname "$0")"
go build -o .dist/backstage-hub ./cmd/hub
