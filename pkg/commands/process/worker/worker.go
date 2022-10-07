package main

import (
	"encoding/json"
	"fmt"

	"net/http"
	"os"
	"runtime"

	"github.com/rs/zerolog/log"

	"github.com/bearer/curio/pkg/commands/process/worker/blamer"
	"github.com/bearer/curio/pkg/commands/process/worker/work"
	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/scanner"
)

func main() {
	if len(os.Args) != 2 {
		log.Err(fmt.Errorf("not enough arguments, usage: ./worker :1234, program <port> ")).Send()
	}

	port := os.Args[1] // :1234

	// sets up logging level based on env variables
	// config.Load()

	log.Printf("scan worker listening on port %s", port)

	err := http.ListenAndServe(port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		switch r.URL.Path {
		case work.RouteStatus:
			var scanRequest work.ProcessRequest
			json.NewDecoder(r.Body).Decode(&scanRequest)

			response := &work.ProcessResponse{}

			if scanRequest.CustomDetectorConfig != nil {
				err := detectors.SetupCustomDetector(scanRequest.CustomDetectorConfig)
				if err != nil {
					response.Error = err
				}
			}

			json.NewEncoder(rw).Encode(response)
		case work.RouteProcess:
			runtime.GC()

			var scanRequest work.ProcessRequest
			json.NewDecoder(r.Body).Decode(&scanRequest)

			blamer := blamer.New(scanRequest.Dir, scanRequest.BlameRevisionsFilePath, scanRequest.PreviousCommitSHA)

			response := &work.ProcessResponse{}
			var filesList []string
			for _, file := range scanRequest.Files {
				filesList = append(filesList, file.FilePath)
			}

			response.Error = scanner.Scan(scanRequest.Dir, filesList, blamer, scanRequest.FilePath)

			json.NewEncoder(rw).Encode(response)
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	}))

	if err != nil {
		log.Printf("error serving %e", err)
	}
}
