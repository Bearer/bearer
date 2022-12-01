package flag

var (
	SkipPolicyFlag = Flag{
		Name:       "skip-policy",
		ConfigName: "policy.skip-policy",
		Value:      []string{},
		Usage:      "specify the comma separated ids of the policies you would like to skip. Runs all other policies.",
	}
	OnlyPolicyFlag = Flag{
		Name:       "only-policy",
		ConfigName: "policy.only-policy",
		Value:      []string{},
		Usage:      "specify the comma separated ids of the policies you would like to run. Skips all other policies.",
	}
)

type PolicyFlagGroup struct {
	SkipPolicyFlag *Flag
	OnlyPolicyFlag *Flag
}

type PolicyOptions struct {
	SkipPolicy map[string]bool `mapstructure:"skip-policy" json:"skip-policy" yaml:"skip-policy"`
	OnlyPolicy map[string]bool `mapstructure:"only-policy" json:"only-policy" yaml:"only-policy"`
}

func NewPolicyFlagGroup() *PolicyFlagGroup {
	return &PolicyFlagGroup{
		SkipPolicyFlag: &SkipPolicyFlag,
		OnlyPolicyFlag: &OnlyPolicyFlag,
	}
}

func (f *PolicyFlagGroup) Name() string {
	return "Policy"
}

func (f *PolicyFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.SkipPolicyFlag,
		f.OnlyPolicyFlag,
	}
}

func (f *PolicyFlagGroup) ToOptions(args []string) PolicyOptions {
	return PolicyOptions{
		SkipPolicy: argsToMap(f.SkipPolicyFlag),
		OnlyPolicy: argsToMap(f.OnlyPolicyFlag),
	}
}

func argsToMap(flag *Flag) map[string]bool {
	strSlice := getStringSlice(flag)

	result := make(map[string]bool)
	for _, str := range strSlice {
		result[str] = true
	}

	return result
}
