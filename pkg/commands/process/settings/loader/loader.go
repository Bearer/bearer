package loader

import (
	"errors"
	"fmt"
	"slices"

	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/commands/process/settings/policies"
	"github.com/bearer/bearer/pkg/commands/process/settings/rules"
	"github.com/bearer/bearer/pkg/engine"
	"github.com/bearer/bearer/pkg/flag"
	flagtypes "github.com/bearer/bearer/pkg/flag/types"
	"github.com/bearer/bearer/pkg/util/ignore"
	"github.com/bearer/bearer/pkg/version_check"
)

func FromOptions(
	opts flagtypes.Options,
	versionMeta *version_check.VersionMeta,
	engine engine.Engine,
	foundLanguageIDs []string,
) (settings.Config, error) {
	policies, err := policies.Load()
	if err != nil {
		return settings.Config{}, fmt.Errorf("failed to load policies: %w", err)
	}

	result, err := rules.Load(
		opts.ExternalRuleDir,
		opts.RuleOptions,
		versionMeta,
		engine,
		opts.ScanOptions.Force,
		foundLanguageIDs,
	)
	if err != nil {
		return settings.Config{}, err
	}

	ignoredFingerprints, _, _, err := ignore.GetIgnoredFingerprints(opts.IgnoreFile, &opts.Target)
	if err != nil {
		return settings.Config{}, err
	}

	config := settings.Config{
		Client: opts.Client,
		Worker: settings.WorkerOptions{
			Timeout:                   settings.Timeout,
			TimeoutFileMinimum:        settings.TimeoutFileMinimum,
			TimeoutFileMaximum:        settings.TimeoutFileMaximum,
			TimeoutFileBytesPerSecond: settings.TimeoutFileBytesPerSecond,
			TimeoutWorkerOnline:       settings.TimeoutWorkerOnline,
			FileSizeMaximum:           settings.FileSizeMaximum,
			ExistingWorker:            settings.ExistingWorker,
		},
		Scan:                opts.ScanOptions,
		Report:              opts.ReportOptions,
		IgnoredFingerprints: ignoredFingerprints,
		NoColor:             opts.NoColor || opts.Output != "",
		DebugProfile:        opts.DebugProfile,
		Debug:               opts.Debug,
		LogLevel:            opts.LogLevel,
		IgnoreFile:          opts.IgnoreFile,
		IgnoreGit:           opts.IgnoreGit,
		Policies:            policies,
		Rules:               result.Rules,
		LoadedRuleCount:     result.LoadedRuleCount,
		BuiltInRules:        result.BuiltInRules,
		CacheUsed:           result.CacheUsed,
		BearerRulesVersion:  result.BearerRulesVersion,
	}

	if config.Scan.Diff {
		if !slices.Contains([]string{flag.ReportSecurity, flag.ReportSaaS}, config.Report.Report) {
			return settings.Config{}, errors.New("diff base branch is only supported for the security report")
		}
	}

	return config, nil
}
