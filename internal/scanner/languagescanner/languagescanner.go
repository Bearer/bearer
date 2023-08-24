package languagescanner

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/scanner/ast"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/util/file"
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
	querySet := query.NewSet(language.SitterLanguage())

	detectorSet, err := detectorset.New(
		querySet,
		language.NewBuiltInDetectors(schemaClassifier, querySet),
		rules,
		language,
	)
	if err != nil {
		querySet.Close()
		return nil, fmt.Errorf("failed to create detector set: %s", err)
	}

	return &Scanner{
		querySet:    querySet,
		language:    language,
		rules:       rules,
		detectorSet: detectorSet,
	}, nil
}

func (scanner *Scanner) LanguageName() string {
	return scanner.language.Name()
}

func (scanner *Scanner) Scan(
	ctx context.Context,
	fileStats *stats.FileStats,
	file *file.FileInfo,
) ([]*detectortypes.Detection, error) {
	tree, err := scanner.parse(ctx, file)
	if tree == nil || err != nil {
		return nil, err
	}

	sharedCache := cache.NewShared(scanner.detectorSet.BuiltinAndSharedRuleIDs())
	rulesDisabledForNodes := mapNodesToDisabledRules(tree.RootNode())

	var detections []*detectortypes.Detection
	for _, ruleID := range scanner.detectorSet.TopLevelRuleIDs() {
		cache := cache.NewCache(sharedCache)
		rule := scanner.rules[ruleID]

		sanitizerRuleID := ""
		if rule != nil {
			sanitizerRuleID = rule.SanitizerRuleID
		}

		ruleDetections, err := rulescanner.Scan(
			ctx,
			scanner.detectorSet,
			file.FileInfo.Name(),
			fileStats,
			tree,
			rulesDisabledForNodes,
			tree.RootNode(),
			cache,
			settings.DefaultScope,
			ruleID,
			sanitizerRuleID,
		)
		if err != nil {
			return nil, err
		}

		detections = append(detections, ruleDetections...)
	}

	return detections, nil
}

func (scanner *Scanner) parse(ctx context.Context, file *file.FileInfo) (*tree.Tree, error) {
	if !slices.Contains(scanner.language.EnryLanguages(), file.Language) {
		return nil, nil
	}

	contentBytes, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s", err)
	}

	builder, err := ast.Parse(ctx, scanner.language, contentBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s", err)
	}

	analyzer := scanner.language.NewAnalyzer(builder)

	if err := scanner.querySet.Query(builder, builder.SitterRootNode()); err != nil {
		return nil, fmt.Errorf("error running ast queries: %w", err)
	}

	if err := analyzeNode(ctx, analyzer, builder.SitterRootNode()); err != nil {
		return nil, fmt.Errorf("error running language analysis: %w", err)
	}

	return builder.Build(), nil
}

func (scanner *Scanner) Close() {
	scanner.querySet.Close()
}

func mapNodesToDisabledRules(rootNode *tree.Node) map[string][]*tree.Node {
	res := make(map[string][]*tree.Node)
	var disabledRules []string
	err := rootNode.Walk(func(node *tree.Node, visitChildren func() error) error {
		if node.Type() == "comment" {
			// reset rules skipped array
			disabledRules = []string{}

			nodeContent := node.Content()
			if strings.Contains(nodeContent, "bearer:disable") {
				ruleIdsStr := strings.Split(nodeContent, "bearer:disable")[1]

				for _, ruleId := range strings.Split(ruleIdsStr, ",") {
					disabledRules = append(disabledRules, strings.TrimSpace(ruleId))
				}
			}

			return visitChildren()
		}

		// add rules skipped and node to result map
		for _, ruleId := range disabledRules {
			res[ruleId] = append(res[ruleId], node)
		}

		// reset rules skipped array
		disabledRules = []string{}
		return visitChildren()
	})

	// walk itself shouldn't trigger an error, and we aren't creating any
	if err != nil {
		panic(err)
	}

	return res
}

func analyzeNode(ctx context.Context, analyzer language.Analyzer, node *sitter.Node) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	visitChildren := func() error {
		childCount := int(node.ChildCount())

		for i := 0; i < childCount; i++ {
			if err := analyzeNode(ctx, analyzer, node.Child(i)); err != nil {
				return err
			}
		}

		return nil
	}

	return analyzer.Analyze(node, visitChildren)
}
