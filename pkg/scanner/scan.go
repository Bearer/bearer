package scanner

import (
	"fmt"
	"os"

	classsification "github.com/bearer/curio/pkg/classification"
	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/report/writer"
	"github.com/bearer/curio/pkg/util/blamer"
)

func Scan(rootDir string, FilesToScan []string, blamer blamer.Blamer, outputPath string) error {
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_TRUNC, 0666)

	if err != nil {
		return fmt.Errorf("fail opening ouput file %w", err)
	}
	defer file.Close()

	rep := writer.JSONLines{
		Blamer:     blamer,
		Classifier: classsification.NewClassifier(&classsification.Config{}),
		File:       file,
	}

	if err := detectors.Extract(rootDir, FilesToScan, &rep); err != nil {
		return err
	}

	return nil
}
