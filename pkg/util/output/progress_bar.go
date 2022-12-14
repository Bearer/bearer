package output

import (
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(filesLength int, config settings.Config, display_type string) *progressbar.ProgressBar {
	hideProgress := config.Scan.Quiet || config.Scan.Debug
	return progressbar.NewOptions(filesLength,
		progressbar.OptionSetVisibility(!hideProgress),
		progressbar.OptionSetWriter(errorWriter),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			errorWriter.Write([]byte("\n")) //nolint:all,errcheck
		}),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString(display_type),
		progressbar.OptionSetDescription(" └"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}
