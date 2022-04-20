#/bin/sh

aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/h9e2j3o7

# docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag public.ecr.aws/h9e2j3o7/zinc:v0.1.3-s3test .

docker buildx build --push --platform linux/amd64 --tag public.ecr.aws/h9e2j3o7/zinc:v0.1.3-s3test .


