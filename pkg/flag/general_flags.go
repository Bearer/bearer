package flag

var (
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
}

// GlobalOptions defines flags and other configuration parameters for all the subcommands
type GeneralOptions struct {
	ConfigFile string `json:"config_file" yaml:"config_file"`
}

func NewGeneralFlagGroup() *GeneralFlagGroup {
	return &GeneralFlagGroup{
		ConfigFile: &ConfigFileFlag,
	}
}

func (f *GeneralFlagGroup) Name() string {
	return "General"
}

func (f *GeneralFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.ConfigFile,
	}
}

func (f *GeneralFlagGroup) ToOptions() GeneralOptions {
	return GeneralOptions{
		ConfigFile: getString(f.ConfigFile),
	}
}
