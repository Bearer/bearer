package saas

import (
	"errors"
	"os"
	"strings"

	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/commands/process/gitrepository"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	saas "github.com/bearer/bearer/pkg/report/output/saas/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
	"github.com/bearer/bearer/pkg/report/output/types"
	"github.com/bearer/bearer/pkg/util/file"
)

func GetReport(
	reportData *types.ReportData,
	config settings.Config,
	gitContext *gitrepository.Context,
	ensureMeta bool,
) error {
	var meta *saas.Meta
	meta, err := getMeta(reportData, config, gitContext)
	if err != nil {
		if ensureMeta {
			return err
		} else {
			meta = &saas.Meta{
				Target:         config.Scan.Target,
				FoundLanguages: reportData.FoundLanguages,
			}
		}
	}

	saasFindingsBySeverity := translateFindingsBySeverity(reportData.FindingsBySeverity)
	saasIgnoredFindingsBySeverity := translateFindingsBySeverity(reportData.IgnoredFindingsBySeverity)

	reportData.SaasReport = &saas.BearerReport{
		Meta:            *meta,
		Findings:        saasFindingsBySeverity,
		IgnoredFindings: saasIgnoredFindingsBySeverity,
		DataTypes:       reportData.Dataflow.Datatypes,
		Components:      reportData.Dataflow.Components,
		Errors:          reportData.Dataflow.Errors,
		Files:           getDiscoveredFiles(config, reportData.Files),
	}

	return nil
}

func translateFindingsBySeverity[F securitytypes.GenericFinding](someFindingsBySeverity map[string][]F) map[string][]saas.SaasFinding {
	saasFindingsBySeverity := make(map[string][]saas.SaasFinding)
	for severity := range someFindingsBySeverity {
		for _, someFinding := range someFindingsBySeverity[severity] {
			finding := someFinding.GetFinding()
			saasFindingsBySeverity[severity] = append(saasFindingsBySeverity[severity], saas.SaasFinding{
				Finding:      finding,
				SeverityMeta: finding.SeverityMeta,
				IgnoreMeta:   someFinding.GetIgnoreMeta(),
			})
		}
	}
	return saasFindingsBySeverity
}

func getDiscoveredFiles(config settings.Config, files []string) []string {
	filenames := make([]string, len(files))

	for i, filename := range files {
		filenames[i] = file.GetFullFilename(config.Scan.Target, filename)
	}

	return filenames
}

func getMeta(
	reportData *types.ReportData,
	config settings.Config,
	gitContext *gitrepository.Context,
) (*saas.Meta, error) {
	if gitContext == nil {
		return nil, errors.New("not a git repository")
	}

	var messages []string
	if gitContext.Branch == "" {
		messages = append(messages,
			"Couldn't determine the name of the branch being scanned. "+
				"Please set the 'BEARER_BRANCH' environment variable.",
		)
	}
	if gitContext.DefaultBranch == "" {
		messages = append(messages,
			"Couldn't determine the default branch of the repository. "+
				"Please set the 'BEARER_DEFAULT_BRANCH' environment variable.",
		)
	}
	if gitContext.CommitHash == "" {
		messages = append(messages,
			"Couldn't determine the hash of the current commit of the repository. "+
				"Please set the 'BEARER_COMMIT' environment variable.",
		)
	}
	if gitContext.OriginURL == "" {
		messages = append(messages,
			"Couldn't determine the origin URL of the repository. "+
				"Please set the 'BEARER_REPOSITORY_URL' environment variable.",
		)
	}

	if len(messages) != 0 {
		return nil, errors.New(strings.Join(messages, "\n"))
	}

	return &saas.Meta{
		ID:                 gitContext.ID,
		Host:               gitContext.Host,
		Username:           gitContext.Owner,
		Name:               gitContext.Name,
		FullName:           gitContext.FullName,
		URL:                gitContext.OriginURL,
		Target:             config.Scan.Target,
		SHA:                gitContext.CommitHash,
		CurrentBranch:      gitContext.Branch,
		DefaultBranch:      gitContext.DefaultBranch,
		DiffBaseBranch:     gitContext.BaseBranch,
		BearerRulesVersion: config.BearerRulesVersion,
		BearerVersion:      build.Version,
		FoundLanguages:     reportData.FoundLanguages,
		GitlabPipelineId:   os.Getenv("CI_PIPELINE_ID"),
		GitlabJobId:        os.Getenv("CI_JOB_ID"),
	}, nil
}
