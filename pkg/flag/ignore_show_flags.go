package flag

var (
	AllFlag = Flag{
		Name:       "all",
		ConfigName: "ignore_show.all",
		Value:      false,
		Usage:      "Show all ignored fingerprints.",
	}
)

type IgnoreShowFlagGroup struct {
	AllFlag *Flag
}

type IgnoreShowOptions struct {
	All bool `mapstructure:"all" json:"all" yaml:"all"`
}

func NewIgnoreShowFlagGroup() *IgnoreShowFlagGroup {
	return &IgnoreShowFlagGroup{
		AllFlag: &AllFlag,
	}
}

func (f *IgnoreShowFlagGroup) Name() string {
	return "IgnoreShow"
}

func (f *IgnoreShowFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.AllFlag,
	}
}

func (f *IgnoreShowFlagGroup) ToOptions() IgnoreShowOptions {
	return IgnoreShowOptions{
		All: getBool(f.AllFlag),
	}
}
