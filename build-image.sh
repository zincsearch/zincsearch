#!/bin/sh
# author xiaojun207
# eg.1 : sh build-image.sh
# eg.2, set image: sh build-image.sh public.ecr.aws/zinclabs/zincsearch

VERSION=`git describe --tags --always` # eg.: 0.2.5
BUILD_DATE=`date +%Y%m%d` # eg.: 20220701
COMMIT_HASH=`git rev-parse HEAD` #
IMAGE="zinclabs/zincsearch"

if [ -n "$1" ]; then
  IMAGE="$1" #
fi

# build image
docker buildx build \
  --build-arg VERSION="${VERSION}" \
  --build-arg COMMIT_HASH="${COMMIT_HASH}" \
  --build-arg BUILD_DATE="${BUILD_DATE}" \
  --tag "$IMAGE:${VERSION}" \
  --tag "$IMAGE:latest" \
  . -f Dockerfile

# push to image rep
if [ -n "$1" ]; then
 docker push "$IMAGE:${VERSION}"
 docker push "$IMAGE:latest"
fi
