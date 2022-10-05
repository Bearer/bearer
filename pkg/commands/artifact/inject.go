package artifact

import (
	"context"

	"github.com/bearer/curio/pkg/scanner"
	"github.com/bearer/curio/pkg/types"
)

func initializeFilesystemScanner(ctx context.Context, path string, artifact types.Artifact) (scanner.Scanner, func(), error) {
	// wire.Build(scanner.StandaloneFilesystemSet)
	return scanner.Scanner{}, nil, nil
}

func initializeRepositoryScanner(ctx context.Context, url string, artifact types.Artifact) (
	scanner.Scanner, func(), error) {
	// wire.Build(scanner.StandaloneRepositorySet)
	return scanner.Scanner{}, nil, nil
}
