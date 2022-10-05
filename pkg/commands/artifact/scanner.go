package artifact

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/bearer/curio/pkg/scanner"
)

// filesystemStandaloneScanner initializes a filesystem scanner in standalone mode
func filesystemStandaloneScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeFilesystemScanner(ctx, conf.Target, conf.Artifact)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a filesystem scanner: %w", err)
	}
	return s, cleanup, nil
}

// filesystemStandaloneScanner initializes a repository scanner in standalone mode
func repositoryStandaloneScanner(ctx context.Context, conf ScannerConfig) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeRepositoryScanner(ctx, conf.Target, conf.Artifact)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a filesystem scanner: %w", err)
	}
	return s, cleanup, nil
}
