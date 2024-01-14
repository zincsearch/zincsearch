#! /bin/sh

if ! command -v golangci-lint &> /dev/null
then
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $PWD/ v1.55.2
    ./golangci-lint run
    rm golangci-lint
else
    golangci-lint run
fi

rc=$?
if [ $rc -ne 0 ]; then
  echo "golangci-lint failed" >&2
  exit $rc
fi
