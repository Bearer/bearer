---
title: Contributing to code
---

# Contributing code

If you're interested in contributing code to Curio, this guide will help you do it. If you haven't already, review the [contributing overview](/contributing/) for different ways you can contribute.

## Set up Curio locally

Curio is written in [Go](https://www.go.dev), so you'll need golang v1.19 or greater installed. Installation instructions for your architecture can be found at [golang downloads page](https://go.dev/dl/).

Next, confirm the installation by running the following command:

```bash
go version
```

You should receive a response with the installed version and architecture.

With Go installed, [Fork and clone](https://docs.github.com/en/get-started/quickstart/contributing-to-projects) the [Curio repository](https://github.com/Bearer/curio) locally using your preferred method. For example:

```bash
git clone https://github.com/Bearer/curio.git
```

Next, navigate into the curio project in your terminal, and install the required go modules:

```bash
go mod download
```

Finally, we use [direnv](https://direnv.net/) to manage env vars in development. This is primarily used for running tests. You can either:

- Set up your own by using `.envrc.example` as a base. You're unlikely to need to set anything other than `GITHUB_WORKSPACE` for integration tests, so feel free to skip this step or export it manually using: ```export GITHUB_WORKSPACE=`pwd` ```.
- Run tests using the script method mentioned [below](#running-tests).


## Running Curio locally from source

To run Curio from source use the following command from the curio directory:

```bash
go run ./cmd/curio/main.go [COMMAND]
```
Use commands and flags as normal in place of `[COMMAND]`.

## Running tests

Running tests is best performed using the [`run_tests.sh` script](https://github.com/Bearer/curio/blob/main/scripts/run_tests.sh). This will configure all needed variables to handle both unit and integration tests.

```
# From the project base
./scripts/run_tests.sh
```

Alternatively, you can run tests manually. Keep in mind you'll need to include the working path to the `GITHUB_WORKSPACE` environment variable if you choose this path.

Run all tests:

``` shell
go test ./...
```

Running classification tests:

```bash
go test ./pkg/classification/... -count=1
```

Running a single specific test:

```bash
go test -run ^TestSchema$ ./pkg/classification/schema -count=1
```

### Integration testing

These tests work by running an instance of the Curio application, with a
specified set of arguments, and then capturing the standard output and standard
error streams. The captured output is compared against a snapshot of the
expected output, using the [cupaloy](https://github.com/bradleyjkemp/cupaloy)
library.

The standard pattern is to create test case(s) using the
`testhelper.NewTestCase` factory function, and then pass these into the
`testhelper.RunTests` function for execution.

Integration tests are grouped within the top-level
[`integration`](/integration) subdirectory; see the existing tests for examples
to follow when adding a new test.

To update snapshots, set the `UPDATE_SNAPSHOTS` env variable to `true` when running a test. For example:

```bash
UPDATE_SNAPSHOTS=true go test ./...
```

**Note:** you will need to export `GITHUB_WORKSPACE` as the full path to the project so integration tests can run correctly. Alternatively, use the `run_tests.sh` script method mentioned in the [testing section](#running-tests) above.

## Generating CLI Documentation

Curio's reference pages are built, in-part, by files created from the CLI source. Whenever updating CLI commands or flags, it's best to rerun the generator below and include the created files in your PR. Upon a new doc build, the docs website will reflect these changes.

From the root of the project run:

```bash
go run ./scripts/gen-doc-yaml.go
```

This will auto generate yaml files for any updated CLI arguments.

## Getting help

If you need help getting set up, or are unsure how to navigate the codebase, join our [discord community]({{meta.links.discord}}) and we'll be happy to point you in the right direction.
