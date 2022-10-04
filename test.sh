#!/bin/bash

export ZINC_FIRST_ADMIN_USER=admin  
export ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123
export ZINC_WAL_SYNC_INTERVAL=10ms
export ZINC_WAL_REDOLOG_NO_SYNC=true
export ZINC_ENABLE_TEXT_KEYWORD_MAPPING=true
export ZINC_JWT_SECRET=SuperObviousSecret123

find ./pkg -type d -name data |xargs rm -fR
find ./test -type d -name data |xargs rm -fR

if [[ $1 == "bench" ]]; then
    go test -v -test.run=NONE -test.bench=Bulk -benchmem ./test/benchmark/ -cpuprofile=./tmp/cpu.pprof -memprofile=./tmp/mem.pprof
    # go tool pprof -http=:9999 ./tmp/mem.pprof
else
    go test -v ./... -test.run=$1
fi

find ./pkg -type d -name data |xargs rm -fR
find ./test -type d -name data |xargs rm -fR
