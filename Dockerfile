# syntax=docker/dockerfile:experimental
############################
# STEP 1 build web dist
############################
FROM node:18.18.2-slim as webBuilder
WORKDIR /web
COPY ./web /web/

RUN npm install
RUN npm run build

############################
# STEP 2 build executable binary
############################
# FROM golang:alpine AS builder
FROM public.ecr.aws/docker/library/golang:1.21 as builder
ARG VERSION
ARG COMMIT_HASH
ARG BUILD_DATE

RUN update-ca-certificates
# RUN apk update && apk add --no-cache git
# Create zincsearch user.
ENV USER=zincsearch
ENV GROUP=zincsearch
ENV UID=10001
ENV GID=10001
# See https://stackoverflow.com/a/55757473/12429735RUN
RUN groupadd --gid "${GID}" "${GROUP}"
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    --gid "${GID}" \
    "${USER}"
# Create default directories for persistent ZincSearch data used in final build stage.
# It follows the Linux filesystem hierarchy pattern
# https://tldp.org/LDP/Linux-Filesystem-Hierarchy/html/var.html
RUN mkdir -p /var/lib/zincsearch /data && chown zincsearch:zincsearch /var/lib/zincsearch /data
WORKDIR $GOPATH/src/github.com/zincsearch/zincsearch/
COPY . .
COPY --from=webBuilder /web/dist web/dist

# Fetch dependencies.
# Using go get.
RUN go mod tidy

ENV VERSION=$VERSION
ENV COMMIT_HASH=$COMMIT_HASH
ENV BUILD_DATE=$BUILD_DATE

RUN CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/zincsearch/zincsearch/pkg/meta.Version=${VERSION} -X github.com/zincsearch/zincsearch/pkg/meta.CommitHash=${COMMIT_HASH} -X github.com/zincsearch/zincsearch/pkg/meta.BuildDate=${BUILD_DATE}" -o zincsearch cmd/zincsearch/main.go
############################
# STEP 3 build a small image
############################
# FROM public.ecr.aws/lts/ubuntu:latest
FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the ssl certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy our static executable.
COPY --from=builder  /go/src/github.com/zincsearch/zincsearch/zincsearch /go/bin/zincsearch

# Create directories that can be used to keep ZincSearch data persistent along with host source or named volumes
COPY --from=builder --chown=zincsearch:zincsearch /var/lib/zincsearch /var/lib/zincsearch
COPY --from=builder --chown=zincsearch:zincsearch /data /data

# Use an unprivileged user.
USER zincsearch:zincsearch
# Port on which the service will be exposed.
EXPOSE 4080
# Run the zincsearch binary.
ENTRYPOINT ["/go/bin/zincsearch"]
