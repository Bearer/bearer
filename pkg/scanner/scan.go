package scanner

import (
	"context"
	"log"

	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/bearer/curio/pkg/fanal/artifact"
	ftypes "github.com/bearer/curio/pkg/fanal/types"
	"github.com/bearer/curio/pkg/scanner/local"
	"github.com/bearer/curio/pkg/types"
)

// StandaloneSuperSet is used in the standalone mode
var StandaloneSuperSet = wire.NewSet(
	local.SuperSet,
	wire.Bind(new(Driver), new(local.Scanner)),
	NewScanner,
)

// StandaloneDockerSet binds docker dependencies
var StandaloneDockerSet = wire.NewSet(
	wire.Value([]image.Option(nil)), // optional functions
	image.NewContainerImage,
	aimage.NewArtifact,
	StandaloneSuperSet,
)

// StandaloneArchiveSet binds archive scan dependencies
var StandaloneArchiveSet = wire.NewSet(
	image.NewArchiveImage,
	aimage.NewArtifact,
	StandaloneSuperSet,
)

// StandaloneFilesystemSet binds filesystem dependencies
var StandaloneFilesystemSet = wire.NewSet(
	flocal.NewArtifact,
	StandaloneSuperSet,
)

// StandaloneRepositorySet binds repository dependencies
var StandaloneRepositorySet = wire.NewSet(
	remote.NewArtifact,
	StandaloneSuperSet,
)

// RemoteSuperSet is used in the client mode
var RemoteSuperSet = wire.NewSet(
	client.NewScanner,
	wire.Value([]client.Option(nil)),
	wire.Bind(new(Driver), new(client.Scanner)),
	NewScanner,
)

// RemoteFilesystemSet binds filesystem dependencies for client/server mode
var RemoteFilesystemSet = wire.NewSet(
	flocal.NewArtifact,
	RemoteSuperSet,
)

// RemoteArchiveSet binds remote archive dependencies
var RemoteArchiveSet = wire.NewSet(
	aimage.NewArtifact,
	image.NewArchiveImage,
	RemoteSuperSet,
)

// Scanner implements the Artifact and Driver operations
type Scanner struct {
	driver   Driver
	artifact artifact.Artifact
}

// Driver defines operations of scanner
type Driver interface {
	Scan(ctx context.Context, target, artifactKey string, blobKeys []string, options types.ScanOptions) (
		results types.Results, osFound *ftypes.OS, err error)
}

// NewScanner is the factory method of Scanner
func NewScanner(driver Driver, ar artifact.Artifact) Scanner {
	return Scanner{driver: driver, artifact: ar}
}

// ScanArtifact scans the artifacts and returns results
func (s Scanner) ScanArtifact(ctx context.Context, options types.ScanOptions) (types.Report, error) {
	artifactInfo, err := s.artifact.Inspect(ctx)
	if err != nil {
		return types.Report{}, xerrors.Errorf("failed analysis: %w", err)
	}
	defer func() {
		if err := s.artifact.Clean(artifactInfo); err != nil {
			log.Logger.Warnf("Failed to clean the artifact %q: %v", artifactInfo.Name, err)
		}
	}()

	results, osFound, err := s.driver.Scan(ctx, artifactInfo.Name, artifactInfo.ID, artifactInfo.BlobIDs, options)
	if err != nil {
		return types.Report{}, xerrors.Errorf("scan failed: %w", err)
	}

	if osFound != nil && osFound.Eosl {
		log.Logger.Warnf("This OS version is no longer supported by the distribution: %s %s", osFound.Family, osFound.Name)
		log.Logger.Warnf("The vulnerability detection may be insufficient because security updates are not provided")
	}

	// Layer makes sense only when scanning container images
	if artifactInfo.Type != ftypes.ArtifactContainerImage {
		removeLayer(results)
	}

	return types.Report{
		SchemaVersion: report.SchemaVersion,
		ArtifactName:  artifactInfo.Name,
		ArtifactType:  artifactInfo.Type,
		Metadata: types.Metadata{
			OS: osFound,

			// Container image
			ImageID:     artifactInfo.ImageMetadata.ID,
			DiffIDs:     artifactInfo.ImageMetadata.DiffIDs,
			RepoTags:    artifactInfo.ImageMetadata.RepoTags,
			RepoDigests: artifactInfo.ImageMetadata.RepoDigests,
			ImageConfig: artifactInfo.ImageMetadata.ConfigFile,
		},
		CycloneDX: artifactInfo.CycloneDX,
		Results:   results,
	}, nil
}

func removeLayer(results types.Results) {
	for i := range results {
		result := results[i]

		for j := range result.Packages {
			result.Packages[j].Layer = ftypes.Layer{}
		}
		for j := range result.Vulnerabilities {
			result.Vulnerabilities[j].Layer = ftypes.Layer{}
		}
		for j := range result.Misconfigurations {
			result.Misconfigurations[j].Layer = ftypes.Layer{}
		}
	}
}
