package privacy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/bearer/bearer/pkg/util/rego"
	"github.com/hhatto/gocloc"
	"golang.org/x/exp/maps"

	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/security"
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
	ThirdParty     string   `json:"third_party,omitempty" yaml:"third_party"`
}

type RuleFailureSummary struct {
	DataSubject              string          `json:"subject_name,omitempty" yaml:"subject_name"`
	DataTypes                map[string]bool `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	CriticalRiskFailureCount int             `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFailureCount     int             `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFailureCount   int             `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFailureCount      int             `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	TriggeredRules           map[string]bool `json:"triggered_rules" yaml:"triggered_rules"`
}

type Input struct {
	Dataflow       *dataflow.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory  `json:"data_categories" yaml:"data_categories"`
}

type Output struct {
	DataType    string `json:"name,omitempty" yaml:"name"`
	DataSubject string `json:"subject_name,omitempty" yaml:"subject_name"`
	LineNumber  int    `json:"line_number,omitempty" yaml:"line_number"`
}

type Subject struct {
	DataSubject              string `json:"subject_name,omitempty" yaml:"subject_name"`
	DataType                 string `json:"name,omitempty" yaml:"name"`
	DetectionCount           int    `json:"detection_count" yaml:"detection_count"`
	CriticalRiskFailureCount int    `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFailureCount     int    `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFailureCount   int    `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFailureCount      int    `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	RulesPassedCount         int    `json:"rules_passed_count" yaml:"rules_passed_count"`
}

type ThirdPartyRuleCounter struct {
	RuleIds         map[string]bool
	Count           int
	SubjectFailures map[string]map[string]bool
}

type ThirdParty struct {
	ThirdParty               string   `json:"third_party,omitempty" yaml:"third_party"`
	DataSubject              string   `json:"subject_name,omitempty" yaml:"subject_name"`
	DataTypes                []string `json:"data_types,omitempty" yaml:"data_types"`
	CriticalRiskFailureCount int      `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFailureCount     int      `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFailureCount   int      `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFailureCount      int      `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	RulesPassedCount         int      `json:"rules_passed_count" yaml:"rules_passed_count"`
}

type Report struct {
	Subjects   []Subject    `json:"subjects,omitempty" yaml:"subjects"`
	ThirdParty []ThirdParty `json:"third_party,omitempty" yaml:"third_party"`
}

func BuildCsvString(dataflow *dataflow.DataFlow, lineOfCodeOutput *gocloc.Result, config settings.Config) (*strings.Builder, error) {
	csvStr := &strings.Builder{}
	csvStr.WriteString("\nSubject,Data Types,Detection Count,Critical Risk Failure,High Risk Failure,Medium Risk Failure,Low Risk Failure,Rules Passed\n")
	result, _, _, err := GetOutput(dataflow, lineOfCodeOutput, config)
	if err != nil {
		return csvStr, err
	}

	for _, subject := range result.Subjects {
		subjectArr := []string{
			subject.DataSubject,
			subject.DataType,
			fmt.Sprint(subject.DetectionCount),
			fmt.Sprint(subject.CriticalRiskFailureCount),
			fmt.Sprint(subject.HighRiskFailureCount),
			fmt.Sprint(subject.MediumRiskFailureCount),
			fmt.Sprint(subject.LowRiskFailureCount),
			fmt.Sprint(subject.RulesPassedCount),
		}
		csvStr.WriteString(strings.Join(subjectArr, ",") + "\n")
	}

	csvStr.WriteString("\n")
	csvStr.WriteString("Third Party,Subject,Data Types,Critical Risk Failure,High Risk Failure,Medium Risk Failure,Low Risk Failure,Rules Passed\n")

	for _, thirdParty := range result.ThirdParty {
		thirdPartyArr := []string{
			thirdParty.ThirdParty,
			thirdParty.DataSubject,
			"\"" + strings.Join(thirdParty.DataTypes, ",") + "\"",
			fmt.Sprint(thirdParty.CriticalRiskFailureCount),
			fmt.Sprint(thirdParty.HighRiskFailureCount),
			fmt.Sprint(thirdParty.MediumRiskFailureCount),
			fmt.Sprint(thirdParty.LowRiskFailureCount),
			fmt.Sprint(thirdParty.RulesPassedCount),
		}
		csvStr.WriteString(strings.Join(thirdPartyArr, ",") + "\n")
	}

	return csvStr, nil
}

func GetOutput(dataflow *dataflow.DataFlow, lineOfCodeOutput *gocloc.Result, config settings.Config) (*Report, *gocloc.Result, *dataflow.DataFlow, error) {
	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Evaluating rules")
	}

	bar := output.GetProgressBar(len(config.Rules), config, "rules")

	subjectRuleFailures := make(map[string]RuleFailureSummary)
	thirdPartyRuleFailures := make(map[string]map[string]RuleFailureSummary)

	localRuleCounter := 0
	thirdPartyRulesCounter := make(map[string]ThirdPartyRuleCounter)

	for _, rule := range config.Rules {
		// increment counters
		if rule.IsLocal {
			localRuleCounter += 1
		}

		if rule.AssociatedRecipe != "" {
			thirdPartyRuleCounter, ok := thirdPartyRulesCounter[rule.AssociatedRecipe]
			if !ok {
				thirdPartyRuleCounter = ThirdPartyRuleCounter{
					RuleIds:         make(map[string]bool),
					SubjectFailures: make(map[string]map[string]bool),
				}
			}

			thirdPartyRuleCounter.Count += 1
			thirdPartyRuleCounter.RuleIds[rule.Id] = true

			thirdPartyRulesCounter[rule.AssociatedRecipe] = thirdPartyRuleCounter
		}

		err := bar.Add(1)
		if err != nil {
			output.StdErrLogger().Msgf("Policy %s failed to write progress bar %s", rule.Id, err)
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
			return nil, nil, nil, err
		}

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs)
			if err != nil {
				return nil, nil, nil, err
			}

			var ruleOutput map[string][]RuleOutput
			err = json.Unmarshal(jsonRes, &ruleOutput)
			if err != nil {
				return nil, nil, nil, err
			}

			for _, ruleOutputFailure := range ruleOutput["local_rule_failure"] {
				ruleSeverity := security.CalculateSeverity(ruleOutputFailure.CategoryGroups, rule.Severity, true)

				key := buildKey(ruleOutputFailure.DataSubject, ruleOutputFailure.DataType)
				subjectRuleFailure, ok := subjectRuleFailures[key]
				if !ok {
					// key not found; create a new failure obj
					subjectRuleFailure = RuleFailureSummary{
						CriticalRiskFailureCount: 0,
						HighRiskFailureCount:     0,
						MediumRiskFailureCount:   0,
						LowRiskFailureCount:      0,
						TriggeredRules:           make(map[string]bool),
					}
				}
				// count severity
				switch ruleSeverity {
				case "critical":
					subjectRuleFailure.CriticalRiskFailureCount += 1
				case "high":
					subjectRuleFailure.HighRiskFailureCount += 1
				case "medium":
					subjectRuleFailure.MediumRiskFailureCount += 1
				case "low":
					subjectRuleFailure.LowRiskFailureCount += 1
				}

				subjectRuleFailure.TriggeredRules[ruleOutputFailure.RuleId] = true
				subjectRuleFailures[key] = subjectRuleFailure

				// update third party failures

				if rule.AssociatedRecipe == "" {
					continue
				}

				thirdPartyFailure, ok := thirdPartyRuleFailures[ruleOutputFailure.ThirdParty]
				if !ok {
					// third party key not found; create empty map
					thirdPartyFailure = make(map[string]RuleFailureSummary)
					thirdPartyRuleFailures[ruleOutputFailure.ThirdParty] = thirdPartyFailure
				}
				thirdPartyDataSubject, ok := thirdPartyFailure[ruleOutputFailure.DataSubject]
				if !ok {
					// data subject key not found; create a new failure obj
					thirdPartyDataSubject = RuleFailureSummary{
						DataSubject:              ruleOutputFailure.DataSubject,
						DataTypes:                make(map[string]bool),
						CriticalRiskFailureCount: 0,
						HighRiskFailureCount:     0,
						MediumRiskFailureCount:   0,
						LowRiskFailureCount:      0,
					}
				}

				// count severity
				switch ruleSeverity {
				case "critical":
					thirdPartyDataSubject.CriticalRiskFailureCount += 1
				case "high":
					thirdPartyDataSubject.HighRiskFailureCount += 1
				case "medium":
					thirdPartyDataSubject.MediumRiskFailureCount += 1
				case "low":
					thirdPartyDataSubject.LowRiskFailureCount += 1
				}

				// add data type to map
				thirdPartyDataSubject.DataTypes[ruleOutputFailure.DataType] = true
				thirdPartyRuleFailures[ruleOutputFailure.ThirdParty][ruleOutputFailure.DataSubject] = thirdPartyDataSubject

				// increment counter
				thirdPartyRuleCounter := thirdPartyRulesCounter[rule.AssociatedRecipe]
				subjectFailure := thirdPartyRuleCounter.SubjectFailures[ruleOutputFailure.DataSubject]
				if !ok {
					subjectFailure = make(map[string]bool)
				}
				subjectFailure[ruleOutputFailure.RuleId] = true
				thirdPartyRuleCounter.SubjectFailures[ruleOutputFailure.DataSubject] = subjectFailure
			}
		}
	}

	if !config.Scan.Quiet {
		output.StdErrLogger().Msgf("Compiling privacy report")
	}

	// get inventory result
	subjectInventory := make(map[string]Subject)
	privacyReportPolicy := config.Policies["privacy_report"]
	rs, err := rego.RunQuery(privacyReportPolicy.Query,
		Input{
			Dataflow:       dataflow,
			DataCategories: db.DefaultWithContext(config.Scan.Context).DataCategories,
		},
		privacyReportPolicy.Modules.ToRegoModules(),
	)

	if err != nil {
		return nil, nil, nil, err
	}

	if len(rs) > 0 {
		jsonRes, err := json.Marshal(rs)
		if err != nil {
			return nil, nil, nil, err
		}

		var outputItems map[string][]Output
		err = json.Unmarshal(jsonRes, &outputItems)
		if err != nil {
			return nil, nil, nil, err
		}

		for _, outputItem := range outputItems["items"] {
			key := buildKey(outputItem.DataSubject, outputItem.DataType)
			subject, ok := subjectInventory[key]
			if !ok {
				// key not found, add a new item
				ruleFailure := subjectRuleFailures[key]
				subject = Subject{
					DataSubject:              outputItem.DataSubject,
					DataType:                 outputItem.DataType,
					CriticalRiskFailureCount: ruleFailure.CriticalRiskFailureCount,
					HighRiskFailureCount:     ruleFailure.HighRiskFailureCount,
					MediumRiskFailureCount:   ruleFailure.MediumRiskFailureCount,
					LowRiskFailureCount:      ruleFailure.LowRiskFailureCount,
					RulesPassedCount:         localRuleCounter - len(ruleFailure.TriggeredRules),
				}
			}
			subject.DetectionCount += 1
			subjectInventory[key] = subject
		}
	}

	var thirdPartyInventory []ThirdParty
	for _, component := range dataflow.Components {
		if component.SubType != "third_party" {
			continue
		}

		thirdPartyFailure, ok := thirdPartyRuleFailures[component.Name]
		if !ok {
			// no failures, therefore no associated data subjects
			thirdPartyInventory = append(thirdPartyInventory, ThirdParty{
				ThirdParty:               component.Name,
				DataSubject:              "Unknown",
				DataTypes:                []string{"Unknown"},
				CriticalRiskFailureCount: 0,
				HighRiskFailureCount:     0,
				MediumRiskFailureCount:   0,
				LowRiskFailureCount:      0,
				RulesPassedCount:         0,
			})
		}

		for _, ruleFailure := range thirdPartyFailure {
			thirdPartyInventory = append(thirdPartyInventory, ThirdParty{
				ThirdParty:               component.Name,
				DataSubject:              ruleFailure.DataSubject,
				DataTypes:                maps.Keys(ruleFailure.DataTypes),
				CriticalRiskFailureCount: ruleFailure.CriticalRiskFailureCount,
				HighRiskFailureCount:     ruleFailure.HighRiskFailureCount,
				MediumRiskFailureCount:   ruleFailure.MediumRiskFailureCount,
				LowRiskFailureCount:      ruleFailure.LowRiskFailureCount,
				RulesPassedCount:         thirdPartyRulesCounter[component.Name].Count - len(thirdPartyRulesCounter[component.Name].SubjectFailures[ruleFailure.DataSubject]),
			})
		}
	}

	return &Report{
		Subjects:   maps.Values(subjectInventory),
		ThirdParty: thirdPartyInventory,
	}, lineOfCodeOutput, dataflow, nil
}

func buildKey(dataSubject string, dataType string) string {
	return dataSubject + ":" + strings.ToUpper(dataType)
}
