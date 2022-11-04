---
layout: layouts/doc.njk
title: Quick Start
---

# Quick Start

Scan your first project in X minutes or less.

## Installation

Curio is available as a standalone executable binary. The latest release is available on the [releases tab](https://github.com/Bearer/curio/releases/latest), or use one of the methods below.

### Install Script

:warning: **Not working till public** :warning:

This script downloads the Curio binary automatically based on your OS and architecture.

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh
```

_Defaults to `./bin` as a bin directory and to the latest releases_

If you need to customize the options, use the following to pass parameters:

```bash
curl -sfL https://raw.githubusercontent.com/Bearer/curio/main/contrib/install.sh | sh -s -- -b /usr/local/bin
```

### Binary

Download the archive file for your operating system/architecture from [here](https://github.com/Bearer/curio/releases/latest/). Unpack the archive, and put the binary somewhere in your $PATH (on UNIX-y systems, /usr/local/bin or the like). Make sure it has permission to execute:

```bash
chmod +x ./curio
```

## Scan your project

Run `curio scan` on a project directory:

```bash
curio scan /path/to/your_project
```

or a single a file:

```bash
curio scan ./curio-ci-test/Pipfile.lock
```

<!-- TODO: insert sample output or video here -->

Additional options for using and configuring the `scan` command can be found in the [scan documentation](docs/reference/commands.md#scan).
