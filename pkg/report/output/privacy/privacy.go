package privacy

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/bearer/bearer/pkg/classification/db"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	globaltypes "github.com/bearer/bearer/pkg/types"
	"github.com/bearer/bearer/pkg/util/output"
	"github.com/bearer/bearer/pkg/util/progressbar"
	"github.com/bearer/bearer/pkg/util/rego"

	"github.com/bearer/bearer/pkg/report/output/privacy/types"
	"github.com/bearer/bearer/pkg/report/output/security"
	outputtypes "github.com/bearer/bearer/pkg/report/output/types"
)

type RuleInput struct {
	RuleId         string                `json:"rule_id" yaml:"rule_id"`
	Rule           *settings.Rule        `json:"rule" yaml:"rule"`
	Dataflow       *outputtypes.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory     `json:"data_categories" yaml:"data_categories"`
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
	CriticalRiskFindingCount int             `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFindingCount     int             `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFindingCount   int             `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFindingCount      int             `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	TriggeredRules           map[string]bool `json:"triggered_rules" yaml:"triggered_rules"`
}

type Input struct {
	Dataflow       *outputtypes.DataFlow `json:"dataflow" yaml:"dataflow"`
	DataCategories []db.DataCategory     `json:"data_categories" yaml:"data_categories"`
}

type Output struct {
	DataType    string `json:"name,omitempty" yaml:"name"`
	DataSubject string `json:"subject_name,omitempty" yaml:"subject_name"`
	LineNumber  int    `json:"line_number,omitempty" yaml:"line_number"`
}

type ThirdPartyRuleCounter struct {
	RuleIds         map[string]bool
	Count           int
	SubjectFailures map[string]map[string]bool
}

const PLACEHOLDER_VALUE = "Unknown"

func BuildCsvString(reportData *outputtypes.ReportData, config settings.Config) (*strings.Builder, error) {
	csvStr := &strings.Builder{}
	csvStr.WriteString("\nSubject,Data Types,Detection Count,Critical Risk Finding,High Risk Finding,Medium Risk Finding,Low Risk Finding,Rules Passed\n")

	for _, subject := range reportData.PrivacyReport.Subjects {
		subjectArr := []string{
			subject.DataSubject,
			subject.DataType,
			fmt.Sprint(subject.DetectionCount),
			fmt.Sprint(subject.CriticalRiskFindingCount),
			fmt.Sprint(subject.HighRiskFindingCount),
			fmt.Sprint(subject.MediumRiskFindingCount),
			fmt.Sprint(subject.LowRiskFindingCount),
			fmt.Sprint(subject.RulesPassedCount),
		}
		csvStr.WriteString(strings.Join(subjectArr, ",") + "\n")
	}

	csvStr.WriteString("\n")
	csvStr.WriteString("Third Party,Subject,Data Types,Critical Risk Finding,High Risk Finding,Medium Risk Finding,Low Risk Finding,Rules Passed\n")

	for _, thirdParty := range reportData.PrivacyReport.ThirdParty {
		thirdPartyArr := []string{
			thirdParty.ThirdParty,
			thirdParty.DataSubject,
			"\"" + strings.Join(thirdParty.DataTypes, ",") + "\"",
			fmt.Sprint(thirdParty.CriticalRiskFindingCount),
			fmt.Sprint(thirdParty.HighRiskFindingCount),
			fmt.Sprint(thirdParty.MediumRiskFindingCount),
			fmt.Sprint(thirdParty.LowRiskFindingCount),
			fmt.Sprint(thirdParty.RulesPassedCount),
		}
		csvStr.WriteString(strings.Join(thirdPartyArr, ",") + "\n")
	}

	return csvStr, nil
}

