package testhelper

import (
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bearer/curio/new/detector/types"
	"github.com/bearer/curio/pkg/classification"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/flag"
	"github.com/bearer/curio/pkg/util/file"
	"github.com/bradleyjkemp/cupaloy"
	"gopkg.in/yaml.v2"
)

var TrimPath = "testdata/testcases/"

func RunTest(t *testing.T,
	rules map[string][]byte,
	testCasesPath string,
	compositionInstantiator func(map[string]settings.Rule, *classification.Classifier) (types.Composition, error)) {
	var rulesConfig map[string]settings.Rule = make(map[string]settings.Rule)

	for ruleName, ruleContent := range rules {
		var parsedRule settings.Rule
		err := yaml.Unmarshal(ruleContent, &parsedRule)
		if err != nil {
			t.Fatal(err)
		}
		rulesConfig[ruleName] = parsedRule
	}

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
		t.Fatal(err)
	}

	composition, err := compositionInstantiator(rulesConfig, classifier)
	if err != nil {
		t.Fatal(err)
	}

	files := []string{}
	filepath.WalkDir(testCasesPath, func(path string, d fs.DirEntry, walkDirErr error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, strings.TrimPrefix(path, TrimPath))
		return nil
	})

	fileInfos := []*file.FileInfo{}

	file.IterateFilesList(
		testCasesPath,
		files,
		func(dir *file.Path) (bool, error) {
			return true, nil
		},
		func(file *file.FileInfo) error {
			fileInfos = append(fileInfos, file)
			return nil
		},
	)

	for _, testCase := range fileInfos {
		t.Run(testCase.RelativePath, func(t *testing.T) {
			customDetections, err := composition.DetectFromFile(testCase)
			if err != nil {
				t.Fatal(err)
			}

			bytes, err := json.MarshalIndent(customDetections, "", "\t")
			if err != nil {
				t.Fatal(err)
			}

			cupaloy.SnapshotT(t, string(bytes))
		})

	}
}
