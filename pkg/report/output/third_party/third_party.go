package third_party

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/commands/process/settings"
	"github.com/bearer/curio/pkg/util/output"
	"github.com/bearer/curio/pkg/util/rego"
	"golang.org/x/exp/maps"

	"github.com/bearer/curio/pkg/report/output/dataflow"
	"github.com/bearer/curio/pkg/report/output/summary"
)

type RuleInput struct {
	RuleId         string             `json:"rule_id" yaml:"rule_id"`
	Rule           *settings.Rule     `json:"rule" yaml:"rule"`
	Dataflow       *dataflow.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory  `json:"data_categories" yaml:"data_categories"`
}

type RuleOutput struct {
	DataType       string   `json:"name,omitempty" yaml:"name"`
	CategoryGroups []string `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	DataSubject    string   `json:"subject_name,omitempty" yaml:"subject_name"`
	RuleId         string   `json:"rule_id,omitempty" yaml:"rule_id"`
	ThirdParty     string   `json:"third_party,omitempty" yaml:"third_party"`
}

type RuleFailureSummary struct {
	DataSubject              string          `json:"subject_name,omitempty" yaml:"subject_name"`
	DataTypes                map[string]bool `json:"data_types" yaml:"data_types"`
	CriticalRiskFailureCount int             `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFailureCount     int             `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFailureCount   int             `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFailureCount      int             `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	TriggeredRulesCount      map[string]int  `json:"triggered_rules" yaml:"triggered_rules"`
}

type InventoryResult struct {
	ThirdParty               string   `json:"third_party,omitempty" yaml:"third_party"`
	DataSubject              string   `json:"subject_name,omitempty" yaml:"subject_name"`
	DataTypes                []string `json:"data_types,omitempty" yaml:"data_types"`
	CriticalRiskFailureCount int      `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFailureCount     int      `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFailureCount   int      `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFailureCount      int      `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	RulesPassedCount         int      `json:"rules_passed_count" yaml:"rules_passed_count"`
}

func BuildCsvString(dataflow *dataflow.DataFlow, config settings.Config) (*strings.Builder, error) {
	csvStr := &strings.Builder{}
	csvStr.WriteString("Third Party,Subject,Data Types,Critical Risk Failure,High Risk Failure,Medium Risk Failure,Low Risk Failure,RulesPassed\n")
	result, err := GetOutput(dataflow, config)
	if err != nil {
		return csvStr, err
	}

	for _, item := range result {
		itemArr := []string{
			item.ThirdParty,
			item.DataSubject,
			"\"" + strings.Join(item.DataTypes, ",") + "\"",
			fmt.Sprint(item.CriticalRiskFailureCount),
			fmt.Sprint(item.HighRiskFailureCount),
			fmt.Sprint(item.MediumRiskFailureCount),
			fmt.Sprint(item.LowRiskFailureCount),
			fmt.Sprint(item.RulesPassedCount),
		}
		csvStr.WriteString(strings.Join(itemArr, ",") + "\n")
	}

	return csvStr, nil
}

func GetOutput(dataflow *dataflow.DataFlow, config settings.Config) ([]InventoryResult, error) {
	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Evaluating rules")
	}

	thirdPartyFailures := make(map[string]map[string]RuleFailureSummary)
	triggeredRules := make(map[string]bool)
	thirdPartyRuleCount := make(map[string]int)
	for _, rule := range config.Rules {
		if !rule.PolicyType() {
			continue
		}

		if rule.AssociatedRecipe == "" {
			// no associated recipe
			continue
		}

		thirdPartyRuleCount[rule.AssociatedRecipe] += 1
		policy := config.Policies[rule.Type]

		// Create a prepared query that can be evaluated.
		rs, err := rego.RunQuery(policy.Query,
			RuleInput{
				RuleId:         rule.Id,
				Rule:           rule,
				Dataflow:       dataflow,
				DataCategories: db.DefaultWithContext(config.Scan.Context).DataCategories,
			},
			policy.Modules.ToRegoModules())
		if err != nil {
			return nil, err
		}

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs)
			if err != nil {
				return nil, err
			}

			var ruleOutput map[string][]RuleOutput
			err = json.Unmarshal(jsonRes, &ruleOutput)
			if err != nil {
				return nil, err
			}

			for _, ruleOutputFailure := range ruleOutput["local_rule_failure"] {
				thirdParty, ok := thirdPartyFailures[ruleOutputFailure.ThirdParty]
				if !ok {
					// third party key not found; create empty map
					thirdParty = make(map[string]RuleFailureSummary)
					thirdPartyFailures[ruleOutputFailure.ThirdParty] = thirdParty
				}

				ruleFailure, ok := thirdParty[ruleOutputFailure.DataSubject]
				if !ok {
					// data subject key not found; create a new failure obj
					ruleFailure = RuleFailureSummary{
						DataSubject:              ruleOutputFailure.DataSubject,
						DataTypes:                make(map[string]bool),
						CriticalRiskFailureCount: 0,
						HighRiskFailureCount:     0,
						MediumRiskFailureCount:   0,
						LowRiskFailureCount:      0,
						TriggeredRulesCount:      make(map[string]int),
					}
				}

				// count severity
				switch summary.FindHighestSeverity(ruleOutputFailure.CategoryGroups, rule.Severity) {
				case "critical":
					ruleFailure.CriticalRiskFailureCount += 1
				case "high":
					ruleFailure.HighRiskFailureCount += 1
				case "medium":
					ruleFailure.MediumRiskFailureCount += 1
				case "low":
					ruleFailure.LowRiskFailureCount += 1
				}

				// add data type to map
				ruleFailure.DataTypes[ruleOutputFailure.DataType] = true

				// increment triggered rules where necessary
				_, ok = triggeredRules[ruleOutputFailure.RuleId+":"+ruleOutputFailure.ThirdParty]
				if !ok {
					triggeredRules[ruleOutputFailure.RuleId+":"+ruleOutputFailure.ThirdParty] = true
					ruleFailure.TriggeredRulesCount[ruleOutputFailure.ThirdParty] += 1
				}

				thirdParty[ruleOutputFailure.DataSubject] = ruleFailure
			}
		}
	}

	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Compiling third party report")
	}

	// get inventory result
	result := []InventoryResult{}
	for _, component := range dataflow.Components {
		if component.SubType != "third_party" {
			continue
		}

		var inventoryItem InventoryResult
		thirdPartyFailure, ok := thirdPartyFailures[component.Name]
		if !ok {
			// no failures, therefore no associated data subjects
			inventoryItem = InventoryResult{
				ThirdParty:               component.Name,
				DataSubject:              "Unknown",
				DataTypes:                []string{"Unknown"},
				CriticalRiskFailureCount: 0,
				HighRiskFailureCount:     0,
				MediumRiskFailureCount:   0,
				LowRiskFailureCount:      0,
				RulesPassedCount:         0,
			}
		}
		for _, ruleFailure := range thirdPartyFailure {
			inventoryItem = InventoryResult{
				ThirdParty:               component.Name,
				DataSubject:              ruleFailure.DataSubject,
				DataTypes:                maps.Keys(ruleFailure.DataTypes),
				CriticalRiskFailureCount: ruleFailure.CriticalRiskFailureCount,
				HighRiskFailureCount:     ruleFailure.HighRiskFailureCount,
				MediumRiskFailureCount:   ruleFailure.MediumRiskFailureCount,
				LowRiskFailureCount:      ruleFailure.LowRiskFailureCount,
				RulesPassedCount:         thirdPartyRuleCount[component.Name] - ruleFailure.TriggeredRulesCount[component.Name],
			}
		}

		result = append(result, inventoryItem)
	}

	return result, nil
}
