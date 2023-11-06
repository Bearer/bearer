package stats

import (
	"time"

	"github.com/bearer/bearer/internal/flag"
	"github.com/bearer/bearer/internal/util/output"

	"github.com/hhatto/gocloc"
	"github.com/schollz/progressbar/v3"
)

func GoclocDetectorOutput(path string, opts flag.Options) (*gocloc.Result, error) {
	clocOpts := gocloc.NewClocOptions()
	clocOpts.SkipDuplicated = true
	output.StdErrLog("Analyzing codebase")

	if !hideProgress(opts) {
		progressBar := getProgressBar()
		defer progressBar.Close()
		clocOpts.OnCode = func(line string) {
			progressBar.Add(1)
		}
	}

	languages := gocloc.NewDefinedLanguages()
	processor := gocloc.NewProcessor(languages, clocOpts)

	return processor.Analyze([]string{path})
}

func hideProgress(opts flag.Options) bool {
	return opts.ScanOptions.HideProgressBar || opts.ScanOptions.Quiet || opts.Debug
}

func getProgressBar() *progressbar.ProgressBar {

	return progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(output.ErrorWriter()),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			output.ErrorWriter().Write([]byte("\n")) //nolint:all,errcheck
		}),
	)
}
