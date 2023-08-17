package evaluator

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/new/detector/detection"
	cachepkg "github.com/bearer/bearer/new/detector/evaluator/cache"
	fileeval "github.com/bearer/bearer/new/detector/evaluator/file"
	"github.com/bearer/bearer/new/detector/evaluator/stats"
	detectorset "github.com/bearer/bearer/new/detector/set"
	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/new/language"
	"github.com/bearer/bearer/new/language/implementation"
	languagetypes "github.com/bearer/bearer/new/language/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/file"
)

type Evaluator struct {
	querySet      *query.Set
	languageName  string
	enryLanguages []string
	lang          languagetypes.Language
	rules         map[string]*settings.Rule
	detectorSet   detectortypes.DetectorSet
}

func New(
	langImplementation implementation.Implementation,
	schemaClassifier *schema.Classifier,
	rules map[string]*settings.Rule,
) (*Evaluator, error) {
	lang := language.New(langImplementation)
	querySet := lang.NewQuerySet()

	detectorSet, err := detectorset.New(
		querySet,
		langImplementation.NewBuiltInDetectors(schemaClassifier, querySet),
		rules,
		langImplementation.Name(),
		lang,
	)
	if err != nil {
		querySet.Close()
		return nil, fmt.Errorf("failed to create detector set: %s", err)
	}

	return &Evaluator{
		querySet:      querySet,
		languageName:  langImplementation.Name(),
		enryLanguages: langImplementation.EnryLanguages(),
		lang:          lang,
		rules:         rules,
		detectorSet:   detectorSet,
	}, nil
}

func (evaluator *Evaluator) LanguageName() string {
	return evaluator.languageName
}

func (evaluator *Evaluator) DetectFromFile(
	ctx context.Context,
	fileStats *stats.FileStats,
	file *file.FileInfo,
) ([]*detection.Detection, error) {
	if !slices.Contains(evaluator.enryLanguages, file.Language) {
		return nil, nil
	}

	contentBytes, err := os.ReadFile(file.AbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s", err)
	}

	tree, err := evaluator.lang.Parse(ctx, contentBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s", err)
	}

	fileEvaluator := fileeval.New(
		ctx,
		evaluator.detectorSet,
		tree,
		file.FileInfo.Name(),
		fileStats,
	)

	sharedCache := cachepkg.NewShared(evaluator.detectorSet.BuiltinAndSharedRuleIDs())

	var detections []*detection.Detection
	for _, ruleID := range evaluator.detectorSet.TopLevelRuleIDs() {
		cache := cachepkg.NewCache(sharedCache)
		rule := evaluator.rules[ruleID]

		sanitizerRuleID := ""
		if rule != nil {
			sanitizerRuleID = rule.SanitizerRuleID
		}

		ruleDetections, err := fileEvaluator.Evaluate(
			tree.RootNode(),
			ruleID,
			sanitizerRuleID,
			cache,
			settings.DefaultScope,
		)
		if err != nil {
			return nil, err
		}

		detections = append(detections, ruleDetections...)
	}

	return detections, nil
}

func (evaluator *Evaluator) Close() {
	evaluator.querySet.Close()
}
