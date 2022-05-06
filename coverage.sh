#! /bin/sh

export ZINC_FIRST_ADMIN_USER=admin  
export ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123

go test -v ./... -race -covermode=atomic -coverprofile=coverage.out

# make sure to set CODECOV_TOKEN env variable before doing this
# codecov -f coverage.out
# or 
# bash <(curl -s https://codecov.io/bash) 


# Example setup https://github.com/lluuiissoo/go-testcoverage/blob/main/.github/workflows/ci.yml

# enable threshold
COVERAGE_THRESHOLD=20

totalCoverage=`go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]'`

echo "Total Coverage is $totalCoverage %"

diff=$(echo "$totalCoverage < $COVERAGE_THRESHOLD" | bc)

if [ $diff -eq 1 ]; then
    echo "Coverage is below threshold of $COVERAGE_THRESHOLD %"
    exit 1
else
    echo "Coverage is above threshold of $COVERAGE_THRESHOLD %"
    exit 0
fi





