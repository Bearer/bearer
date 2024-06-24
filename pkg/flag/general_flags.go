package flag

import (
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/api"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
)

const (
	ErrorLogLevel = "error"
	InfoLogLevel  = "info"
	DebugLogLevel = "debug"
	TraceLogLevel = "trace"
)

type generalFlagGroup struct{ flagGroupBase }

var GeneralFlagGroup = &generalFlagGroup{flagGroupBase{name: "General"}}

var (
	HostFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:            "host",
		ConfigName:      "host",
		Value:           "api.cycode.com",
		Usage:           "Version check hostname.",
		DisableInConfig: true,
		Hide:            true,
	})

	APIKeyFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:            "api-key",
		ConfigName:      "api-key",
		Value:           "",
		Usage:           "Legacy.",
		DisableInConfig: true,
	})

	ConfigFileFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:            "config-file",
		ConfigName:      "config-file",
		Value:           "bearer.yml",
		Usage:           "Load configuration from the specified path.",
		DisableInConfig: true,
	})

	DisableVersionCheckFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:       "disable-version-check",
		ConfigName: "disable-version-check",
		Value:      false,
		Usage:      "Disable Bearer version checking",
	})

	NoColorFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:       "no-color",
		ConfigName: "report.no-color",
		Value:      false,
		Usage:      "Disable color in output",
	})

	IgnoreFileFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:            "ignore-file",
		ConfigName:      "ignore-file",
		Value:           "bearer.ignore",
		Usage:           "Load ignore file from the specified path.",
		DisableInConfig: true,
	})

	DebugFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:            "debug",
		ConfigName:      "debug",
		Value:           false,
		Usage:           "Enable debug logs. Equivalent to --log-level=debug",
		DisableInConfig: true,
	})

	LogLevelFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:       "log-level",
		ConfigName: "log-level",
		Value:      "info",
		Usage:      "Set log level (error, info, debug, trace)",
	})

	DebugProfileFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:            "debug-profile",
		ConfigName:      "debug-profile",
		Value:           false,
		Usage:           "Generate profiling data for debugging",
		Hide:            true,
		DisableInConfig: true,
	})

	IgnoreGitFlag = GeneralFlagGroup.add(flagtypes.Flag{
		Name:            "ignore-git",
		ConfigName:      "ignore-git",
		Value:           false,
		Usage:           "Ignore Git listing",
		Hide:            true,
		DisableInConfig: true,
	})
)

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GeneralOptions struct {
	ConfigFile          string `json:"config_file" yaml:"config_file"`
	Client              *api.API
	DisableVersionCheck bool
	NoColor             bool   `mapstructure:"no_color" json:"no_color" yaml:"no_color"`
	IgnoreFile          string `mapstructure:"ignore_file" json:"ignore_file" yaml:"ignore_file"`
	Debug               bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	LogLevel            string `mapstructure:"log-level" json:"log-level" yaml:"log-level"`
	DebugProfile        bool
	IgnoreGit           bool `mapstructure:"ignore-git" json:"ignore-git" yaml:"ignore-git"`
}

func (generalFlagGroup) SetOptions(options *flagtypes.Options, args []string) error {
	var client *api.API
	apiKey := getString(APIKeyFlag)
	if apiKey != "" {
		log.Debug().Msgf("API Key is no longer used please remove it from your config")
	}

	debug := getBool(DebugFlag)
	logLevel := getString(LogLevelFlag)
	if debug {
		logLevel = DebugLogLevel
	}

	options.GeneralOptions = flagtypes.GeneralOptions{
		Client:              client,
		ConfigFile:          getString(ConfigFileFlag),
		DisableVersionCheck: getBool(DisableVersionCheckFlag),
		NoColor:             getBool(NoColorFlag),
		IgnoreFile:          getString(IgnoreFileFlag),
		Debug:               debug,
		LogLevel:            logLevel,
		IgnoreGit:           getBool(IgnoreGitFlag),
		DebugProfile:        getBool(DebugProfileFlag),
	}

	return nil
}
