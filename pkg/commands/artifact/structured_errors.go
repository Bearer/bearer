package artifact

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/bearer/bearer/pkg/classification"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/scanner/languagescanner"
	"github.com/bearer/bearer/pkg/scanner/ruleset"
	"github.com/bearer/bearer/pkg/structurederrors"
	outputhandler "github.com/bearer/bearer/pkg/util/output"
)

// checkRules sets up the scanner for each language in the same way as the workers, so that a rule
// the engine cannot compile is recorded against its rule and pattern. Without this, the error only
// surfaces as a message when the workers start.
func checkRules(engine engine.Engine, config settings.Config) error {
	if !config.StructuredErrors.Enabled() || !slices.Contains(config.Scan.Scanner, flag.ScannerSAST) {
		return nil
	}

	classifier, err := classification.NewClassifier(&classification.Config{Config: config})
	if err != nil {
		return err
	}

	for _, language := range engine.GetLanguages() {
		languageScanner, err := languagescanner.New(language, classifier.Schema, config.Rules)
		if err != nil {
			config.StructuredErrors.Add(engineErrorItem(language.ID(), config.Rules, err))
			return fmt.Errorf("error creating %s language scanner: %w", language.ID(), err)
		}

		languageScanner.Close()
	}

	return nil
}

func engineErrorItem(languageID string, rules map[string]*settings.Rule, err error) structurederrors.Item {
	var patternError *ruleset.PatternError
	if !errors.As(err, &patternError) {
		return structurederrors.Item{
			Severity: structurederrors.SeverityError,
			Code:     structurederrors.CodeRuleSetInvalid,
			Message:  err.Error(),
			Language: languageID,
		}
	}

	code := structurederrors.CodePatternInvalid
	if patternError.IsFilter {
		code = structurederrors.CodeFilterInvalid
	}

	item := structurederrors.Item{
		Severity:     structurederrors.SeverityError,
		Code:         code,
		Message:      patternError.Error(),
		RuleID:       patternError.RuleID,
		PatternIndex: structurederrors.Index(patternError.PatternIndex),
		Language:     languageID,
	}

	if rule := rules[patternError.RuleID]; rule != nil && patternError.PatternIndex < len(rule.Patterns) {
		item.Pattern = rule.Patterns[patternError.PatternIndex].Pattern
	}

	return item
}

// writeStructuredErrorsReport writes a report holding only the structured errors, for when rule
// errors stop the scan before there is a normal report to put them in
func writeStructuredErrorsReport(errorCollector *structurederrors.Collector, outputPath string) error {
	if !errorCollector.Enabled() {
		return nil
	}

	report, err := security.StructuredErrorsOutput(errorCollector)
	if err != nil {
		return err
	}

	// written in the same way as a normal report
	logger := outputhandler.StdOutLog
	if outputPath != "" {
		reportFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("error creating output file %w", err)
		}
		defer func() { _ = reportFile.Close() }()

		logger = outputhandler.PlainLogger(reportFile)
	}

	logger(report)
	return nil
}
