# Curio

### Installation

to install modules locally

`go mod download`

### Testing

running classification tests

`go test ./pkg/classification/... -count=1`

running a single specific test

`go test -run ^TestSchema$ ./pkg/classification/schema -count=1`
