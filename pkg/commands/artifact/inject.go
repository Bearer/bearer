//go:build wireinject
// +build wireinject

package artifact

import (
	"context"

	"github.com/google/wire"

	"github.com/bearer/curio/pkg/fanal/artifact"
	"github.com/bearer/curio/pkg/fanal/cache"
	"github.com/bearer/curio/pkg/fanal/types"
	"github.com/bearer/curio/pkg/rpc/client"
	"github.com/bearer/curio/pkg/scanner"
)

// initializeFilesystemScanner is for filesystem scanning in standalone mode
func initializeFilesystemScanner(ctx context.Context, path string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.Cache, artifactOption artifact.Option) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneFilesystemSet)
	return scanner.Scanner{}, nil, nil
}

func initializeRepositoryScanner(ctx context.Context, url string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.Cache, artifactOption artifact.Option) (
	scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneRepositorySet)
	return scanner.Scanner{}, nil, nil
}

func initializeSBOMScanner(ctx context.Context, filePath string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.Cache, artifactOption artifact.Option) (
	scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneSBOMSet)
	return scanner.Scanner{}, nil, nil
}

// initializeRemoteFilesystemScanner is for filesystem scanning in client/server mode
func initializeRemoteFilesystemScanner(ctx context.Context, path string, artifactCache cache.ArtifactCache,
	remoteScanOptions client.ScannerOption, artifactOption artifact.Option) (scanner.Scanner, func(), error) {
	wire.Build(scanner.RemoteFilesystemSet)
	return scanner.Scanner{}, nil, nil
}