func AddReportData(reportData *outputtypes.ReportData, config settings.Config) error {
	if !config.Scan.Quiet {
		output.StdErrLog("Evaluating rules")
	}

	bar := progressbar.GetProgressBar(len(config.Rules), &config)

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
			output.StdErrLog(fmt.Sprintf("Policy %s failed to write progress bar %s", rule.Id, err))
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
				Dataflow:       reportData.Dataflow,
				DataCategories: db.DefaultWithContext(config.Scan.Context).DataCategories,
			},
			policy.Modules.ToRegoModules())
		if err != nil {
			return err
		}

		if len(rs) > 0 {
			jsonRes, err := json.Marshal(rs)
			if err != nil {
				return err
			}

			var ruleOutput map[string][]RuleOutput
			err = json.Unmarshal(jsonRes, &ruleOutput)
			if err != nil {
				return err
			}

			for _, ruleOutputFailure := range ruleOutput["local_rule_failure"] {
				ruleSeverity := security.CalculateSeverity(ruleOutputFailure.CategoryGroups, rule.GetSeverity(), true)

				key := buildKey(ruleOutputFailure.DataSubject, ruleOutputFailure.DataType)
				subjectRuleFailure, ok := subjectRuleFailures[key]
				if !ok {
					// key not found; create a new failure obj
					subjectRuleFailure = RuleFailureSummary{
						CriticalRiskFindingCount: 0,
						HighRiskFindingCount:     0,
						MediumRiskFindingCount:   0,
						LowRiskFindingCount:      0,
						TriggeredRules:           make(map[string]bool),
					}
				}

				// count severity
				switch ruleSeverity.DisplaySeverity {
				case globaltypes.LevelCritical:
					subjectRuleFailure.CriticalRiskFindingCount += 1
				case globaltypes.LevelHigh:
					subjectRuleFailure.HighRiskFindingCount += 1
				case globaltypes.LevelMedium:
					subjectRuleFailure.MediumRiskFindingCount += 1
				case globaltypes.LevelLow:
					subjectRuleFailure.LowRiskFindingCount += 1
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
						CriticalRiskFindingCount: 0,
						HighRiskFindingCount:     0,
						MediumRiskFindingCount:   0,
						LowRiskFindingCount:      0,
					}
				}

				// count severity
				switch ruleSeverity.DisplaySeverity {
				case globaltypes.LevelCritical:
					thirdPartyDataSubject.CriticalRiskFindingCount += 1
				case globaltypes.LevelHigh:
					thirdPartyDataSubject.HighRiskFindingCount += 1
				case globaltypes.LevelMedium:
					thirdPartyDataSubject.MediumRiskFindingCount += 1
				case globaltypes.LevelLow:
					thirdPartyDataSubject.LowRiskFindingCount += 1
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
		output.StdErrLog("Compiling privacy report")
	}

	// get inventory result
	subjectInventory := make(map[string]types.Subject)
	privacyReportPolicy := config.Policies["privacy_report"]
	rs, err := rego.RunQuery(privacyReportPolicy.Query,
		Input{
			Dataflow:       reportData.Dataflow,
			DataCategories: db.DefaultWithContext(config.Scan.Context).DataCategories,
		},
		privacyReportPolicy.Modules.ToRegoModules(),
	)

	if err != nil {
		return err
	}

	if len(rs) > 0 {
		jsonRes, err := json.Marshal(rs)
		if err != nil {
			return err
		}

		var outputItems map[string][]Output
		err = json.Unmarshal(jsonRes, &outputItems)
		if err != nil {
			return err
		}

		for _, outputItem := range outputItems["items"] {
			key := buildKey(outputItem.DataSubject, outputItem.DataType)
			subject, ok := subjectInventory[key]
			if !ok {
				// key not found, add a new item
				if outputItem.DataSubject == "" {
					outputItem.DataSubject = PLACEHOLDER_VALUE
				}
				ruleFailure := subjectRuleFailures[key]
				subject = types.Subject{
					DataSubject:              outputItem.DataSubject,
					DataType:                 outputItem.DataType,
					CriticalRiskFindingCount: ruleFailure.CriticalRiskFindingCount,
					HighRiskFindingCount:     ruleFailure.HighRiskFindingCount,
					MediumRiskFindingCount:   ruleFailure.MediumRiskFindingCount,
					LowRiskFindingCount:      ruleFailure.LowRiskFindingCount,
					RulesPassedCount:         localRuleCounter - len(ruleFailure.TriggeredRules),
				}
			}
			subject.DetectionCount += 1
			subjectInventory[key] = subject
		}
	}

	var thirdPartyInventory []types.ThirdParty
	for _, component := range reportData.Dataflow.Components {
		if component.SubType != "third_party" {
			continue
		}

		thirdPartyFailure, ok := thirdPartyRuleFailures[component.Name]
		if !ok {
			// no failures, therefore no associated data subjects
			thirdPartyInventory = append(thirdPartyInventory, types.ThirdParty{
				ThirdParty:               component.Name,
				DataSubject:              PLACEHOLDER_VALUE,
				DataTypes:                []string{PLACEHOLDER_VALUE},
				CriticalRiskFindingCount: 0,
				HighRiskFindingCount:     0,
				MediumRiskFindingCount:   0,
				LowRiskFindingCount:      0,
				RulesPassedCount:         0,
			})
		}

		for _, ruleFailure := range thirdPartyFailure {
			thirdPartyInventory = append(thirdPartyInventory, types.ThirdParty{
				ThirdParty:               component.Name,
				DataSubject:              ruleFailure.DataSubject,
				DataTypes:                slices.Collect(maps.Keys(ruleFailure.DataTypes)),
				CriticalRiskFindingCount: ruleFailure.CriticalRiskFindingCount,
				HighRiskFindingCount:     ruleFailure.HighRiskFindingCount,
				MediumRiskFindingCount:   ruleFailure.MediumRiskFindingCount,
				LowRiskFindingCount:      ruleFailure.LowRiskFindingCount,
				RulesPassedCount:         thirdPartyRulesCounter[component.Name].Count - len(thirdPartyRulesCounter[component.Name].SubjectFailures[ruleFailure.DataSubject]),
			})
		}
	}

	subjects := slices.Collect(maps.Values(subjectInventory))
	sortInventory(subjects, thirdPartyInventory)

	reportData.PrivacyReport = &types.Report{
		Subjects:   subjects,
		ThirdParty: thirdPartyInventory,
	}
	return nil
}

func sortInventory(subjectInventory []types.Subject, thirdPartyInventory []types.ThirdParty) {
	// sort subject
	sort.Slice(subjectInventory, func(i, j int) bool {
		if subjectInventory[i].DataSubject != subjectInventory[j].DataSubject {
			// order placeholder subjects last of the list
			if subjectInventory[i].DataSubject == PLACEHOLDER_VALUE {
				return false
			}
			if subjectInventory[j].DataSubject == PLACEHOLDER_VALUE {
				return true
			}
			return subjectInventory[i].DataSubject < subjectInventory[j].DataSubject
		}
		return subjectInventory[i].DataType < subjectInventory[j].DataType
	})

	// sort third party
	sort.Slice(thirdPartyInventory, func(i, j int) bool {
		if thirdPartyInventory[i].ThirdParty != thirdPartyInventory[j].ThirdParty {
			return thirdPartyInventory[i].ThirdParty < thirdPartyInventory[j].ThirdParty
		}
		return thirdPartyInventory[i].DataSubject < thirdPartyInventory[j].DataSubject
	})
}

func buildKey(dataSubject string, dataType string) string {
	return dataSubject + ":" + strings.ToUpper(dataType)
}
