#!/bin/sh

reflex -d none -s -R vendor. -r \.go$ -- go run cmd/zinc/main.go
