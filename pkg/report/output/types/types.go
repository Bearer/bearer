package types

import (
	dataflowtypes "github.com/bearer/bearer/pkg/report/output/dataflow/types"
	privacytypes "github.com/bearer/bearer/pkg/report/output/privacy/types"
	saastypes "github.com/bearer/bearer/pkg/report/output/saas/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
	statstypes "github.com/bearer/bearer/pkg/report/output/stats/types"
)

type ReportData struct {
	ReportFailed       bool
	SendToCloud        bool
	Files              []string
	Detectors          []any
	Dataflow           *DataFlow
	FindingsBySeverity map[string][]securitytypes.Finding
	PrivacyReport      *privacytypes.Report
	Stats              *statstypes.Stats
	SaasReport         *saastypes.BearerReport
}

type DataFlow struct {
	Datatypes    []dataflowtypes.Datatype     `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	Risks        []dataflowtypes.RiskDetector `json:"risks,omitempty" yaml:"risks,omitempty"`
	Components   []dataflowtypes.Component    `json:"components,omitempty" yaml:"components,omitempty"`
	Dependencies []dataflowtypes.Dependency   `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Errors       []dataflowtypes.Error        `json:"errors,omitempty" yaml:"errors,omitempty"`
}

type GenericFormatter interface {
	Format(format string) (*string, error) // TODO: ensure format is an expected format (from report flags)
}
