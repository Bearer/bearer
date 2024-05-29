package flag

import flagtypes "github.com/bearer/bearer/pkg/flag/types"

type ruleFlagGroup struct{ flagGroupBase }

var RuleFlagGroup = &ruleFlagGroup{flagGroupBase{name: "Rule"}}

var (
	DisableDefaultRulesFlag = RuleFlagGroup.add(flagtypes.Flag{
		Name:       "disable-default-rules",
		ConfigName: "rule.disable-default-rules",
		Value:      false,
		Usage:      "Disables all default and built-in rules.",
	})
	SkipRuleFlag = RuleFlagGroup.add(flagtypes.Flag{
		Name:       "skip-rule",
		ConfigName: "rule.skip-rule",
		Value:      []string{},
		Usage:      "Specify the comma-separated ids of the rules you would like to skip. Runs all other rules.",
	})
	OnlyRuleFlag = RuleFlagGroup.add(flagtypes.Flag{
		Name:       "only-rule",
		ConfigName: "rule.only-rule",
		Value:      []string{},
		Usage:      "Specify the comma-separated ids of the rules you would like to run. Skips all other rules.",
	})
)

type RuleOptions struct {
	DisableDefaultRules bool            `mapstructure:"disable-default-rules" json:"disable-default-rules" yaml:"disable-default-rules"`
	SkipRule            map[string]bool `mapstructure:"skip-rule" json:"skip-rule" yaml:"skip-rule"`
	OnlyRule            map[string]bool `mapstructure:"only-rule" json:"only-rule" yaml:"only-rule"`
}

func (ruleFlagGroup) SetOptions(options *flagtypes.Options, args []string) error {
	options.RuleOptions = flagtypes.RuleOptions{
		DisableDefaultRules: getBool(DisableDefaultRulesFlag),
		SkipRule:            argsToMap(SkipRuleFlag),
		OnlyRule:            argsToMap(OnlyRuleFlag),
	}

	return nil
}
