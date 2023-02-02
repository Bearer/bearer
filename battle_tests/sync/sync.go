package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/bearer/curio/battle_tests/build"
	"github.com/bearer/curio/battle_tests/config"
	repodb "github.com/bearer/curio/battle_tests/db"
	metricsscan "github.com/bearer/curio/battle_tests/metrics_scan"
	"github.com/bearer/curio/battle_tests/rediscli"
	"github.com/bearer/curio/battle_tests/sheet"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func DoWork(ctx context.Context, items []repodb.ItemWithLanguage, docID string, sheetClient *sheet.GoogleSheets) context.Context {
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

			log.Debug().Msgf("picked up work for %s", repository.FullName)
			metricsscan.ScanRepository(repository.HtmlUrl, repository.Language, metricsReport)
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

				err := uploadReportToS3(build.S3Bucket, docID, metrics)
				if err != nil {
					if awsErr, ok := err.(awserr.Error); ok {
						log.Error().Msgf("Failed to upload file to S3 code:%s message:%s", awsErr.Code(), awsErr.Message())
					} else {
						log.Error().Msgf("Failed to upload file to S3: %e", err)
					}
				}
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

func uploadReportToS3(bucketName string, documentID string, metrics *metricsscan.MetricsReport) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("eu-west-1"),
		},
	})

	if err != nil {
		log.Error().Msgf("Failed to create AWS session: %e", err)
		return err
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		log.Error().Msgf("Failed to serialize metrics to JSON: %e", err)
		return err
	}

	reader := bytes.NewReader(data)
	uploader := s3manager.NewUploader(sess)
	key_filename := strings.Replace(strings.Replace(metrics.URL, "https://", "", 1), "/", "_", -1) + ".json"
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(documentID + "/" + key_filename),
		Body:   reader,
	})

	if err != nil {
		return err
	}

	log.Debug().Msgf("S3 File Uploaded: %s", result.Location)

	return nil
}
