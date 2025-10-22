package types

import (
	dataflowtypes "github.com/bearer/bearer/pkg/report/output/dataflow/types"
	privacytypes "github.com/bearer/bearer/pkg/report/output/privacy/types"
	saastypes "github.com/bearer/bearer/pkg/report/output/saas/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
	statstypes "github.com/bearer/bearer/pkg/report/output/stats/types"
)

type ReportData struct {
	ReportFailed              bool
	Files                     []string
	FoundLanguages            map[string]int32 // language => loc e.g. { "Ruby": 6742, "JavaScript": 122 }
	LanguageFiles             map[string]int32 // language => file count
	TotalLanguageFiles        int32
	Detectors                 []any
	Dataflow                  *DataFlow
	RawFindings               []securitytypes.RawFinding `json:"findings"`
	FindingsBySeverity        map[string][]securitytypes.Finding
	IgnoredFindingsBySeverity map[string][]securitytypes.IgnoredFinding
	PrivacyReport             *privacytypes.Report
	Stats                     *statstypes.Stats
	SaasReport                *saastypes.BearerReport
	ExpectedDetections        []securitytypes.ExpectedDetection
}

type DataFlow struct {
	Languages          []LanguageStats              `json:"languages,omitempty" yaml:"languages,omitempty"`
	Datatypes          []dataflowtypes.Datatype     `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	ExpectedDetections []dataflowtypes.RiskDetector `json:"expected_detections,omitempty" yaml:"expected_detections,omitempty"`
	Risks              []dataflowtypes.RiskDetector `json:"risks,omitempty" yaml:"risks,omitempty"`
	Components         []dataflowtypes.Component    `json:"components,omitempty" yaml:"components,omitempty"`
	Dependencies       []dataflowtypes.Dependency   `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Errors             []dataflowtypes.Error        `json:"errors,omitempty" yaml:"errors,omitempty"`
	Paths              []dataflowtypes.Path         `json:"paths,omitempty" yaml:"paths,omitempty"`
	TotalFiles         int32                        `json:"total_files" yaml:"total_files"`
}

type LanguageStats struct {
	Language string `json:"language" yaml:"language"`
	Lines    int32  `json:"lines" yaml:"lines"`
	Files    int32  `json:"files" yaml:"files"`
}

type GenericFormatter interface {
	Format(format string) (string, error) // TODO: ensure format is an expected format (from report flags)
}
