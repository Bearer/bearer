package types

type Report struct {
	Subjects   []Subject    `json:"subjects,omitempty" yaml:"subjects"`
	ThirdParty []ThirdParty `json:"third_party,omitempty" yaml:"third_party"`
}

type ThirdParty struct {
	ThirdParty               string   `json:"third_party,omitempty" yaml:"third_party"`
	DataSubject              string   `json:"subject_name,omitempty" yaml:"subject_name"`
	DataTypes                []string `json:"data_types,omitempty" yaml:"data_types"`
	CriticalRiskFindingCount int      `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFindingCount     int      `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFindingCount   int      `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFindingCount      int      `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	RulesPassedCount         int      `json:"rules_passed_count" yaml:"rules_passed_count"`
}

type Subject struct {
	DataSubject              string `json:"subject_name,omitempty" yaml:"subject_name"`
	DataType                 string `json:"name,omitempty" yaml:"name"`
	DetectionCount           int    `json:"detection_count" yaml:"detection_count"`
	CriticalRiskFindingCount int    `json:"critical_risk_failure_count" yaml:"critical_risk_failure_count"`
	HighRiskFindingCount     int    `json:"high_risk_failure_count" yaml:"high_risk_failure_count"`
	MediumRiskFindingCount   int    `json:"medium_risk_failure_count" yaml:"medium_risk_failure_count"`
	LowRiskFindingCount      int    `json:"low_risk_failure_count" yaml:"low_risk_failure_count"`
	RulesPassedCount         int    `json:"rules_passed_count" yaml:"rules_passed_count"`
}
