package flag

import (
	"github.com/bearer/bearer/api"
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
		Value:           "",
		Usage:           "Load configuration from the specified path.",
		DisableInConfig: true,
	}
)

type GeneralFlagGroup struct {
	ConfigFile *Flag
	APIKey     *Flag
	Host       *Flag
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GeneralOptions struct {
	ConfigFile string `json:"config_file" yaml:"config_file"`
	Client     *api.API
}

func NewGeneralFlagGroup() *GeneralFlagGroup {
	return &GeneralFlagGroup{
		ConfigFile: &ConfigFileFlag,
		APIKey:     &APIKeyFlag,
		Host:       &HostFlag,
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
			log.Error().Msgf("couldn't initialize client -> %s", err.Error())
		} else {
			log.Debug().Msgf("Initialized client for report")
			return GeneralOptions{
				ConfigFile: getString(f.ConfigFile),
				Client:     client,
			}
		}
	}

	return GeneralOptions{
		ConfigFile: getString(f.ConfigFile),
	}
}
