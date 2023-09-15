package flag

import (
	"fmt"

	"github.com/bearer/bearer/api"
	pointer "github.com/bearer/bearer/internal/util/pointers"
	"github.com/rs/zerolog/log"
)

const (
	ErrorLogLevel = "error"
	InfoLogLevel  = "info"
	DebugLogLevel = "debug"
	TraceLogLevel = "trace"
)

var (
	HostFlag = Flag{
		Name:            "host",
		ConfigName:      "host",
		Value:           "my.bearer.sh",
		Usage:           "Specify the Host for sending the report.",
		DisableInConfig: true,
		Hide:            true,
	}
	APIKeyFlag = Flag{
		Name:            "api-key",
		ConfigName:      "api-key",
		Value:           "",
		Usage:           "Use your Bearer API Key to send the report to Bearer.",
		DisableInConfig: true,
		Hide:            true,
	}
	ConfigFileFlag = Flag{
		Name:            "config-file",
		ConfigName:      "config-file",
		Value:           "bearer.yml",
		Usage:           "Load configuration from the specified path.",
		DisableInConfig: true,
	}
	DisableVersionCheckFlag = Flag{
		Name:       "disable-version-check",
		ConfigName: "disable-version-check",
		Value:      false,
		Usage:      "Disable Bearer version checking",
	}
	NoColorFlag = Flag{
		Name:       "no-color",
		ConfigName: "report.no-color",
		Value:      false,
		Usage:      "Disable color in output",
	}
	IgnoreFileFlag = Flag{
		Name:            "ignore-file",
		ConfigName:      "ignore-file",
		Value:           "bearer.ignore",
		Usage:           "Load ignore file from the specified path.",
		DisableInConfig: true,
	}
	DebugFlag = Flag{
		Name:            "debug",
		ConfigName:      "debug",
		Value:           false,
		Usage:           "Enable debug logs. Equivalent to --log-level=debug",
		DisableInConfig: true,
	}
	LogLevelFlag = Flag{
		Name:       "log-level",
		ConfigName: "log-level",
		Value:      "info",
		Usage:      "Set log level (error, info, debug, trace)",
	}
	DebugProfileFlag = Flag{
		Name:            "debug-profile",
		ConfigName:      "debug-profile",
		Value:           false,
		Usage:           "Generate profiling data for debugging",
		Hide:            true,
		DisableInConfig: true,
	}
)

type GeneralFlagGroup struct {
	ConfigFile          *Flag
	IgnoreFile          *Flag
	APIKey              *Flag
	Host                *Flag
	DisableVersionCheck *Flag
	NoColor             *Flag
	DebugFlag           *Flag
	LogLevelFlag        *Flag
	DebugProfile        *Flag
}

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
}

func NewGeneralFlagGroup() *GeneralFlagGroup {
	return &GeneralFlagGroup{
		ConfigFile:          &ConfigFileFlag,
		APIKey:              &APIKeyFlag,
		Host:                &HostFlag,
		DisableVersionCheck: &DisableVersionCheckFlag,
		NoColor:             &NoColorFlag,
		IgnoreFile:          &IgnoreFileFlag,
		DebugFlag:           &DebugFlag,
		LogLevelFlag:        &LogLevelFlag,
		DebugProfile:        &DebugProfileFlag,
	}
}

func (f *GeneralFlagGroup) Name() string {
	return "General"
}

func (f *GeneralFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.ConfigFile,
		f.APIKey,
		f.Host,
		f.DisableVersionCheck,
		f.NoColor,
		f.IgnoreFile,
		f.DebugFlag,
		f.LogLevelFlag,
		f.DebugProfile,
	}
}

func (f *GeneralFlagGroup) ToOptions() GeneralOptions {
	var client *api.API
	apiKey := getString(f.APIKey)
	if apiKey != "" {
		client = api.New(api.API{
			Host:  getString(f.Host),
			Token: apiKey,
		})

		_, err := client.Hello()
		if err != nil {
			log.Debug().Msgf("couldn't initialize client -> %s", err.Error())
			client.Error = pointer.String(fmt.Sprintf("API key does not appear to be valid for %s.", client.Host))
		} else {
			log.Debug().Msgf("Initialized client for report")
		}
	}

	debug := getBool(f.DebugFlag)
	logLevel := getString(f.LogLevelFlag)
	if debug {
		logLevel = DebugLogLevel
	}

	return GeneralOptions{
		Client:              client,
		ConfigFile:          getString(f.ConfigFile),
		DisableVersionCheck: getBool(f.DisableVersionCheck),
		NoColor:             getBool(f.NoColor),
		IgnoreFile:          getString(f.IgnoreFile),
		Debug:               debug,
		LogLevel:            logLevel,
		DebugProfile:        getBool(f.DebugProfile),
	}
}
