package build

var (
	// These should be set via go build -ldflags -X 'xxxx'.
	CurioVersion  = "dev"
	BattleTestSHA = "devSHA"
	Attempt       = "1"
	Language      = "all"
)
