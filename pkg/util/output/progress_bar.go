package output

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(filesLength int) *progressbar.ProgressBar {
	return progressbar.NewOptions(filesLength,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			log.Info().Msgf("scan completed")
		}),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("files"),
		progressbar.OptionSetDescription("Scanning repository..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}
