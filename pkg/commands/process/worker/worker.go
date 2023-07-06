package worker

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"os/signal"

	"net/http"
	"runtime"

	"github.com/bearer/bearer/new/detector/evaluator/stats"
	customdetector "github.com/bearer/bearer/new/scanner"
	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/debugprofile"
	config "github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/worker/work"
	"github.com/bearer/bearer/pkg/detectors"
	"github.com/bearer/bearer/pkg/scanner"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

var ErrorTimeoutReached = errors.New("file processing time exceeded")

type Worker struct {
	debug     bool
	classifer *classification.Classifier
	scanners  []string
}

func (worker *Worker) Setup(config config.Config) error {
	worker.debug = config.Scan.Debug
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

func (worker *Worker) Scan(ctx context.Context, scanRequest work.ProcessRequest) work.ProcessResponse {
	var fileStats *stats.FileStats
	if worker.debug {
		fileStats = stats.NewFileStats()
	}

	err := scanner.Scan(
		ctx,
		scanRequest.Dir,
		scanRequest.File.FilePath,
		scanRequest.ReportPath,
		worker.classifer,
		fileStats,
		worker.scanners,
	)

	if ctx.Err() != nil {
		err = ErrorTimeoutReached
	}

	var errorString string
	if err != nil {
		errorString = err.Error()
	}

	return work.ProcessResponse{
		FileStats: fileStats,
		Error:     errorString,
	}
}

func Start(port string) error {
	worker := Worker{}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	server := &http.Server{
		Addr: `localhost:` + port,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
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

				scanCtx, cancelScan := context.WithTimeout(ctx, scanRequest.File.Timeout)
				response := worker.Scan(scanCtx, scanRequest)

				cancelScan()

				json.NewEncoder(rw).Encode(response) //nolint:all,errcheck
			default:
				rw.WriteHeader(http.StatusNotFound)
			}
		}),
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()

		debugprofile.Stop()

		if err := server.Shutdown(context.Background()); err != nil {
			log.Debug().Msgf("error shutting down server: %s", err)
		}

		customdetector.Close()

		close(done)
	}()

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	<-done
	return nil
}
