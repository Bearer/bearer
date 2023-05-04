#!/bin/bash

set -eu

echo
echo "Cloning $REPOSITORY_URL"
git clone --depth=1 --single-branch "$REPOSITORY_URL" /tmp/repository
cd /tmp/repository

echo
echo "Scanning"
bearer scan . "--host=$API_HOST" --api-key "$API_KEY"
