package artifact

import (
	"context"
	"sort"

	"github.com/bearer/curio/pkg/fanal/types"
)

type Option struct {
	// AnalyzerGroup     analyzer.Group // It is empty in OSS
	// DisabledAnalyzers []analyzer.Type
	SkipFiles       []string
	SkipDirs        []string
	FilePatterns    []string
	NoProgress      bool
	Offline         bool
	InsecureSkipTLS bool
	AppDirs         []string
	RepoBranch      string
	RepoCommit      string
	RepoTag         string
}

func (o *Option) Sort() {
	// sort.Slice(o.DisabledAnalyzers, func(i, j int) bool {
	// 	return o.DisabledAnalyzers[i] < o.DisabledAnalyzers[j]
	// })
	sort.Strings(o.SkipFiles)
	sort.Strings(o.SkipDirs)
	sort.Strings(o.FilePatterns)
}

type Artifact interface {
	Inspect(ctx context.Context) (reference types.ArtifactReference, err error)
	Clean(reference types.ArtifactReference) error
}
