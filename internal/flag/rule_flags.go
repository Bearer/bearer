package flag

var (
	DisableDefaultRulesFlag = Flag{
		Name:       "disable-default-rules",
		ConfigName: "rule.disable-default-rules",
		Value:      false,
		Usage:      "Disables all default and built-in rules.",
	}
	SkipRuleFlag = Flag{
		Name:       "skip-rule",
		ConfigName: "rule.skip-rule",
		Value:      []string{},
		Usage:      "Specify the comma-separated ids of the rules you would like to skip. Runs all other rules.",
	}
	OnlyRuleFlag = Flag{
		Name:       "only-rule",
		ConfigName: "rule.only-rule",
		Value:      []string{},
		Usage:      "Specify the comma-separated ids of the rules you would like to run. Skips all other rules.",
	}
)

type RuleFlagGroup struct {
	DisableDefaultRulesFlag *Flag
	SkipRuleFlag            *Flag
	OnlyRuleFlag            *Flag
}

type RuleOptions struct {
	DisableDefaultRules bool            `mapstructure:"disable-default-rules" json:"disable-default-rules" yaml:"disable-default-rules"`
	SkipRule            map[string]bool `mapstructure:"skip-rule" json:"skip-rule" yaml:"skip-rule"`
	OnlyRule            map[string]bool `mapstructure:"only-rule" json:"only-rule" yaml:"only-rule"`
}

func NewRuleFlagGroup() *RuleFlagGroup {
	return &RuleFlagGroup{
		DisableDefaultRulesFlag: &DisableDefaultRulesFlag,
		SkipRuleFlag:            &SkipRuleFlag,
		OnlyRuleFlag:            &OnlyRuleFlag,
	}
}

func (f *RuleFlagGroup) Name() string {
	return "Rule"
}

func (f *RuleFlagGroup) Flags() []*Flag {
	return []*Flag{
		f.DisableDefaultRulesFlag,
		f.SkipRuleFlag,
		f.OnlyRuleFlag,
	}
}

func (f *RuleFlagGroup) ToOptions(args []string) RuleOptions {
	return RuleOptions{
		DisableDefaultRules: getBool(f.DisableDefaultRulesFlag),
		SkipRule:            argsToMap(f.SkipRuleFlag),
		OnlyRule:            argsToMap(f.OnlyRuleFlag),
	}
}
