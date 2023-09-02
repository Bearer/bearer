package testhelper

import (
	"context"
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/internal/classification"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/scanner/ast"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/filecontext"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/scanner/rulescanner"
	"github.com/bearer/bearer/internal/scanner/ruleset"
)

type result struct {
	Node    int
	Content string
	Data    interface{}
}

func RunTest(
	t *testing.T,
	name string,
	language language.Language,
	detectorType string,
	fileName string,
) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	t.Run(name, func(tt *testing.T) {
		querySet := query.NewSet(language.ID(), language.SitterLanguage())

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

		ruleSet, err := ruleset.New(language.ID(), make(map[string]*settings.Rule))
		if err != nil {
			tt.Fatalf("failed to create rule set: %s", err)
		}

		detectorSet, err := detectorset.New(
			querySet,
			language.NewBuiltInDetectors(classifier.Schema, querySet),
			ruleSet,
			language,
		)
		if err != nil {
			tt.Fatalf("failed to create detector set: %s", err)
		}

		if err := querySet.Compile(); err != nil {
			tt.Fatalf("failed to compile queries: %s", err)
		}

		fileContext := filecontext.New(
			context.Background(),
			detectorSet,
			fileName,
			nil,
		)

		contentBytes, err := os.ReadFile(fileName)
		if err != nil {
			tt.Fatalf("failed to read file: %s", err)
		}

		tree, err := ast.ParseAndAnalyze(context.Background(), language, querySet, contentBytes)
		if err != nil {
			tt.Fatalf("failed to parse file: %s", err)
		}

		rule, err := ruleSet.RuleByID(detectorType)
		if err != nil {
			tt.Fatalf("failed to lookup rule: %s", err)
		}

		detections, err := rulescanner.ScanTopLevelRule(fileContext, nil, tree, rule)
		if err != nil {
			tt.Fatalf("failed to scan with rule scanner: %s", err)
		}

		results := make([]result, len(detections))
		for i, detection := range detections {
			node := detection.MatchNode
			results[i] = result{
				Node:    node.ID,
				Content: node.Content(),
				Data:    detection.Data,
			}
		}

		yamlResults, err := yaml.Marshal(results)
		if err != nil {
			tt.Fatalf("failed to marshal results: %s", err)
		}

		cupaloy.SnapshotT(tt, tree.RootNode().Dump(), string(yamlResults))
	})
}
