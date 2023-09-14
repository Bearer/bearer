package progressbar

import (
	"github.com/bearer/bearer/internal/commands/process/settings"
	"github.com/bearer/bearer/internal/util/output"
	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(filesLength int, config settings.Config, display_type string) *progressbar.ProgressBar {
	hideProgress := config.Scan.Quiet || config.Debug
	return progressbar.NewOptions(filesLength,
		progressbar.OptionSetVisibility(!hideProgress),
		progressbar.OptionSetWriter(output.ErrorWriter()),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			output.ErrorWriter().Write([]byte("\n")) //nolint:all,errcheck
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
