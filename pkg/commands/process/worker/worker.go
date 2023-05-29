package worker

import (
	"encoding/json"

	"net/http"
	"runtime"

	customdetector "github.com/bearer/bearer/new/scanner"
	"github.com/bearer/bearer/pkg/classification"
	config "github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/worker/blamer"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/detectors"
	"github.com/bearer/bearer/pkg/scanner"
	"golang.org/x/exp/slices"
)

type Worker struct {
	classifer *classification.Classifier
	scanners  []string
}

func (worker *Worker) Setup(config config.Config) error {
	worker.scanners = config.Scan.Scanner

	if slices.Contains(worker.scanners, "sast") {
		classifier, err := classification.NewClassifier(&classification.Config{Config: config})
		if err != nil {
			return err
		}

		err = detectors.SetupLegacyDetector(config.BuiltInRules)
		if err != nil {
			return err
		}

		err = customdetector.Setup(&config, classifier)
		if err != nil {
			return err
		}

		worker.classifer = classifier
	}

	return nil
}

func (worker *Worker) Scan(scanRequest work.ProcessRequest) error {
	blamer := blamer.New(scanRequest.Dir, scanRequest.BlameRevisionsFilePath, scanRequest.PreviousCommitSHA)

	return scanner.Scan(
		scanRequest.Dir,
		[]string{scanRequest.File.FilePath},
		blamer,
		scanRequest.ReportPath,
		worker.classifer,
		worker.scanners,
	)
}

func Start(port string) error {
	worker := Worker{}

	err := http.ListenAndServe(`localhost:`+port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close() //golint:all,errcheck

		switch r.URL.Path {
		case work.RouteInitialize:
			var config config.Config
			json.NewDecoder(r.Body).Decode(&config) //nolint:all,errcheck

			response := work.InitializeResponse{}

			err := worker.Setup(config)
			if err != nil {
				response.Error = err.Error()
			}

			json.NewEncoder(rw).Encode(response) //nolint:all,errcheck
		case work.RouteProcess:
			runtime.GC()
			var scanRequest work.ProcessRequest
			json.NewDecoder(r.Body).Decode(&scanRequest) //nolint:all,errcheck

			response := work.ProcessResponse{}

			err := worker.Scan(scanRequest)
			if err != nil {
				response.Error = err.Error()
			}

			json.NewEncoder(rw).Encode(response) //nolint:all,errcheck
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}))

	return err
}
