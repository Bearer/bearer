---
title: Commands
---

Curio offers a number of commands to use and customize the CLI to your needs.

- [scan](#scan): Scans a repository, project, or file.

In addition to flags available for specific commands, the following flags are globally available for all commands:

- `--config [path]`, `-c [path]` (string): Instructs curio to use the path to a `yaml` config file.
- `--debug`, `-d`: Places curio in debug mode.
- `--generate-default-config`: Writes the default config to `curio-default.yaml`.
- `--quiet`, `-q`: Suppresses progress bar and log output.
- `--version`, `-v`: Displays the installed version number.

## scan

Scans a repository, project, or file.

```bash
curio scan [FLAGS] [PATH]
```

`PATH` (string) (optional) (default: current working directory): A path to a directory or file.

### `FLAGS`

#### Report Flags

- `-f`, `--format` format (json) (default "json")
- `--report` specify the kind of report (detectors) (default "detectors")

#### Scan Flags

- `--skip-path` specify the comma separated files and directories to skip (supports \* syntax), eg. --skip-path users/\*.go,users/admin.sql
- `--debug` enable debug logs
- `--disable-domain-resolution` skip attempt to resolve detected domains during classification (default false)
- `--domain-resolution-timeout` set timeout when attempting to resolve detected domains during classification
- `--internal-domains` define regular expressions for better classification of private or unreachable domains eg. --internal-domains="*.my-company.com,private.sh"

#### Worker Flags

- `--file-size-max` ignore files with file size larger than this config (default 25000000)
- `--files-to-batch` number of files to batch to worker (default 1)
- `--memory-max` if memory needed to scan a file surpasses this limit, skip the file (default 800000000)
- `--timeout` time allowed to complete scan (default 10m0s)
- `--timeout-file-max` maximum timeout assigned to scanning file, this config superseeds timeout-second-per-bytes (default 5m0s)
- `--timeout-file-min` minimum timeout assigned to scanning file, this config superseeds timeout-second-per-bytes (default 5s)
- `--timeout-file-second-per-bytes` number of file size bytes producing a second of timeout assigned to scanning a file (default 10000)
- `--timeout-worker-online` maximum time for worker process to come online (default 1m0s)
- `--workers` number of processing workers to spawn (default 1)

### Usage

Below are usage examples for the `scan` command:

```bash
# Scan a local project
curio scan /path/to/project

# Scan an individual file
curio scan /path/to/Pipfile.lock

# Scan the current directory with a custom config file
curio scan --config /path/to/custom-config.yaml
```
