package worker

import (
	"encoding/json"

	"net/http"
	"runtime"

	customdetector "github.com/bearer/curio/new/scanner"
	"github.com/bearer/curio/pkg/classification"
	config "github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/commands/process/worker/blamer"
	"github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/scanner"
)

type Worker struct {
	classifer *classification.Classifier
}

func (worker *Worker) Setup(config config.Config) error {
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

	return nil
}

func (worker *Worker) Scan(scanRequest work.ProcessRequest) error {
	blamer := blamer.New(scanRequest.Dir, scanRequest.BlameRevisionsFilePath, scanRequest.PreviousCommitSHA)

	var filesList []string
	for _, file := range scanRequest.Files {
		filesList = append(filesList, file.FilePath)
	}

	return scanner.Scan(scanRequest.Dir, filesList, blamer, scanRequest.ReportPath, worker.classifer)
}

func Start(port string) error {
	worker := Worker{}

	err := http.ListenAndServe(`localhost:`+port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close() //golint:all,errcheck

		switch r.URL.Path {
		case work.RouteStatus:
			var config config.Config
			json.NewDecoder(r.Body).Decode(&config) //nolint:all,errcheck

			response := work.StatusResponse{}

			err := worker.Setup(config)
			if err != nil {
				response.ClassifierError = err.Error()
			}

			json.NewEncoder(rw).Encode(response) //nolint:all,errcheck
		case work.RouteProcess:
			runtime.GC()
			var scanRequest work.ProcessRequest
			json.NewDecoder(r.Body).Decode(&scanRequest) //nolint:all,errcheck

			response := work.ProcessResponse{}

			err := worker.Scan(scanRequest)
			if err != nil {
				response.Error = err
			}

			json.NewEncoder(rw).Encode(response) //nolint:all,errcheck
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}))

	return err
}
