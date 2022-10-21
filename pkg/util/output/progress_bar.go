package output

import (
	"os"

	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(filesLength int) *progressbar.ProgressBar {
	return progressbar.NewOptions(filesLength,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			PlainLogger(os.Stderr).Msgf("\n")
		}),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("files"),
		progressbar.OptionSetDescription("scanning repository..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}
