#!/bin/bash

cd web
npm run build
cd ..

go install github.com/rakyll/statik@latest
statik -src=./web/dist

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o zinc


