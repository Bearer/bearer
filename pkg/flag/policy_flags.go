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
	SkipPolicy []string `json:"skip_policy"`
	OnlyPolicy []string `json:"only_policy"`
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
		SkipPolicy: getStringSlice(f.SkipPolicyFlag),
		OnlyPolicy: getStringSlice(f.OnlyPolicyFlag),
	}
}
