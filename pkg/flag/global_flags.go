package flag

import (
	"time"

	"github.com/bearer/curio/pkg/util/cache"
	"github.com/spf13/cobra"
)

var (
	ConfigFileFlag = Flag{
		Name:       "config",
		ConfigName: "config",
		Shorthand:  "c",
		Value:      "bearer.yaml",
		Usage:      "config path",
		Persistent: true,
	}
	ShowVersionFlag = Flag{
		Name:       "version",
		ConfigName: "version",
		Shorthand:  "v",
		Value:      false,
		Usage:      "show version",
		Persistent: true,
	}
	QuietFlag = Flag{
		Name:       "quiet",
		ConfigName: "quiet",
		Shorthand:  "q",
		Value:      false,
		Usage:      "suppress progress bar and log output",
		Persistent: true,
	}
	DebugFlag = Flag{
		Name:       "debug",
		ConfigName: "debug",
		Shorthand:  "d",
		Value:      false,
		Usage:      "debug mode",
		Persistent: true,
	}
	InsecureFlag = Flag{
		Name:       "insecure",
		ConfigName: "insecure",
		Value:      false,
		Usage:      "allow insecure server connections when using TLS",
		Persistent: true,
	}
	TimeoutFlag = Flag{
		Name:       "timeout",
		ConfigName: "timeout",
		Value:      10 * time.Minute,
		Usage:      "Time allowed to complete scan",
		Persistent: true,
	}
	TimeoutFileMinimumFlag = Flag{
		Name:       "timeout-file-min",
		ConfigName: "timeout-file-min",
		Value:      5 * time.Second,
		Usage:      "what is the minimum timeout assigned to scanning file, this config superseeds timeout-second-per-bytes",
		Persistent: true,
	}
	TimeoutFileMaximumFlag = Flag{
		Name:       "timeout-file-max",
		ConfigName: "timeout-file-max",
		Value:      300 * time.Second, // 5 mins
		Usage:      "what is the maximum timeout assigned to scanning file, this config superseeds timeout-second-per-bytes",
		Persistent: true,
	}
	TimeoutFileSecondPerBytesFlag = Flag{
		Name:       "timeout-file-second-per-bytes",
		ConfigName: "timeout-file-second-per-bytes",
		Value:      10 * 1000, // 10kb/s
		Usage:      "how many file size bytes produces a second of timeout assigned to scanning file",
		Persistent: true,
	}
	FileSizeMaximumFlag = Flag{
		Name:       "file-size-max",
		ConfigName: "file-size-max",
		Value:      25 * 1000 * 1000, // 25 MB
		Usage:      "ignore files with file size larger than this config",
		Persistent: true,
	}
	MemoryMaximumFlag = Flag{
		Name:       "memory-max",
		ConfigName: "memory-max",
		Value:      800 * 1000 * 1000, // 800 MB
		Usage:      "if memory needed to scan a file suprasses this limit, skip the file",
		Persistent: true,
	}
	CacheDirFlag = Flag{
		Name:       "cache-dir",
		ConfigName: "cache.dir",
		Value:      cache.DefaultDir(),
		Usage:      "cache directory",
		Persistent: true,
	}
	GenerateDefaultConfigFlag = Flag{
		Name:       "generate-default-config",
		ConfigName: "generate-default-config",
		Value:      false,
		Usage:      "write the default config to curio-default.yaml",
		Persistent: true,
	}
)

// GlobalFlagGroup composes global flags
type GlobalFlagGroup struct {
	ConfigFile                *Flag
	ShowVersion               *Flag // spf13/cobra can't override the logic of version printing like VersionPrinter in urfave/cli. -v needs to be defined ourselves.
	Quiet                     *Flag
	Debug                     *Flag
	Insecure                  *Flag
	Timeout                   *Flag
	TimeoutFileMinimum        *Flag
	TimeoutFileMaximum        *Flag
	TimeoutFileSecondPerBytes *Flag
	FileSizeMaximum           *Flag
	MemoryMaximum             *Flag
	CacheDir                  *Flag
	GenerateDefaultConfig     *Flag
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GlobalOptions struct {
	ConfigFile                string
	ShowVersion               bool
	Quiet                     bool
	Debug                     bool
	Insecure                  bool
	Timeout                   time.Duration
	TimeoutFileMinimum        time.Duration
	TimeoutFileMaximum        time.Duration
	TimeoutFileSecondPerBytes int
	FileSizeMaximum           int
	MemoryMaximum             int
	CacheDir                  string
	GenerateDefaultConfig     bool
}

func NewGlobalFlagGroup() *GlobalFlagGroup {
	return &GlobalFlagGroup{
		ConfigFile:                &ConfigFileFlag,
		ShowVersion:               &ShowVersionFlag,
		Quiet:                     &QuietFlag,
		Debug:                     &DebugFlag,
		Insecure:                  &InsecureFlag,
		Timeout:                   &TimeoutFlag,
		TimeoutFileMinimum:        &TimeoutFileMinimumFlag,
		TimeoutFileMaximum:        &TimeoutFileMaximumFlag,
		TimeoutFileSecondPerBytes: &TimeoutFileSecondPerBytesFlag,
		FileSizeMaximum:           &FileSizeMaximumFlag,
		MemoryMaximum:             &MemoryMaximumFlag,
		CacheDir:                  &CacheDirFlag,
		GenerateDefaultConfig:     &GenerateDefaultConfigFlag,
	}
}

func (f *GlobalFlagGroup) flags() []*Flag {
	return []*Flag{f.ConfigFile, f.ShowVersion, f.Quiet, f.Debug, f.Insecure, f.TimeoutFileMinimum, f.TimeoutFileMaximum, f.TimeoutFileSecondPerBytes, f.FileSizeMaximum, f.MemoryMaximum, f.CacheDir, f.GenerateDefaultConfig}
}

func (f *GlobalFlagGroup) AddFlags(cmd *cobra.Command) {
	for _, flag := range f.flags() {
		addFlag(cmd, flag)
	}
}

func (f *GlobalFlagGroup) Bind(cmd *cobra.Command) error {
	for _, flag := range f.flags() {
		if err := bind(cmd, flag); err != nil {
			return err
		}
	}
	return nil
}

func (f *GlobalFlagGroup) ToOptions() GlobalOptions {
	return GlobalOptions{
		ConfigFile:                getString(f.ConfigFile),
		ShowVersion:               getBool(f.ShowVersion),
		Quiet:                     getBool(f.Quiet),
		Debug:                     getBool(f.Debug),
		Insecure:                  getBool(f.Insecure),
		Timeout:                   getDuration(f.Timeout),
		TimeoutFileMinimum:        getDuration(f.TimeoutFileMinimum),
		TimeoutFileMaximum:        getDuration(f.TimeoutFileMaximum),
		TimeoutFileSecondPerBytes: getInt(f.TimeoutFileSecondPerBytes),
		FileSizeMaximum:           getInt(f.FileSizeMaximum),
		MemoryMaximum:             getInt(f.MemoryMaximum),
		CacheDir:                  getString(f.CacheDir),
		GenerateDefaultConfig:     getBool(f.GenerateDefaultConfig),
	}
}
