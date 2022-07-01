#!/bin/sh
# author xiaojun207
# eg. : sh build-image.sh

VERSION="0.2.5" # 版本
COMMIT_HASH="79b9cbf66548b0c9728c1613fa2c64aeb19db81e" #
BUILD_DATE="20220630" #
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
