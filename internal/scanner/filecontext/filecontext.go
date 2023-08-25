package filecontext

import (
	"context"
	"time"

	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/stats"
)

type Context struct {
	ctx         context.Context
	rules       map[string]*settings.Rule
	detectorSet detectorset.Set
	filename    string
	stats       *stats.FileStats
}

func New(
	ctx context.Context,
	rules map[string]*settings.Rule,
	detectorSet detectorset.Set,
	filename string,
	stats *stats.FileStats,
) *Context {
	return &Context{
		ctx:         ctx,
		rules:       rules,
		detectorSet: detectorSet,
		filename:    filename,
		stats:       stats,
	}
}

func (fileContext *Context) Err() error {
	return fileContext.ctx.Err()
}

func (fileContext *Context) Rules() map[string]*settings.Rule {
	return fileContext.rules
}

func (fileContext *Context) DetectAt(
	node *tree.Node,
	ruleID string,
	detectorContext detectortypes.Context,
) ([]*detectortypes.Detection, error) {
	return fileContext.detectorSet.DetectAt(node, ruleID, detectorContext)
}

func (fileContext *Context) RuleStats(ruleID string, startTime time.Time) {
	fileContext.stats.Rule(ruleID, startTime)
}

func (fileContext *Context) Filename() string {
	return fileContext.filename
}
