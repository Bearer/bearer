package ruleset

import (
	"fmt"
	"maps"
	"slices"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/customdetectors"
	"github.com/bearer/bearer/pkg/util/set"
)

type RuleType int

const (
	RuleTypeTopLevel RuleType = iota
	RuleTypeShared
	RuleTypeBuiltin
	RuleTypeAuxiliary
)

type Set struct {
	rules     []*Rule
	rulesByID map[string]*Rule
}

type Rule struct {
	index         int
	id            string
	ruleType      RuleType
	sanitizerRule *Rule
	patterns      []settings.RulePattern
}

func New(languageID string, settingsRules map[string]*settings.Rule) (*Set, error) {
	languageRules := getLanguageRules(settingsRules, languageID)
	triggerRuleIDs := getTriggerRuleIDs(languageRules)

	rulesByID := make(map[string]*Rule)
	var rules []*Rule

	for _, rule := range builtinRules {
		if rulesByID[rule.id] != nil {
			return nil, fmt.Errorf("duplicate built-in rule '%s'", rule.id)
		}

		rules = append(rules, rule)
		rulesByID[rule.id] = rule
	}

	for _, settingsRule := range languageRules {
		rule := &Rule{
			index:    len(rules),
			id:       settingsRule.Id,
			ruleType: getRuleType(triggerRuleIDs, settingsRule),
			patterns: settingsRule.Patterns,
		}

		if rulesByID[rule.id] != nil {
			return nil, fmt.Errorf("duplicate rule '%s'", rule.id)
		}

		rules = append(rules, rule)
		rulesByID[rule.id] = rule
	}

	for _, rule := range rules {
		if rule.ruleType == RuleTypeBuiltin {
			continue
		}

		settingsRule := settingsRules[rule.id]
		if settingsRule.SanitizerRuleID == "" {
			continue
		}

		sanitizerRule := rulesByID[settingsRule.SanitizerRuleID]
		if sanitizerRule == nil {
			return nil, fmt.Errorf("invalid rule id for sanitizer '%s'", settingsRule.SanitizerRuleID)
		}

		rule.sanitizerRule = sanitizerRule
	}

	return &Set{
		rules:     rules,
		rulesByID: rulesByID,
	}, nil
}

func getLanguageRules(settingsRules map[string]*settings.Rule, languageID string) []*settings.Rule {
	var result []*settings.Rule

	ruleIDs := slices.Sorted(maps.Keys(settingsRules))

	for _, ruleID := range ruleIDs {
		settingsRule := settingsRules[ruleID]
		if slices.Contains(settingsRule.Languages, languageID) {
			result = append(result, settingsRule)
		}
	}

	return result
}

func getTriggerRuleIDs(languageRules []*settings.Rule) set.Set[string] {
	triggerRuleIDs := set.New[string]()

	for _, settingsRule := range languageRules {
		triggerRuleIDs.AddAll(settingsRule.Trigger.RequiredDetections)
	}

	return triggerRuleIDs
}

func getRuleType(triggerRuleIDs set.Set[string], settingsRule *settings.Rule) RuleType {
	switch {
	case settingsRule.Type == customdetectors.TypeShared:
		return RuleTypeShared
	case !settingsRule.IsAuxilary || triggerRuleIDs.Has(settingsRule.Id):
		return RuleTypeTopLevel
	default:
		return RuleTypeAuxiliary
	}
}

func (set *Set) RuleByIndex(idx uint64) (*Rule, error) {
	return set.Rules()[idx], nil
}

func (set *Set) RuleByID(id string) (*Rule, error) {
	rule, exists := set.rulesByID[id]
	if !exists {
		return nil, fmt.Errorf("invalid rule id '%s'", id)
	}

	return rule, nil
}

func (set *Set) Rules() []*Rule {
	return set.rules
}

func (rule *Rule) Index() int {
	return rule.index
}

func (rule *Rule) ID() string {
	return rule.id
}

func (rule *Rule) Type() RuleType {
	return rule.ruleType
}

func (rule *Rule) SanitizerRule() *Rule {
	return rule.sanitizerRule
}

func (rule *Rule) Patterns() []settings.RulePattern {
	return rule.patterns
}
