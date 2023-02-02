package subjects

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
	LineNumber     int      `json:"line_number,omitempty" yaml:"line_number"`
	RuleId         string   `json:"rule_id,omitempty" yaml:"rule_id"`
}

type RuleFailureSummary struct {
	CriticalRiskFailureCount int             `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFailureCount     int             `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFailureCount   int             `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFailureCount      int             `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	TriggeredRules           map[string]bool `json:"triggered_rules" yaml:"triggered_rules"`
}

type InventoryInput struct {
	Dataflow       *dataflow.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory  `json:"data_categories" yaml:"data_categories"`
}

type InventoryOutput struct {
	DataType    string `json:"name,omitempty" yaml:"name"`
	DataSubject string `json:"subject_name,omitempty" yaml:"subject_name"`
	LineNumber  int    `json:"line_number,omitempty" yaml:"line_number"`
}

type InventoryResult struct {
	DataSubject              string `json:"subject_name,omitempty" yaml:"subject_name"`
	DataType                 string `json:"name,omitempty" yaml:"name"`
	DetectionCount           int    `json:"detection_count" yaml:"detection_count"`
	CriticalRiskFailureCount int    `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFailureCount     int    `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFailureCount   int    `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFailureCount      int    `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	RulesPassedCount         int    `json:"rules_passed_count" yaml:"rules_passed_count"`
}

func BuildCsvString(dataflow *dataflow.DataFlow, config settings.Config) (*strings.Builder, error) {
	csvStr := &strings.Builder{}
	csvStr.WriteString("Subject,Data Types,Detection Count,Critical Risk Failure,High Risk Failure,Medium Risk Failure,Low Risk Failure,RulesPassed\n")
	result, err := GetOutput(dataflow, config)
	if err != nil {
		return csvStr, err
	}

	for _, item := range result {
		itemArr := []string{
			item.DataSubject,
			item.DataType,
			fmt.Sprint(item.DetectionCount),
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

	bar := output.GetProgressBar(len(config.Rules), config, "rules")

	result := make(map[string]InventoryResult)
	ruleFailures := make(map[string]RuleFailureSummary)
	localRuleCount := 0
	for _, rule := range config.Rules {
		if rule.Trigger == "local" {
			localRuleCount += 1
		}

		err := bar.Add(1)
		if err != nil {
			output.StdErrLogger().Msgf("Policy %s failed to write progress bar %e", rule.Id, err)
		}

		if !rule.PolicyType() {
			continue
		}

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
				key := buildKey(ruleOutputFailure.DataSubject, ruleOutputFailure.DataType)
				ruleFailure, ok := ruleFailures[key]
				if !ok {
					// key not found; create a new failure obj
					ruleFailure = RuleFailureSummary{
						CriticalRiskFailureCount: 0,
						HighRiskFailureCount:     0,
						MediumRiskFailureCount:   0,
						LowRiskFailureCount:      0,
						TriggeredRules:           make(map[string]bool),
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

				ruleFailure.TriggeredRules[ruleOutputFailure.RuleId] = true
				ruleFailures[key] = ruleFailure
			}
		}
	}

	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Compiling inventory report")
	}

	// get inventory result
	inventoryReportPolicy := config.Policies["inventory_report"]
	rs, err := rego.RunQuery(inventoryReportPolicy.Query,
		InventoryInput{
			Dataflow:       dataflow,
			DataCategories: db.DefaultWithContext(config.Scan.Context).DataCategories,
		},
		inventoryReportPolicy.Modules.ToRegoModules())
	if err != nil {
		return nil, err
	}
	if len(rs) > 0 {
		jsonRes, err := json.Marshal(rs)
		if err != nil {
			return nil, err
		}

		var inventoryOutput map[string][]InventoryOutput
		err = json.Unmarshal(jsonRes, &inventoryOutput)
		if err != nil {
			return nil, err
		}

		for _, outputItem := range inventoryOutput["report_items"] {
			key := buildKey(outputItem.DataSubject, outputItem.DataType)
			inventoryItem, ok := result[key]
			if !ok {
				// key not found, add a new item
				ruleFailure := ruleFailures[key]
				inventoryItem = InventoryResult{
					DataSubject:              outputItem.DataSubject,
					DataType:                 outputItem.DataType,
					CriticalRiskFailureCount: ruleFailure.CriticalRiskFailureCount,
					HighRiskFailureCount:     ruleFailure.HighRiskFailureCount,
					MediumRiskFailureCount:   ruleFailure.MediumRiskFailureCount,
					LowRiskFailureCount:      ruleFailure.LowRiskFailureCount,
					RulesPassedCount:         localRuleCount - len(ruleFailure.TriggeredRules),
				}
			}
			inventoryItem.DetectionCount += 1

			result[key] = inventoryItem
		}
	}

	return maps.Values(result), nil
}

func buildKey(dataSubject string, dataType string) string {
	return dataSubject + ":" + strings.ToUpper(dataType)
}
