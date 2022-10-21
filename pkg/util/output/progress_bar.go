package output

import (
	"os"

	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(filesLength int, config settings.Config) *progressbar.ProgressBar {
	return progressbar.NewOptions(filesLength,
		progressbar.OptionSetVisibility(!config.Scan.Debug),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			PlainLogger(os.Stderr).Msgf("\n")
		}),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("files"),
		progressbar.OptionSetDescription("scanning..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}
