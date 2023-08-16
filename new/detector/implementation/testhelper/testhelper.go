package testhelper

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"

	cachepkg "github.com/bearer/bearer/new/detector/evaluator/cache"
	fileeval "github.com/bearer/bearer/new/detector/evaluator/file"
	"github.com/bearer/bearer/new/detector/set"
	"github.com/bearer/bearer/new/language"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
)

type result struct {
	Position string
	Content  string
	Data     interface{}
}

func RunTest(
	t *testing.T,
	name string,
	langImplementation implementation.Implementation,
	detectorType string,
	fileName string,
) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	t.Run(name, func(tt *testing.T) {
		classifier, err := classification.NewClassifier(&classification.Config{
			Config: settings.Config{
				Scan: flag.ScanOptions{
					DisableDomainResolution: true,
					DomainResolutionTimeout: 0,
					Context:                 flag.Context(flag.Empty),
				},
			},
		})
		if err != nil {
			tt.Fatalf("failed to create classifier: %s", err)
		}

		querySet := query.NewSet(langImplementation.SitterLanguage())
		lang := language.New(langImplementation)

		detectorSet, err := set.New(
			querySet,
			langImplementation.NewBuiltInDetectors(classifier.Schema, querySet),
			make(map[string]*settings.Rule),
			langImplementation.Name(),
			lang,
		)
		if err != nil {
			tt.Fatalf("failed to create detector set: %s", err)
		}

		contentBytes, err := os.ReadFile(fileName)
		if err != nil {
			tt.Fatalf("failed to read file: %s", err)
		}

		tree, err := lang.Parse(context.Background(), contentBytes)
		if err != nil {
			tt.Fatalf("failed to parse content: %s", err)
		}

		fileEvaluator := fileeval.New(
			context.Background(),
			detectorSet,
			tree,
			fileName,
			nil,
		)
		if err != nil {
			tt.Fatalf("failed to create file evaluator: %s", err)
		}

		detections, err := fileEvaluator.Evaluate(
			tree.RootNode(),
			detectorType,
			"",
			cachepkg.NewCache(cachepkg.NewShared(detectorSet.BuiltinAndSharedRuleIDs())),
			settings.NESTED_SCOPE,
			true,
		)
		if err != nil {
			tt.Fatalf("failed to detect: %s", err)
		}

		results := make([]result, len(detections))
		for i, detection := range detections {
			node := detection.MatchNode
			results[i] = result{
				Position: fmt.Sprintf("%d:%d", node.ContentStart.Line, node.ContentStart.Column),
				Content:  node.Content(),
				Data:     detection.Data,
			}
		}

		yamlResults, err := yaml.Marshal(results)
		if err != nil {
			tt.Fatalf("failed to marshal results: %s", err)
		}

		cupaloy.SnapshotT(tt, string(yamlResults))
	})
}
