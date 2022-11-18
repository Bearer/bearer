package sync

import (
	"context"
	"time"

	"github.com/bearer/curio/battle_tests/build"
	"github.com/bearer/curio/battle_tests/config"
	repodb "github.com/bearer/curio/battle_tests/db"
	metricsscan "github.com/bearer/curio/battle_tests/metrics_scan"
	"github.com/bearer/curio/battle_tests/rediscli"
	"github.com/bearer/curio/battle_tests/sheet"
	"github.com/rs/zerolog/log"
)

func GetDocumentID(sheetClient *sheet.GoogleSheets) (documentID string, err error) {
	workerCount, err := rediscli.WorkerOnline()
	if err != nil {
		log.Err(err)
		return
	}

	if workerCount == 1 {
		log.Debug().Msgf("workerCount is 1... creating document %s in %s", build.CurioVersion, config.Runtime.Drive.ParentFolderId)
		doc := sheetClient.CreateDocument(build.CurioVersion, config.Runtime.Drive.ParentFolderId)

		log.Debug().Msgf("doc %s", doc.ID)
		err = rediscli.SetDocument(doc.ID)

		if err != nil {
			log.Error().Msgf("setting DocID in redis failed %s", err.Error())
			return
		}

		return doc.ID, err
	}

	for {
		documentID, err = rediscli.GetDocument()

		if err != nil || documentID == "" {
			log.Debug().Msgf("document couldn't be found... retrying")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return documentID, err
	}
}

func DoWork(ctx context.Context, items []repodb.Item, docID string, sheetClient *sheet.GoogleSheets) context.Context {
	selfContext, selfDone := context.WithCancel(ctx)
	go func() {
		for {
			repoCounter, err := rediscli.PickUpWork()
			if err != nil {
				log.Error().Msgf("repoCounter failed %s", err.Error())
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if repoCounter > len(items) {
				err := WorkerOffline(docID, sheetClient)

				if err != nil {
					log.Error().Msgf("error when setting worker offline %e", err)
				}

				selfDone()
				return
			}
			repository := items[repoCounter-1]

			metricsReport := make(chan *metricsscan.MetricsReport, 1)

			log.Debug().Msgf("picked up work for %s", repository.URL())
			metricsscan.ScanRepository(repository.URL(), metricsReport)
			// Uncomment this line if you want to fake the process
			// metricsscan.FakeScanRepository(repository.URL(), metricsReport)

			select {
			case <-ctx.Done():
				err := WorkerOffline(docID, sheetClient)

				if err != nil {
					log.Error().Msgf("error when setting worker offline %e", err)
				}

				return
			case metrics := <-metricsReport:
				sheetClient.InsertMetricsMustPass(docID, metrics)
			}
		}
	}()
	return selfContext
}

func WorkerOffline(docID string, sheetClient *sheet.GoogleSheets) error {
	log.Debug().Msgf("setting worker offline")
	workerCount, err := rediscli.WorkerOffline()
	if err != nil {
		return err
	}

	log.Debug().Msgf("worker count after offline %d...", workerCount)

	return nil
}
