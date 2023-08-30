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
	ruleID string,
	detectorContext detectortypes.Context,
) (*detectorset.Result, error) {
	return fileContext.detectorSet.DetectAt(node, ruleID, detectorContext)
}

func (fileContext *Context) RuleStats(ruleID string, startTime time.Time) {
	fileContext.stats.Rule(ruleID, startTime)
}

func (fileContext *Context) Filename() string {
	return fileContext.filename
}
