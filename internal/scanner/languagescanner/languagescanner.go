package languagescanner

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/filecontext"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/util/file"

	"github.com/bearer/bearer/internal/scanner/cache"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/rulescanner"
	"github.com/bearer/bearer/internal/scanner/stats"
)

type Scanner struct {
	language    language.Language
	ruleSet     *ruleset.Set
	querySet    *query.Set
	detectorSet detectorset.Set
}

func New(
	language language.Language,
	schemaClassifier *schema.Classifier,
	rules map[string]*settings.Rule,
) (*Scanner, error) {
	ruleSet, err := ruleset.New(language.ID(), rules)
	if err != nil {
		return nil, fmt.Errorf("error creating rule set: %w", err)
	}

	querySet := query.NewSet(language.ID(), language.SitterLanguage())

	detectorSet, err := detectorset.New(schemaClassifier, language, ruleSet, querySet)
	if err != nil {
		querySet.Close()
		return nil, fmt.Errorf("failed to create detector set: %w", err)
	}

	if err = querySet.Compile(); err != nil {
		querySet.Close()
		return nil, fmt.Errorf("error compiling query set: %w", err)
	}

	return &Scanner{
		language:    language,
		ruleSet:     ruleSet,
		querySet:    querySet,
		detectorSet: detectorSet,
	}, nil
}

func (scanner *Scanner) LanguageID() string {
	return scanner.language.ID()
}

func (scanner *Scanner) Scan(
	ctx context.Context,
	fileStats *stats.FileStats,
	fileInfo *file.FileInfo,
) ([]*detectortypes.Detection, error) {
	if !slices.Contains(scanner.language.EnryLanguages(), fileInfo.Language) {
		return nil, nil
	}

	contentBytes, err := os.ReadFile(fileInfo.AbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	tree, err := ast.ParseAndAnalyze(ctx, scanner.language, scanner.ruleSet, scanner.querySet, contentBytes)
	if err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("tree (%d nodes):\n%s", tree.NodeCount(), tree.RootNode().Dump())
	}

	fileContext := filecontext.New(
		ctx,
		scanner.ruleSet,
		scanner.detectorSet,
		fileInfo.FileInfo.Name(),
		fileStats,
	)

	return scanner.evaluateRules(fileContext, tree)
}

func (scanner *Scanner) evaluateRules(
	fileContext *filecontext.Context,
	tree *tree.Tree,
) ([]*detectortypes.Detection, error) {
	sharedCache := cache.NewShared(scanner.ruleSet.Rules())

	var detections []*detectortypes.Detection
	for _, rule := range scanner.ruleSet.Rules() {
		if rule.Type() != ruleset.RuleTypeTopLevel {
			continue
		}

		ruleDetections, err := rulescanner.ScanTopLevelRule(
			fileContext,
			cache.NewCache(tree, sharedCache),
			tree,
			rule,
		)
		if err != nil {
			return nil, err
		}

		detections = append(detections, ruleDetections...)
	}

	return detections, nil
}

func (scanner *Scanner) Close() {
	scanner.querySet.Close()
}
