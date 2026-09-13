#! /bin/sh

cd "$(dirname "$0")/.." || exit 1

if ! command -v golangci-lint &> /dev/null
then
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $PWD/ v2.13.2
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
