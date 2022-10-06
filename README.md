<div align="center">
[![GitHub Release][release-img]][release]
[![Test][test-img]][test]
[![GitHub All Releases][github-all-releases-img]][release]
</div>

# Curio

## Development

### Installation

to install modules locally

`go mod download`

### Testing

running classification tests

`go test ./pkg/classification/... -count=1`

running a single specific test

`go test -run ^TestSchema$ ./pkg/classification/schema -count=1`

[test]: https://github.com/Bearer/curio/actions/workflows/test.yml
[test-img]: https://github.com/Bearer/curio/actions/workflows/test.yml/badge.svg
[release]: https://github.com/Bearer/curio/releases
[release-img]: https://img.shields.io/github/release/Bearer/curio.svg?logo=github
[github-all-releases-img]: https://img.shields.io/github/downloads/Bearer/curio/total?logo=github
