package languagescanner

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/scanner/ast"
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/pkg/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/pkg/scanner/detectors/types"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/scanner/variableshape"
	"github.com/bearer/bearer/pkg/util/file"

	"github.com/bearer/bearer/pkg/scanner/cache"
	"github.com/bearer/bearer/pkg/scanner/detectorset"
	"github.com/bearer/bearer/pkg/scanner/rulescanner"
	"github.com/bearer/bearer/pkg/scanner/stats"
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

	variableShapeSet, err := variableshape.NewSet(language, ruleSet)
	if err != nil {
		return nil, fmt.Errorf("error creating variable shape set: %w", err)
	}

	querySet := query.NewSet(language.ID(), language.SitterLanguage())

	detectorSet, err := detectorset.New(schemaClassifier, language, ruleSet, variableShapeSet, querySet)
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
) ([]*detectortypes.Detection, []*detectortypes.Detection, error) {
	if !slices.Contains(scanner.language.EnryLanguages(), fileInfo.Language) {
		return nil, nil, nil
	}

	contentBytes, err := os.ReadFile(fileInfo.AbsolutePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	tree, err := ast.ParseAndAnalyze(ctx, scanner.language, scanner.ruleSet, scanner.querySet, contentBytes)
	if err != nil {
		return nil, nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("tree (%d nodes):\n%s", tree.NodeCount(), tree.RootNode().Dump())
	}

	sharedCache := cache.NewShared(scanner.ruleSet.Rules())
	traversalCache := traversalstrategy.NewCache(tree.NodeCount())
	cache := cache.NewCache(tree, sharedCache)
	ruleScanner := rulescanner.New(
		ctx,
		scanner.detectorSet,
		fileInfo.Name(),
		fileStats,
		traversalCache,
		cache,
	)

	detections, err := scanner.evaluateRules(ruleScanner, cache, tree)
	expectedDetections, _ := scanner.ExpectedDetections(tree)

	return detections, expectedDetections, err
}

func (scanner *Scanner) ExpectedDetections(tree *tree.Tree) ([]*detectortypes.Detection, error) {
	var detections []*detectortypes.Detection
	nodes := tree.Nodes()
	for i := range tree.Nodes() {
		node := &nodes[i]
		if len(node.ExpectedRules()) > 0 {
			for _, expectedRule := range node.ExpectedRules() {
				rule, _ := scanner.ruleSet.RuleByID(expectedRule)
				detections = append(detections, []*detectortypes.Detection{
					{
						RuleID:    rule.ID(),
						MatchNode: node,
					},
				}...)
			}
		}
	}

	return detections, nil
}

func (scanner *Scanner) evaluateRules(
	ruleScanner *rulescanner.Scanner,
	cache *cache.Cache,
	tree *tree.Tree,
) (
	[]*detectortypes.Detection,
	error,
) {
	var detections []*detectortypes.Detection
	for _, rule := range scanner.ruleSet.Rules() {
		if rule.Type() != ruleset.RuleTypeTopLevel {
			continue
		}

		cache.Clear()
		ruleDetections, err := ruleScanner.Scan(tree.RootNode(), rule, traversalstrategy.NestedStrict)
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
