package structurederrors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bearer/bearer/pkg/structurederrors"
)

func TestReport(t *testing.T) {
	collector := structurederrors.New()
	collector.RegisterRule("aux", "rules/b.yml", "parent")
	collector.Add(structurederrors.Item{
		Severity:     structurederrors.SeverityError,
		Code:         structurederrors.CodePatternInvalid,
		RuleID:       "aux",
		PatternIndex: structurederrors.Index(1),
	})
	collector.Add(structurederrors.Item{
		Severity: structurederrors.SeverityError,
		Code:     structurederrors.CodeRuleYAMLInvalid,
		RuleFile: "rules/a.yml",
	})

	report := collector.Report()
	require.NotNil(t, report)
	assert.Equal(t, structurederrors.SchemaVersion, report.SchemaVersion)
	require.Len(t, report.Items, 2)

	assert.Equal(t, structurederrors.CodeRuleYAMLInvalid, report.Items[0].Code)
	assert.Equal(t, "rules/b.yml", report.Items[1].RuleFile)
	assert.Equal(t, "parent", report.Items[1].ParentRuleID)
}

func TestReportEmpty(t *testing.T) {
	report := structurederrors.New().Report()
	require.NotNil(t, report)
	assert.Equal(t, []structurederrors.Item{}, report.Items)
}

func TestDisabled(t *testing.T) {
	var collector *structurederrors.Collector

	collector.RegisterRule("rule", "rule.yml", "")
	collector.Add(structurederrors.Item{Code: structurederrors.CodePatternInvalid})

	assert.False(t, collector.Enabled())
	assert.Nil(t, collector.Report())
}
