#!/bin/bash

VERSION="${VERSION:-dev}"

docker build -t "docker.io/peek1e/file-api:$VERSION" -f src/docker/Dockerfile-api . &
docker build -t "docker.io/peek1e/file-web:$VERSION" -f src/docker/Dockerfile-web . &
docker build -t "docker.io/peek1e/file-cleaner:$VERSION" -f src/docker/Dockerfile-cleaner . &
docker build -t "docker.io/peek1e/file-migration:$VERSION" -f src/docker/Dockerfile-migration . &
