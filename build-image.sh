#!/bin/sh
# author xiaojun207
# eg. : sh build-image.sh 0.2.5 79b9cbf66548b0c9728c1613fa2c64aeb19db81e 20220630 zinclabs/zinc

VERSION="$1"
COMMIT_HASH="$2"
BUILD_DATE="$3"
IMAGE="$4"

# build image
docker buildx build \
  --build-arg VERSION="v${VERSION}" \
  --build-arg COMMIT_HASH="${COMMIT_HASH}" \
  --build-arg BUILD_DATE="${BUILD_DATE}" \
  --tag "$IMAGE:${VERSION}" \
  --tag "$IMAGE:latest" \
  . -f Dockerfile

# push to image rep
docker push "$IMAGE:${VERSION}"
docker push "$IMAGE:latest"
