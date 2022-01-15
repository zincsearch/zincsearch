#!/bin/bash

rm -fR test/data

FIRST_ADMIN_USER="admin" FIRST_ADMIN_PASSWORD="Complexpass#123" go test -v ./...
