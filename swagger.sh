#!/bin/sh

if ! command -v swag &> /dev/null
then
    echo "swag could not be found. Make sure it is installed and available on the environment"
    echo "please refer: https://github.com/swaggo/swag"
    exit
fi

echo "Running swag command"
swag init -g cmd/zincsearch/main.go
