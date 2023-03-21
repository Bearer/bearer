package flag

var (
	APIKeyFlag = Flag{
		Name:            "api-key",
		ConfigName:      "api-key",
		Value:           "",
		Usage:           "Use your Bearer API Key to send the report to Bearer.",
		DisableInConfig: true,
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
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GeneralOptions struct {
	ConfigFile string `json:"config_file" yaml:"config_file"`
	APIKey     string `json:"api_key" yaml:"api_key"`
}

func NewGeneralFlagGroup() *GeneralFlagGroup {
	return &GeneralFlagGroup{
		ConfigFile: &ConfigFileFlag,
		APIKey:     &APIKeyFlag,
	}
}

func (f *GeneralFlagGroup) Name() string {
	return "General"
}

func (f *GeneralFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.ConfigFile,
		f.APIKey,
	}
}

func (f *GeneralFlagGroup) ToOptions() GeneralOptions {
	return GeneralOptions{
		ConfigFile: getString(f.ConfigFile),
		APIKey:     getString(f.APIKey),
	}
}
