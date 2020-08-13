#!/bin/bash
cd "$(dirname "$0")"

# Generate swagger doc!
$GOBIN/swag init -g cmd/server/main.go
