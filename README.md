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
