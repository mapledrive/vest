#!/bin/sh

# build main image
docker build --force-rm -t vest-main . && \

# build build image
docker build --force-rm --no-cache -t vest-build -f ./internal/docker/build/Dockerfile .

# start container
docker run -d -p 443:443 -p 8080:8080 vest-build
