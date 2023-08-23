package stats

import (
	"regexp"
	"strings"

	"github.com/hhatto/gocloc"
	"github.com/jessevdk/go-flags"
)

type CmdOptions struct {
	Byfile         bool   `long:"by-file" description:"report results for every encountered source file"`
	SortTag        string `long:"sort" default:"code" description:"sort based on a certain column"`
	OutputType     string `long:"output-type" default:"default" description:"output type [values: default,cloc-xml,sloccount,json]"`
	ExcludeExt     string `long:"exclude-ext" description:"exclude file name extensions (separated commas)"`
	IncludeLang    string `long:"include-lang" description:"include language name (separated commas)"`
	Match          string `long:"match" description:"include file name (regex)"`
	NotMatch       string `long:"not-match" description:"exclude file name (regex)"`
	MatchDir       string `long:"match-d" description:"include dir name (regex)"`
	NotMatchDir    string `long:"not-match-d" description:"exclude dir name (regex)"`
	Debug          bool   `long:"debug" description:"dump debug log for developer"`
	SkipDuplicated bool   `long:"skip-duplicated" description:"skip duplicated files"`
	ShowLang       bool   `long:"show-lang" description:"print about all languages and extensions"`
}

func GoclocDetectorOutput(path string) (*gocloc.Result, error) {
	var opts CmdOptions
	clocOpts := gocloc.NewClocOptions()
	args := []string{
		"--output-type=json",
		path,
	}

	paths, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return nil, err
	}

	languages := gocloc.NewDefinedLanguages()

	// setup option for exclude extensions
	for _, ext := range strings.Split(opts.ExcludeExt, ",") {
		e, ok := gocloc.Exts[ext]
		if ok {
			clocOpts.ExcludeExts[e] = struct{}{}
		} else {
			clocOpts.ExcludeExts[ext] = struct{}{}
		}
	}

	// directory and file matching options
	if opts.Match != "" {
		clocOpts.ReMatch = regexp.MustCompile(opts.Match)
	}
	if opts.NotMatch != "" {
		clocOpts.ReNotMatch = regexp.MustCompile(opts.NotMatch)
	}
	if opts.MatchDir != "" {
		clocOpts.ReMatchDir = regexp.MustCompile(opts.MatchDir)
	}
	if opts.NotMatchDir != "" {
		clocOpts.ReNotMatchDir = regexp.MustCompile(opts.NotMatchDir)
	}

	// setup option for include languages
	for _, lang := range strings.Split(opts.IncludeLang, ",") {
		if _, ok := languages.Langs[lang]; ok {
			clocOpts.IncludeLangs[lang] = struct{}{}
		}
	}

	clocOpts.Debug = opts.Debug
	clocOpts.SkipDuplicated = opts.SkipDuplicated

	processor := gocloc.NewProcessor(languages, clocOpts)
	return processor.Analyze(paths)
}
