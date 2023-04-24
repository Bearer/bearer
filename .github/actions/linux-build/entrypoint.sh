#!/bin/bash

set -eu -o pipefail

echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin

if [[ ! -z "${GITHUB_USER-}"]]; then
  echo "$GITHUB_TOKEN" | docker login ghcr.io --username "$GITHUB_USER" --password-stdin
fi

git config --global --add safe.directory '*'

exec goreleaser "$@"
