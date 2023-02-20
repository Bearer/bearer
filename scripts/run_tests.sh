#!/bin/sh -
#
# Run the Curio test suite

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
  rm -f ./curio || do_error "Failed to clean up"
}

trap do_cleanup 1 2 3 6

do_info "Building Curio binary..."
go build -a -o ./curio ./cmd/curio/main.go || do_error "Failed to build Curio binary"

[ -f curio ] || do_error "No Curio binary found"

TEST_ARGS=$DEFAULT_TEST_ARGS
[ $# -eq 0 ] || TEST_ARGS="$@"

do_info "Running tests..."
USE_BINARY=1 GITHUB_WORKSPACE=`pwd` go test $TEST_ARGS
TEST_STATUS=$?

do_cleanup

[ $TEST_STATUS -eq 0 ] || do_error "Tests failed"
