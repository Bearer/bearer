#!/bin/bash

set -eu -o pipefail

echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin

git config --global --add safe.directory '*'

exec goreleaser "$@"
