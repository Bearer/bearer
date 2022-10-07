package scanner

import (
	"context"
	"fmt"
	"os"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/report"
	"github.com/bearer/curio/pkg/types"
	"github.com/bearer/curio/pkg/util/blamer"
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

func Scan(rootDir string, FilesToScan []string, blamer blamer.Blamer, outputPath string) error {
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_TRUNC, 0666)

	if err != nil {
		return fmt.Errorf("fail opening ouput file %w", err)
	}
	defer file.Close()

	rep := report.JsonLinesReport{
		Blamer: blamer,
		File:   file,
	}

	if err := detectors.Extract(rootDir, FilesToScan, &rep); err != nil {
		return err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		rep.AddFillerLine()
	}

	return nil
}
