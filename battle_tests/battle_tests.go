package main

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"

	"os"
	"os/signal"
	"syscall"

	"github.com/bearer/bearer/battle_tests/build"
	"github.com/bearer/bearer/battle_tests/config"
	"github.com/bearer/bearer/battle_tests/db"
	"github.com/bearer/bearer/battle_tests/rediscli"
	"github.com/bearer/bearer/battle_tests/sheet"
	"github.com/bearer/bearer/battle_tests/sync"
)

type Version struct {
	Version string
}

func main() {
	config.Load()
	rediscli.Setup()
	err := rediscli.Init()

	if err != nil {
		log.Debug().Msgf("failed to init redis")
	}

	log.Debug().Msgf("Binary version %s", build.Version)
	log.Debug().Msgf("battle tests SHA %s", build.BattleTestSHA)
	log.Debug().Msgf("S3 bucket %s", build.S3Bucket)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)

	db := db.UnmarshalRaw()
	log.Debug().Msgf("Selected %d repos for test", len(db))

	sheetClient := sheet.New()

	programCtx := context.TODO()
	docID, err := sync.GetDocumentID(sheetClient)
	if err != nil {
		log.Debug().Msgf("failed to get document id")
		log.Err(err).Send()
		return
	}

	log.Debug().Msgf("DocID is %s", docID)

	workerCtx := sync.DoWork(programCtx, db, docID, sheetClient)

	select {
	case <-shutdown:
		workerCtx.Done()
	case <-workerCtx.Done():
		return
	}
}
