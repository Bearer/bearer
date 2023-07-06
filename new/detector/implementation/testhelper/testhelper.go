package testhelper

import (
	"context"
	"fmt"
	"testing"

	"github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/util/file"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type result struct {
	Position string
	Content  string
	Data     interface{}
}

func RunTest(
	t *testing.T,
	name string,
	compositionInstantiator func(map[string]*settings.Rule, *classification.Classifier) (types.Composition, error),
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

		composition, err := compositionInstantiator(make(map[string]*settings.Rule), classifier)
		if err != nil {
			tt.Fatalf("failed to create composition: %s", err)
		}
		defer composition.Close()

		fileInfo, err := file.FileInfoFromPath(fileName)
		if err != nil {
			tt.Fatalf("failed to create file info for %s: %s", fileName, err)
		}

		detections, err := composition.DetectFromFileWithTypes(context.Background(), nil, fileInfo, []string{detectorType}, nil)
		if err != nil {
			tt.Fatalf("failed to detect: %s", err)
		}

		results := make([]result, len(detections))
		for i, detection := range detections {
			node := detection.MatchNode
			results[i] = result{
				Position: fmt.Sprintf("%d:%d", node.StartLineNumber(), node.StartColumnNumber()),
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
