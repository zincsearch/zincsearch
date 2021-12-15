#!/bin/bash

cd web
npm run build
cd ..

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o zinc cmd/zinc/main.go


