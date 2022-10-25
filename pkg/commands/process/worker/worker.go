package worker

import (
	"encoding/json"

	"net/http"
	"runtime"

	classsification "github.com/bearer/curio/pkg/classification"
	config "github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/commands/process/worker/blamer"
	"github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/scanner"
)

func Start(port string, config config.Config) error {
	response := work.StatusResponse{}
	classifier, err := classsification.NewClassifier(&classsification.Config{Config: config})
	if err != nil {
		response.ClassifierError = err.Error()
	}

	err = detectors.SetupCustomDetector(config.CustomDetector)
	if err != nil {
		response.CustomDetectorError = err.Error()
	}

	err = http.ListenAndServe(`localhost`+port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close() //golint:all,errcheck

		switch r.URL.Path {
		case work.RouteStatus:
			json.NewEncoder(rw).Encode(response) //nolint:all,errcheck
		case work.RouteProcess:
			runtime.GC()

			var scanRequest work.ProcessRequest
			json.NewDecoder(r.Body).Decode(&scanRequest) //nolint:all,errcheck

			blamer := blamer.New(scanRequest.Dir, scanRequest.BlameRevisionsFilePath, scanRequest.PreviousCommitSHA)

			response := &work.ProcessResponse{}
			var filesList []string
			for _, file := range scanRequest.Files {
				filesList = append(filesList, file.FilePath)
			}

			response.Error = scanner.Scan(scanRequest.Dir, filesList, blamer, scanRequest.FilePath, classifier)

			json.NewEncoder(rw).Encode(response) //nolint:all,errcheck
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}))

	return err
}
