#/bin/sh

aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/zinclabs

# docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag public.ecr.aws/zinclabs/zincsearch:test .

docker buildx build --push --platform linux/amd64 --tag public.ecr.aws/zinclabs/zincsearch:test .

