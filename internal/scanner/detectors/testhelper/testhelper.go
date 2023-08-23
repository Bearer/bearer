package testhelper

import (
	"context"
	"fmt"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"

	"github.com/bearer/bearer/internal/classification"
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/scanner/language"
	"github.com/bearer/bearer/internal/scanner/languagescanner"
	"github.com/bearer/bearer/internal/util/file"
)

type result struct {
	Position string
	Content  string
	Data     interface{}
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

		languageScanner, err := languagescanner.New(language, classifier.Schema, make(map[string]*settings.Rule))
		if err != nil {
			tt.Fatalf("failed to create language scanner: %s", err)
		}

		fileInfo, err := file.FileInfoFromPath(fileName)
		if err != nil {
			tt.Fatalf("failed to create file info: %s", err)
		}

		detections, err := languageScanner.Scan(context.Background(), nil, fileInfo)
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
