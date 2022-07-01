#!/bin/sh
# author xiaojun207
# eg. : sh build-image.sh

VERSION=`git describe --tags --always`
BUILD_DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p-GMT'`
COMMIT_HASH=`git rev-parse HEAD`
IMAGE="zinclabs/zinc"

# build image
docker buildx build \
  --build-arg VERSION="v${VERSION}" \
  --build-arg COMMIT_HASH="${COMMIT_HASH}" \
  --build-arg BUILD_DATE="${BUILD_DATE}" \
  --tag "$IMAGE:${VERSION}" \
  --tag "$IMAGE:latest" \
  . -f Dockerfile

# push to image rep
# docker push "$IMAGE:${VERSION}"
# docker push "$IMAGE:latest"
