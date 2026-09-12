#!/bin/sh

cd "$(dirname "$0")/.." || exit 1

reflex -d none -s -R vendor. -r \.go$ -- go run cmd/zincsearch/main.go
