# Contributing to ZincSearch

## Setting up development environment

### Prerequisite

ZincSearch uses Go (For server) and VueJS (For Web UI)

You must have following installed:

1. Git
2. Go 1.16 + (We recommend go 1.19+)
3. nodejs v14+ and npm v6+

## Building from source code

### Lets clone the repo and get started

```shell
git clone https://github.com/zinclabs/zincsearch
cd zincsearch
```

### Now let's build the UI

```shell
cd web
npm install
npm run build
cd ..
```

Output will be stored in web/dist folder. web/dist will be embedded in ZincSearch binary when ZincSearch go application is built.

It is important that you build the web app every time you make any changes to javascript code as the built code is then embedded in go application.

### Time to build the go application now

Download the dependencies

```shell
go mod tidy # this will download the go libraries used by zincsearch
```

Simple:

```shell
go build -o zincsearch cmd/zincsearch/main.go # will build the ZincSearch binary
```

Advanced

```shell
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X github.com/zinclabs/zincsearch/pkg/meta.Version=${VERSION} -X github.com/zinclabs/zincsearch/pkg/meta.CommitHash=${COMMIT_HASH} -X github.com/zinclabs/zincsearch/pkg/meta.BuildDate=${BUILD_DATE}" -o zincsearch cmd/zincsearch/main.go
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
go mod tidy
ZINC_FIRST_ADMIN_USER=admin ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 go run cmd/zincsearch/main.go
```

This will start the ZincSearch API server on port 4080

environment variables ZINC_FIRST_ADMIN_USER and ZINC_FIRST_ADMIN_PASSWORD are required only first time when ZincSearch is started.

### UI

```shell
cd web
npm install
npm run dev
```

This will start UI server on port 8080

In order for you to effectively use the UI you would want to have the ZincSearch API server running in a separate window that will accept requests from the UI.

## Swagger

The server also exposes a Swagger API endpoint which you can see by visiting the `/swagger/index.html` path. It uses [gin-swagger](https://github.com/swaggo/gin-swagger) to mark API endpoints with comment annotations and [swag](https://github.com/swaggo/swag) to generate the API spec from the annotations to Swagger Documentation 2.0.

If you update the annotations, you need to also regenerate the Swagger documentation by running the `swagger.sh` script located at the base project folder:

````bash
./swagger.sh
2022/05/31 10:18:13 Generate swagger docs....
2022/05/31 10:18:13 Generate general API Info, search dir:./
2022/05/31 10:18:13 Generating auth.LoginRequest
2022/05/31 10:18:13 Generating auth.LoginSuccess
2022/05/31 10:18:13 Generating auth.SimpleUser
2022/05/31 10:18:13 Generating auth.LoginError
2022/05/31 10:18:13 create docs.go at  docs/docs.go
2022/05/31 10:18:13 create swagger.json at  docs/swagger.json
2022/05/31 10:18:13 create swagger.yaml at  docs/swagger.yaml
```

## Build docker image

Make sure that you have [docker](https://docs.docker.com/get-docker/).

Simple build:

```shell
docker build --tag zinc:latest . -f Dockerfile
````

Multi-arch build

In order to build multi-arch builds you will need [buildx](https://docs.docker.com/buildx/working-with-buildx/) installed. You will need to pass the platform flag for the platform that you want to build.

```shell
docker buildx build --platform linux/amd64 --tag zinc:latest-linux-amd64 . -f Dockerfile.hub
```

# Checks in CI pipeline

We check for following in CI pipeline for any pull requests.

1. Unit test code coverage for go code.
    - If code coverage is less than 81% (according to go test) the CI tests will fail.
    - You can test coverage yourself by running `./coverage.sh` 
    - We use codecov for visualizing code coverage of go code, codecov updates coverage for every PR through a comment. It allows you to see missing coverage for any lines.
1. Linting in Javascript for GUI
    - We run eslint for javacript anf any linting failures will result in build failures.
    - You can test for linting failures by running `./lint.sh` in web folder.
    - You can also fix automatically fixable linting error by running `npm run lint-autofix`


## How to contribute code

1. Fork the repository on github (e.g. awesomedev/zincsearch)
1. Clone the repo from the forked repository ( e.g. awesomedev/zincsearch) to your machine.
1. create a new branch locally. 
1. Make the changes to code.
1. Push the code to your repo.
1. Create a PR
1. Make sure that the automatic CI checks pass for your PR.
