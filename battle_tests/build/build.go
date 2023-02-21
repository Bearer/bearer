package build

var (
	// These should be set via go build -ldflags -X 'xxxx'.
	Version       = "dev"
	BattleTestSHA = "devSHA"
	Attempt       = "1"
	Language      = "all"
	S3Bucket      = "bearer-battle-test-reports"
)
