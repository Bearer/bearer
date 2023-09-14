#!/bin/sh -
#
# Run the project test suite

# Ensures environment variables are not going to conflict with tests
unset ${!BEARER_*}

DEFAULT_TEST_ARGS="-count=1 -v ./..."

do_info() {
  printf "INFO: $*\n"
}

do_error() {
  printf "ERROR: $*\n" 1>&2
  exit 1
}

do_cleanup() {
  do_info "Cleaning up"
  rm -f ./bearer || do_error "Failed to clean up"
}

trap do_cleanup 1 2 3 6

do_info "Building binary..."
go build -a -o ./bearer ./cmd/bearer/main.go || do_error "Failed to build binary"

[ -f bearer ] || do_error "No binary found"

TEST_ARGS=$DEFAULT_TEST_ARGS
[ $# -eq 0 ] || TEST_ARGS="$@"

do_info "Running tests..."
USE_BINARY=1 GITHUB_WORKSPACE=`pwd` go test $TEST_ARGS
TEST_STATUS=$?

do_cleanup

[ $TEST_STATUS -eq 0 ] || do_error "Tests failed"
