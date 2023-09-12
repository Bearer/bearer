package types

import (
	dataflowtypes "github.com/bearer/bearer/pkg/report/output/dataflow/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
	ignoretypes "github.com/bearer/bearer/pkg/util/ignore/types"
)

type Meta struct {
	ID                 string `json:"id" yaml:"id"`
	Host               string `json:"host" yaml:"host"`
	Username           string `json:"username" yaml:"username"`
	Name               string `json:"name" yaml:"name"`
	URL                string `json:"url" yaml:"url"`
	FullName           string `json:"full_name" yaml:"full_name"`
	Target             string `json:"target" yaml:"target"`
	SHA                string `json:"sha" yaml:"sha"`
	CurrentBranch      string `json:"current_branch" yaml:"current_branch"`
	DefaultBranch      string `json:"default_branch" yaml:"default_branch"`
	DiffBaseBranch     string `json:"diff_base_branch,omitempty" yaml:"diff_base_branch,omitempty"`
	SignedID           string `json:"signed_id,omitempty" yaml:"signed_id,omitempty"`
	BearerRulesVersion string `json:"bearer_rules_version,omitempty" yaml:"bearer_rules_version,omitempty"`
	BearerVersion      string `json:"bearer_version,omitempty" yaml:"bearer_version,omitempty"`
}

type BearerReport struct {
	Meta            Meta                      `json:"meta" yaml:"meta"`
	Findings        map[string][]SaasFinding  `json:"findings" yaml:"findings"`
	IgnoredFindings map[string][]SaasFinding  `json:"ignored_findings" yaml:"ignored_findings"`
	DataTypes       []dataflowtypes.Datatype  `json:"data_types" yaml:"data_types"`
	Components      []dataflowtypes.Component `json:"components" yaml:"components"`
	Errors          []dataflowtypes.Error     `json:"errors" yaml:"errors"`
	Files           []string                  `json:"files" yaml:"files"`
	// Dependencies []dataflowtypes.Dependency    `json:"dependencies" yaml:"dependencies"`
}

type SaasFinding struct {
	securitytypes.Finding
	SeverityMeta securitytypes.SeverityMeta      `json:"severity_meta" yaml:"severity_meta"`
	IgnoreMeta   *ignoretypes.IgnoredFingerprint `json:"ignore_meta,omitempty" yaml:"ignore_meta,omitempty"`
}
