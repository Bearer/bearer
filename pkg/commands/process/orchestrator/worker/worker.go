package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"

	"net/http"
	"runtime"

	"github.com/bearer/bearer/new/detector/evaluator/stats"
	"github.com/bearer/bearer/new/scanner"
	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/debugprofile"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/work"
	config "github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/detectors"
	"github.com/bearer/bearer/pkg/report/writer"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

var ErrorTimeoutReached = errors.New("file processing time exceeded")

type Worker struct {
	debug           bool
	classifer       *classification.Classifier
	enabledScanners []string
	sastScanner     *scanner.Scanner
}

func (worker *Worker) Setup(config config.Config) error {
	worker.debug = config.Debug
	worker.enabledScanners = config.Scan.Scanner

	if slices.Contains(worker.enabledScanners, "sast") {
		classifier, err := classification.NewClassifier(&classification.Config{Config: config})
		if err != nil {
			return err
		}

		err = detectors.SetupLegacyDetector(config.BuiltInRules)
		if err != nil {
			return err
		}

		sastScanner, err := scanner.New(classifier.Schema, config.Rules)
		if err != nil {
			return err
		}

		worker.sastScanner = sastScanner
		worker.classifer = classifier
	}

	return nil
}

func (worker *Worker) Scan(ctx context.Context, scanRequest work.ProcessRequest) (*stats.FileStats, error) {
	var fileStats *stats.FileStats
	if worker.debug {
		fileStats = stats.NewFileStats()
	}

	file, err := os.Create(scanRequest.ReportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open output file %w", err)
	}
	defer file.Close()

	err = detectors.Extract(
		ctx,
		scanRequest.Dir,
		scanRequest.File.FilePath,
		&writer.Detectors{
			Classifier: worker.classifer,
			File:       file,
		},
		fileStats,
		worker.enabledScanners,
		worker.sastScanner,
	)

	if ctx.Err() != nil {
		return fileStats, ErrorTimeoutReached
	}

	return fileStats, err
}

func (worker *Worker) Close() {
	if worker.sastScanner != nil {
		worker.sastScanner.Close()
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
				fileStats, err := worker.Scan(scanCtx, scanRequest)
				var errorString string
				if err != nil {
					errorString = err.Error()
				}

				cancelScan()

				json.NewEncoder(rw).Encode(work.ProcessResponse{ //nolint:all,errcheck
					FileStats: fileStats,
					Error:     errorString,
				})
			case work.RouteReduceMemory:
				log.Trace().Msgf("attempting to reduce memory usage")
				runtime.GC()
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

		worker.Close()

		close(done)
	}()

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	<-done
	return nil
}
