# Contributing to Zinc
## Setting up development environment

### Prerequisite

Zinc uses Go (For server) and VueJS (For Web UI)

You must have follwing installed:

1. Git
2. Go 1.16 +
3. nodejs v14+ and npm v6+

## Building from source code

### Lets clone the repo and get started

```shell
git clone https://github.com/prabhatsharma/zinc
cd zinc
```

### Now let's build the UI

```shell
cd web
npm install
npm run build
cd ..
```

Output will be stored in web/dist folder. web/dist will be embedded in zinc binary when zinc go application is built. 

It is important that you build the web app every time you make any changes to javascript code as the built code is then embedded in go application.

### Time to build the go application now

Download the dependencies

```shell
go get -d -v # this will download the go libraries used by zinc
```

Simple:
```shell
go build -o zinc cmd/zinc/main.go # will build the zinc binary
```

Advanced

```shell
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X github.com/prabhatsharma/zinc/pkg/meta/v1.Version=${VERSION} -X github.com/prabhatsharma/zinc/pkg/meta/v1.CommitHash=${COMMIT_HASH} -X github.com/prabhatsharma/zinc/pkg/meta/v1.BuildDate=${BUILD_DATE}" -o zinc cmd/zinc/main.go
```

Setting GOOS and GOARCH appropriately allows for cross platform compilation. Check [Official docs](https://go.dev/doc/install/source#environment) for all possible values and combinations. This [gist](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63) is also great.


Setting CGO_ENABLED=0 allows for static linking, which results in a single binary output that has no dependencies.

setting up ldflags allows for passing values like version number to the binary at build time instead of hardcoding the value in cource code. Generally the version number is set by the CI pipeline during build by using the git tag.

## Developing

Once you have the source code cloned you can start development.

There are 2 areas of development. 

1. UI
1. Server


### Server

```shell
go get -d -v
ZINC_FIRST_ADMIN_USER=admin ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 go run cmd/zinc/main.go
```

This will start the Zinc API server on port 4080

environment variables ZINC_FIRST_ADMIN_USER and ZINC_FIRST_ADMIN_PASSWORD are required only first time when zinc is started.

### UI

```shell
cd web
npm install
npm run serve
```
This will start UI server on port 8080

In order for you to effectively use the UI you would wnat to have the Zinc API server running in a seperate window that will accept requests from the UI.


## Build docker image

Make sure that you have [docker](https://docs.docker.com/get-docker/). 

Simple build:

```shell
docker build --tag zinc:latest . -f Dockerfile.hub
```
Multi-arch build

In order to build multi-srach builds you will need [buildx](https://docs.docker.com/buildx/working-with-buildx/) installed. You will need to pass the platform flag for the platform that you want to build.

```shell
docker buildx build --platform linux/amd64 --tag zinc:latest-linux-amd64 . -f Dockerfile.hub
```
