# Contributing

Thanks for your interest in contributing. This guide will help you get started.

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

<!-- TODO: Add details on building, running, etc from source. -->

### Testing

Running classification tests:

```bash
go test ./pkg/classification/... -count=1
```

Running a single specific test:

```bash
go test -run ^TestSchema$ ./pkg/classification/schema -count=1
```

### Submitting your changes as a pull request

<!-- add guidelines/format for contributions (branch naming, PR templates, etc) -->
