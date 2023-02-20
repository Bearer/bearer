package settings_test

import (
	"testing"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/stretchr/testify/assert"
)

func TestRuleSeverityMapping(t *testing.T) {
	opts := flag.Options{
		ScanOptions: flag.ScanOptions{},
		RuleOptions: flag.RuleOptions{},
		RepoOptions: flag.RepoOptions{},
		ReportOptions: flag.ReportOptions{
			Report: "summary",
			Severity: map[string]bool{
				"critical": true,
				"high":     true,
				"medium":   true,
				"low":      true,
				"warning":  true,
			},
		},
		GeneralOptions: flag.GeneralOptions{},
	}
	config, err := settings.FromOptions(opts)
	if err != nil {
		t.Fatalf("failed to generate config:%s", err)
	}

	// test random rule's severity has been mapped
	// as expected
	loggerRule := config.Rules["ruby_rails_logger"]
	assert.Equal(t, loggerRule.Severity["Personal Data"] != "", true)
}
