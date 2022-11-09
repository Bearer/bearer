package output

import (
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(filesLength int, config settings.Config) *progressbar.ProgressBar {
	return progressbar.NewOptions(filesLength,
		progressbar.OptionSetVisibility(!config.Scan.Debug),
		progressbar.OptionSetWriter(outputWriter),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			StdErrLogger().Msgf("\n")
		}),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("files"),
		progressbar.OptionSetDescription("Scanning target..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}
