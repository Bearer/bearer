package flag

import (
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
	ConfigFile            *Flag
	ShowVersion           *Flag
	Quiet                 *Flag
	Debug                 *Flag
	GenerateDefaultConfig *Flag
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GlobalOptions struct {
	ConfigFile            string
	ShowVersion           bool
	Quiet                 bool
	Debug                 bool
	GenerateDefaultConfig bool
}

func NewGlobalFlagGroup() *GlobalFlagGroup {
	return &GlobalFlagGroup{
		ConfigFile:            &ConfigFileFlag,
		ShowVersion:           &ShowVersionFlag,
		Quiet:                 &QuietFlag,
		Debug:                 &DebugFlag,
		GenerateDefaultConfig: &GenerateDefaultConfigFlag,
	}
}

func (f *GlobalFlagGroup) flags() []*Flag {
	return []*Flag{
		f.ConfigFile,
		f.ShowVersion,
		f.Quiet,
		f.Debug,
		f.GenerateDefaultConfig,
	}
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
		ConfigFile:            getString(f.ConfigFile),
		ShowVersion:           getBool(f.ShowVersion),
		Quiet:                 getBool(f.Quiet),
		Debug:                 getBool(f.Debug),
		GenerateDefaultConfig: getBool(f.GenerateDefaultConfig),
	}
}
