package filecontext

import (
	"context"
	"time"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/scanner/stats"
)

type Context struct {
	ctx         context.Context
	ruleSet     *ruleset.Set
	detectorSet detectorset.Set
	filename    string
	stats       *stats.FileStats
}

func New(
	ctx context.Context,
	ruleSet *ruleset.Set,
	detectorSet detectorset.Set,
	filename string,
	stats *stats.FileStats,
) *Context {
	return &Context{
		ctx:         ctx,
		ruleSet:     ruleSet,
		detectorSet: detectorSet,
		filename:    filename,
		stats:       stats,
	}
}

func (fileContext *Context) Filename() string {
	return fileContext.filename
}

func (fileContext *Context) Rule(ruleIndex int) *ruleset.Rule {
	return fileContext.ruleSet.Rules()[ruleIndex]
}

func (fileContext *Context) RuleStats(rule *ruleset.Rule, startTime time.Time) {
	if fileContext.stats != nil {
		fileContext.stats.Rule(rule.ID(), startTime)
	}
}

func (fileContext *Context) DetectAt(
	node *tree.Node,
	rule *ruleset.Rule,
	detectorContext detectortypes.Context,
) (*detectorset.Result, error) {
	if fileContext.ctx.Err() != nil {
		return nil, fileContext.ctx.Err()
	}

	return fileContext.detectorSet.DetectAt(node, rule, detectorContext)
}
