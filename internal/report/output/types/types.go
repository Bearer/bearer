package types

import (
	dataflowtypes "github.com/bearer/bearer/internal/report/output/dataflow/types"
	privacytypes "github.com/bearer/bearer/internal/report/output/privacy/types"
	saastypes "github.com/bearer/bearer/internal/report/output/saas/types"
	securitytypes "github.com/bearer/bearer/internal/report/output/security/types"
	statstypes "github.com/bearer/bearer/internal/report/output/stats/types"
)

type ReportData struct {
	ReportFailed              bool
	Files                     []string
	FoundLanguages            map[string]int32 // language => loc e.g. { "Ruby": 6742, "JavaScript": 122 }
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
	Datatypes          []dataflowtypes.Datatype     `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	ExpectedDetections []dataflowtypes.RiskDetector `json:"expected_detections,omitempty" yaml:"expected_detections,omitempty"`
	Risks              []dataflowtypes.RiskDetector `json:"risks,omitempty" yaml:"risks,omitempty"`
	Components         []dataflowtypes.Component    `json:"components,omitempty" yaml:"components,omitempty"`
	Dependencies       []dataflowtypes.Dependency   `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Errors             []dataflowtypes.Error        `json:"errors,omitempty" yaml:"errors,omitempty"`
}

type GenericFormatter interface {
	Format(format string) (string, error) // TODO: ensure format is an expected format (from report flags)
}
