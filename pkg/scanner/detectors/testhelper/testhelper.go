package testhelper

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/scanner/ast"
	"github.com/bearer/bearer/pkg/scanner/ast/query"
	"github.com/bearer/bearer/pkg/scanner/ast/traversalstrategy"
	"github.com/bearer/bearer/pkg/scanner/detectorset"
	"github.com/bearer/bearer/pkg/scanner/language"
	"github.com/bearer/bearer/pkg/scanner/rulescanner"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/scanner/variableshape"
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
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
		FormatTimestamp: func(i interface{}) string {
			timestamp, _ := time.Parse(time.RFC3339, i.(string))
			return timestamp.Format("2006-01-02 15:04:05")
		},
	})

	t.Run(name, func(tt *testing.T) {
		classifier, err := classification.NewClassifier(&classification.Config{
			Config: settings.Config{
				Scan: flagtypes.ScanOptions{
					DisableDomainResolution: true,
					DomainResolutionTimeout: 0,
					Context:                 flagtypes.Context(flag.Empty),
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

		variableShapeSet, err := variableshape.NewSet(language, ruleSet)
		if err != nil {
			tt.Fatalf("failed to create variable shape set: %s", err)
		}

		querySet := query.NewSet(language.ID(), language.SitterLanguage())
		detectorSet, err := detectorset.New(
			classifier.Schema,
			language,
			ruleSet,
			variableShapeSet,
			querySet,
		)
		if err != nil {
			tt.Fatalf("failed to create detector set: %s", err)
		}

		if err := querySet.Compile(); err != nil {
			tt.Fatalf("failed to compile queries: %s", err)
		}

		contentBytes, err := os.ReadFile(fileName)
		if err != nil {
			tt.Fatalf("failed to read file: %s", err)
		}

		tree, err := ast.ParseAndAnalyze(context.Background(), language, ruleSet, querySet, contentBytes)
		if err != nil {
			tt.Fatalf("failed to parse file: %s", err)
		}

		ruleScanner := rulescanner.New(
			context.Background(),
			detectorSet,
			fileName,
			nil,
			traversalstrategy.NewCache(tree.NodeCount()),
			nil,
		)

		rule, err := ruleSet.RuleByID(detectorType)
		if err != nil {
			tt.Fatalf("failed to lookup rule: %s", err)
		}

		detections, err := ruleScanner.Scan(tree.RootNode(), rule, traversalstrategy.NestedStrict)
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
