package scanner

import (
	"context"

	"github.com/bearer/curio/pkg/types"
)

// Scanner implements the Artifact
type Scanner struct {
	artifact types.Artifact
}

// NewScanner is the factory method of Scanner
func NewScanner(ar types.Artifact) Scanner {
	return Scanner{artifact: ar}
}

// ScanArtifact scans the artifacts and returns results
func (s Scanner) ScanArtifact(ctx context.Context, options types.ScanOptions) (types.Report, error) {
	return types.Report{
		Artifact: s.artifact,
	}, nil
}
