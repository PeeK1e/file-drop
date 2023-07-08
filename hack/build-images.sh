#!/bin/bash

set -m

printf "Building Images"
printf "%-*s\n" "25" "" | tr ' ' '-'

VERSION="${VERSION:-dev}"

docker build -t "docker.io/peek1e/file-api:$VERSION" -f src/docker/Dockerfile-api . &
docker build -t "docker.io/peek1e/file-web:$VERSION" -f src/docker/Dockerfile-web . &
docker build -t "docker.io/peek1e/file-cleaner:$VERSION" -f src/docker/Dockerfile-cleaner . &
docker build -t "docker.io/peek1e/file-migrations:$VERSION" -f src/docker/Dockerfile-migrations . &

while true; do fg 2> /dev/null; [ $? == 1 ] && break; done

printf "%-*s\n" "25" "" | tr ' ' '-'
printf "Done :)\n"
