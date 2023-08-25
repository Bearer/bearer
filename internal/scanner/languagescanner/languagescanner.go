package languagescanner

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

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
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/bearer/bearer/internal/scanner/cache"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/rulescanner"
	"github.com/bearer/bearer/internal/scanner/stats"
)

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
	tree, err := scanner.parseAndAnalyze(ctx, fileInfo)
	if tree == nil || err != nil {
		return nil, err
	}

	if log.Trace().Enabled() {
		log.Trace().Msgf("tree:\n%s", tree.RootNode().Dump())
	}

	sharedCache := cache.NewShared(scanner.detectorSet.BuiltinAndSharedRuleIDs())
	fileContext := filecontext.New(
		ctx,
		scanner.rules,
		scanner.detectorSet,
		fileInfo.FileInfo.Name(),
		fileStats,
	)

	var detections []*detectortypes.Detection
	for _, ruleID := range scanner.detectorSet.TopLevelRuleIDs() {
		ruleDetections, err := rulescanner.Scan(
			fileContext,
			cache.NewCache(sharedCache),
			settings.NESTED_STRICT_SCOPE,
			ruleID,
			tree.RootNode(),
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

func (scanner *Scanner) parseAndAnalyze(ctx context.Context, fileInfo *file.FileInfo) (*tree.Tree, error) {
	if !slices.Contains(scanner.language.EnryLanguages(), fileInfo.Language) {
		return nil, nil
	}

	contentBytes, err := os.ReadFile(fileInfo.AbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	builder, err := ast.Parse(ctx, scanner.language, contentBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	return scanner.analyze(ctx, builder)
}

func (scanner *Scanner) analyze(ctx context.Context, builder *tree.Builder) (*tree.Tree, error) {
	if err := scanner.querySet.Query(ctx, builder, builder.SitterRootNode()); err != nil {
		return nil, fmt.Errorf("error running ast queries: %w", err)
	}

	analyzer := scanner.language.NewAnalyzer(builder)
	if err := analyzeNode(ctx, builder, analyzer, builder.SitterRootNode()); err != nil {
		return nil, fmt.Errorf("error running language analysis: %w", err)
	}

	return builder.Build(), nil
}

func analyzeNode(
	ctx context.Context,
	builder *tree.Builder,
	analyzer language.Analyzer,
	node *sitter.Node,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visitChildren := func() error {
		childCount := int(node.ChildCount())

		var disabledRules []string
		for i := 0; i < childCount; i++ {
			child := node.Child(i)
			if !child.IsNamed() {
				continue
			}

			disabledRules = addDisabledRules(builder, disabledRules, node)
			if err := analyzeNode(ctx, builder, analyzer, child); err != nil {
				return err
			}
		}

		return nil
	}

	return analyzer.Analyze(node, visitChildren)
}

func addDisabledRules(builder *tree.Builder, disabledRules []string, node *sitter.Node) []string {
	if node.Type() == "comment" {
		nextDisabledRules := disabledRules

		nodeContent := builder.ContentFor(node)
		if strings.Contains(nodeContent, "bearer:disable") {
			ruleIdsStr := strings.Split(nodeContent, "bearer:disable")[1]

			for _, ruleId := range strings.Split(ruleIdsStr, ",") {
				nextDisabledRules = append(nextDisabledRules, strings.TrimSpace(ruleId))
			}
		}

		return nextDisabledRules
	}

	builder.AddDisabledRules(node, disabledRules)

	return nil
}
