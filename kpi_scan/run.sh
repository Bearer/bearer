#!/bin/bash

set -eu

echo "Finding latest Bearer version"
curl -H "Accept: application/vnd.github.v3+json" \
     https://api.github.com/repos/bearer/bearer/releases/latest > /tmp/bearer-release.json

BEARER_URL=$(jq -r '.assets[] | select(.name | test("bearer_.*_linux_amd64.tar.gz")).browser_download_url' /tmp/bearer-release.json)
if [[ -z "$BEARER_URL" ]]; then
  echo "Can't find bearer URL"
  exit 1
fi

echo
echo "Downloading Bearer from $BEARER_URL"
curl --location --output /tmp/bearer.tar.gz "$BEARER_URL"
tar --extract --gunzip --file /tmp/bearer.tar.gz --directory /tmp/

echo
echo "Cloning $REPOSITORY_URL"
git clone "$REPOSITORY_URL" /tmp/repository
cd /tmp/repository

echo
echo "Scanning"
/tmp/bearer scan . "--host=$API_HOST" --api-key "$API_KEY"
