#!/bin/bash

export ZINC_FIRST_ADMIN_USER="admin"
export ZINC_FIRST_ADMIN_PASSWORD="Complexpass#123" 

find ./pkg -name data -type dir|xargs rm -fR
find ./test -name data -type dir|xargs rm -fR

if [[ $1 == "bench" ]]; then
    go test -v -test.run=NONE -test.bench=Bulk -benchmem ./test/benchmark/ -cpuprofile=./tmp/cpu.pprof -memprofile=./tmp/mem.pprof
    # go tool pprof -http=:9999 ./tmp/mem.pprof
else
    go test -v ./... -test.run=$1
fi

find ./pkg -name data -type dir|xargs rm -fR
find ./test -name data -type dir|xargs rm -fR
