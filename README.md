<div align="center">

  <a href="https://www.bearer.com" rel="nofollow">
    <img alt="Bearer logo" data-canonical-src="https://www.bearer.com/assets/bearer-logo.svg" src="https://www.bearer.com/assets/bearer-logo.svg" width="250">
  </a>

  <hr/>
  
  [![GitHub Release][release-img]][release]
  [![Test][test-img]][test]
  [![GitHub All Releases][github-all-releases-img]][release]
</div>

# Curio

## Example usage

Scan a project:

```bash
  $ curio scan /path/to/your_project
```

Scan a single file:

```bash

  $ curio scan ./curio-ci-test/Pipfile.lock
```

## Scan command flags:

### Report Flags:

- `-f`, `--format` format (table, json, jsonline) (default "jsonlines")

### Scan Flags:

- `--skip` specify the comma separated files and directories to skip (supports \* syntax), eg. --skip users/\*.go,users/admin.sql

### Worker Flags:

- `--file-size-max` ignore files with file size larger than this config (default 25000000)
- `--files-to-batch` number of files to batch to worker (default 1)
- `--memory-max` if memory needed to scan a file surpasses this limit, skip the file (default 800000000)
- `--timeout` time allowed to complete scan (default 10m0s)
- `--timeout-file-max` maximum timeout assigned to scanning file, this config superseeds timeout-second-per-bytes (default 5m0s)
- `--timeout-file-min` minimum timeout assigned to scanning file, this config superseeds timeout-second-per-bytes (default 5s)
- `--timeout-file-second-per-bytes` number of file size bytes producing a second of timeout assigned to scanning a file (default 10000)
- `--timeout-worker-online` maximum time for worker process to come online (default 1m0s)
- `--workers` number of processing workers to spawn (default 1)

## Development

### Installation

Install modules:

```bash
go mod download
```

### Testing

Running classification tests:

```bash
go test ./pkg/classification/... -count=1
```

Running a single specific test:

```bash
go test -run ^TestSchema$ ./pkg/classification/schema -count=1
```

---

[test]: https://github.com/Bearer/curio/actions/workflows/test.yml
[test-img]: https://github.com/Bearer/curio/actions/workflows/test.yml/badge.svg
[release]: https://github.com/Bearer/curio/releases
[release-img]: https://img.shields.io/github/release/Bearer/curio.svg?logo=github
[github-all-releases-img]: https://img.shields.io/github/downloads/Bearer/curio/total?logo=github
