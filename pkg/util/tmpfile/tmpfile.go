package tmpfile

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

var ErrCreateFailed = errors.New("failed to create file")

func Create(tmpDir string, ext string) string {
	outputFile, err := os.CreateTemp(tmpDir, "*"+ext)
	if err != nil {
		log.Fatal().Msgf("got create fail error %s %s", err, ErrCreateFailed)
	}
	outputFile.Close()

	return outputFile.Name()
}
