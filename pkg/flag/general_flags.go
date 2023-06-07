package flag

import (
	"fmt"

	"github.com/bearer/bearer/api"
	pointer "github.com/bearer/bearer/pkg/util/pointers"
	"github.com/rs/zerolog/log"
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
	DebugProfileFlag = Flag{
		Name:       "debug-profile",
		ConfigName: "debug-profile",
		Value:      false,
		Usage:      "Generate profiling data for debugging",
		Hide:       true,
	}
)

type GeneralFlagGroup struct {
	ConfigFile          *Flag
	APIKey              *Flag
	Host                *Flag
	DisableVersionCheck *Flag
	NoColor             *Flag
	DebugProfile        *Flag
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GeneralOptions struct {
	ConfigFile          string `json:"config_file" yaml:"config_file"`
	Client              *api.API
	DisableVersionCheck bool
	NoColor             bool `mapstructure:"no_color" json:"no_color" yaml:"no_color"`
	DebugProfile        bool
}

func NewGeneralFlagGroup() *GeneralFlagGroup {
	return &GeneralFlagGroup{
		ConfigFile:          &ConfigFileFlag,
		APIKey:              &APIKeyFlag,
		Host:                &HostFlag,
		DisableVersionCheck: &DisableVersionCheckFlag,
		NoColor:             &NoColorFlag,
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

	return GeneralOptions{
		Client:              client,
		ConfigFile:          getString(f.ConfigFile),
		DisableVersionCheck: getBool(f.DisableVersionCheck),
		NoColor:             getBool(f.NoColor),
		DebugProfile:        getBool(f.DebugProfile),
	}
}
