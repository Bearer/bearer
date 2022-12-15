# Contributing

Thanks for your interest in contributing. This guide will help you get started.

- [Basics](#basics)
- [Development](#development)
  - [Installing Curio](#installing-curio)
  - [Running Curio locally from source](#running-curio-locally-from-source)
- [Testing](#testing)
  - [Integration testing](#integration-testing)
- [Generating CLI Documentation](#generating-cli-documentation)
- [Submitting your changes as a pull request](#submitting-your-changes-as-a-pull-request)

## Basics

If you've discovered a bug or thing you have and aren't sure how to fix it, [open an issue](https://github.com/Bearer/curio/issues/new/choose) with details on how to recreate the problem.

## Development

Curio is written in [Go](https://www.go.dev), so you'll need golang v1.19 or greater installed. Installation instructions for your architecture can be found at [golang downloads page](https://go.dev/dl/).

Next, confirm the installation by running the following command:

```bash
go version
```

You should receive a response with the installed version and architecture.

### Installing Curio

[Fork and clone](https://docs.github.com/en/get-started/quickstart/contributing-to-projects) the repository locally using your preferred method. For example:

```bash
git clone https://github.com/Bearer/curio.git
```

Next, navigate into the curio project in your terminal, and install the required go modules:

```bash
go mod download
```

### Running Curio locally from source

To run Curio from source use the following command from the curio directory:

```bash
go run ./cmd/curio/main.go [COMMAND]
```
Use commands and flags as normal in place of `[COMMAND]`.


## Testing

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

## Generating CLI Documentation

Curio's reference pages are built, in-part, by files created from the CLI source. Whenever updating CLI commands or flags, it's best to rerun the generator below and include the created files in your PR. Upon a new doc build, the docs website will reflect these changes.

From the root of the project run:

```bash
go run ./scripts/gen-doc-yaml.go
```

This will auto generate yaml files for any updated CLI arguments.

## Submitting your changes as a pull request

To Contribute new code, fixes, or doc changes to Curio it is best to by searching existing [issues](https://github.com/Bearer/curio/issues), or by joining the [community discord](https://discord.gg/eaHZBJUXRF). Here are some tips to get you started:

- If you're unfamiliar with the process, check out GitHub's guide: [Creating a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request)
- Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) whenever creating a pull request.
- As a best practice, we ask as much as possible that you sign your commits. Especially, we require our core contributors to sign their commits. Check out GitHub's page on [signing commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits) for a complete guide.

When in doubt, don't hesitate to raise an issue or jump into the discord.
