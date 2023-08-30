package languagescanner

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/filecontext"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/util/file"
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/scanner/cache"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/rulescanner"
	"github.com/bearer/bearer/internal/scanner/stats"
)

const minimumNodeCountForCache = 20_000

type Scanner struct {
	querySet    *query.Set
	language    language.Language
	rules       map[string]*settings.Rule
	detectorSet detectorset.Set
}

func New(
	language language.Language,
	schemaClassifier *schema.Classifier,
	rules map[string]*settings.Rule,
) (*Scanner, error) {
	querySet := query.NewSet(language.ID(), language.SitterLanguage())

	detectorSet, err := detectorset.New(
		querySet,
		language.NewBuiltInDetectors(schemaClassifier, querySet),
		rules,
		language,
	)
	if err != nil {
		querySet.Close()
		return nil, fmt.Errorf("failed to create detector set: %w", err)
	}

	return &Scanner{
		querySet:    querySet,
		language:    language,
		rules:       rules,
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

	tree, err := ast.ParseAndAnalyze(ctx, scanner.language, scanner.querySet, contentBytes)
	if err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("tree (%d nodes):\n%s", tree.NodeCount(), tree.RootNode().Dump())
	}

	fileContext := filecontext.New(
		ctx,
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
	var sharedCache *cache.Shared
	if tree.NodeCount() > minimumNodeCountForCache {
		sharedCache = cache.NewShared(scanner.detectorSet.BuiltinAndSharedRuleIDs())
		log.Trace().Msg("cache enabled")
	}

	var detections []*detectortypes.Detection
	for _, ruleID := range scanner.detectorSet.TopLevelRuleIDs() {
		var ruleCache *cache.Cache
		if sharedCache != nil {
			ruleCache = cache.NewCache(sharedCache)
		}

		ruleDetections, err := rulescanner.ScanTopLevelRule(
			fileContext,
			ruleCache,
			tree,
			ruleID,
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
