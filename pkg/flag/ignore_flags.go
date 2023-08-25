package flag

var (
	BearerIgnoreFileFlag = Flag{
		Name:            "bearer-ignore-file",
		ConfigName:      "bearer-ignore-file",
		Value:           "bearer.ignore",
		Usage:           "Load bearer.ignore file from the specified path.",
		DisableInConfig: true,
	}
)

type IgnoreFlagGroup struct {
	BearerIgnoreFileFlag *Flag
}

type IgnoreOptions struct {
	BearerIgnoreFile string `mapstructure:"ignore_bearer_ignore_file" json:"ignore_bearer_ignore_file" yaml:"ignore_bearer_ignore_file"`
}

func NewIgnoreFlagGroup() *IgnoreFlagGroup {
	return &IgnoreFlagGroup{
		BearerIgnoreFileFlag: &BearerIgnoreFileFlag,
	}
}

func (f *IgnoreFlagGroup) Name() string {
	return "Ignore"
}

func (f *IgnoreFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.BearerIgnoreFileFlag,
	}
}

func (f *IgnoreFlagGroup) ToOptions() IgnoreOptions {
	return IgnoreOptions{
		BearerIgnoreFile: getString(f.BearerIgnoreFileFlag),
	}
}
