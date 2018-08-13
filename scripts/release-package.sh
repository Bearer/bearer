#!/usr/bin/env bash
set -e

echo "Fetching Bearer source code"
git pull origin master

echo "Starting publishing"
yarn lerna-publish