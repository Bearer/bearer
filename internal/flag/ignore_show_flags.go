package flag

type ignoreShowFlagGroup struct{ flagGroupBase }

var IgnoreShowFlagGroup = &ignoreShowFlagGroup{flagGroupBase{name: "Ignore Show"}}

var (
	AllFlag = IgnoreShowFlagGroup.add(Flag{
		Name:       "all",
		ConfigName: "ignore_show.all",
		Value:      false,
		Usage:      "Show all ignored fingerprints.",
	})
)

type IgnoreShowOptions struct {
	All bool `mapstructure:"all" json:"all" yaml:"all"`
}

func (ignoreShowFlagGroup) SetOptions(options *Options, args []string) error {
	options.IgnoreShowOptions = IgnoreShowOptions{
		All: getBool(AllFlag),
	}

	return nil
}
