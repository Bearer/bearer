package filecontext

import (
	"context"
	"time"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/stats"
)

type Context struct {
	ctx         context.Context
	detectorSet detectorset.Set
	filename    string
	stats       *stats.FileStats
}

func New(
	ctx context.Context,
	detectorSet detectorset.Set,
	filename string,
	stats *stats.FileStats,
) *Context {
	return &Context{
		ctx:         ctx,
		detectorSet: detectorSet,
		filename:    filename,
		stats:       stats,
	}
}

func (fileContext *Context) Err() error {
	return fileContext.ctx.Err()
}

func (fileContext *Context) DetectAt(
	node *tree.Node,
	detectorID int,
	detectorContext detectortypes.Context,
) (*detectorset.Result, error) {
	return fileContext.detectorSet.DetectAt(node, detectorID, detectorContext)
}

func (fileContext *Context) RuleStats(detectorID int, startTime time.Time) {
	if fileContext.stats != nil {
		fileContext.stats.Rule(fileContext.RuleIDFor(detectorID), startTime)
	}
}

func (fileContext *Context) RuleIDFor(detectorID int) string {
	return fileContext.detectorSet.RuleIDFor(detectorID)
}

func (fileContext *Context) DetectorIDFor(ruleID string) int {
	return fileContext.detectorSet.DetectorIDFor(ruleID)
}

func (fileContext *Context) Filename() string {
	return fileContext.filename
}
