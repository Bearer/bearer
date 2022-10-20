package build

var (
	// These should be set via go build -ldflags -X 'xxxx'.
	Version   = "dev"
	CommitSHA = "devSHA"
)
