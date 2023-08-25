package saas

import (
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gitsight/go-vcsurl"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"

	"github.com/bearer/bearer/api"
	"github.com/bearer/bearer/api/s3"
	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	saas "github.com/bearer/bearer/pkg/report/output/saas/types"
	securitytypes "github.com/bearer/bearer/pkg/report/output/security/types"
	"github.com/bearer/bearer/pkg/report/output/types"
	"github.com/bearer/bearer/pkg/util/file"
	util "github.com/bearer/bearer/pkg/util/output"
	pointer "github.com/bearer/bearer/pkg/util/pointers"
)

func GetReport(reportData *types.ReportData, config settings.Config, ensureMeta bool) error {
	var meta *saas.Meta
	meta, err := getMeta(config)
	if err != nil {
		if ensureMeta {
			return err
		} else {
			meta = &saas.Meta{
				Target: config.Scan.Target,
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

func GetVCSInfo(config settings.Config) (*vcsurl.VCS, error) {
	gitRemote, err := getRemote(config.Scan.Target)
	if err != nil {
		return nil, err
	}

	info, err := vcsurl.Parse(*gitRemote)
	if err != nil {
		log.Debug().Msgf("couldn't parse origin url %s", err)
		return nil, err
	}

	return info, nil
}

func SendReport(config settings.Config, reportData *types.ReportData) {
	if reportData.SaasReport == nil {
		err := GetReport(reportData, config, true)
		if err != nil {
			errorMessage := fmt.Sprintf("Unable to calculate Metadata. %s", err)
			log.Debug().Msgf(errorMessage)
			config.Client.Error = &errorMessage
		}
	}

	tmpDir, filename, err := createBearerGzipFileReport(config, reportData)
	if err != nil {
		config.Client.Error = pointer.String("Could not compress report.")
		log.Debug().Msgf("error creating report %s", err)
	}

	defer os.RemoveAll(*tmpDir)

	err = sendReportToBearer(config.Client, &reportData.SaasReport.Meta, filename)
	if err != nil {
		config.Client.Error = pointer.String("Report upload failed.")
		log.Debug().Msgf("error sending report to Bearer cloud: %s", err)
	}
}

func translateFindingsBySeverity(findingBySeverity map[string][]securitytypes.Finding) map[string][]saas.SaasFinding {
	saasFindingsBySeverity := make(map[string][]saas.SaasFinding)
	for _, severity := range maps.Keys(findingBySeverity) {
		for _, finding := range findingBySeverity[severity] {
			saasFindingsBySeverity[severity] = append(saasFindingsBySeverity[severity], saas.SaasFinding{Finding: finding, SeverityMeta: finding.SeverityMeta})
		}
	}
	return saasFindingsBySeverity
}

func sendReportToBearer(client *api.API, meta *saas.Meta, filename *string) error {
	fileUploadOffer, err := s3.UploadS3(&s3.UploadRequestS3{
		Api:             client,
		FilePath:        *filename,
		FilePrefix:      "bearer_security_report",
		ContentType:     "application/json",
		ContentEncoding: "gzip",
	})
	if err != nil {
		return err
	}

	meta.SignedID = fileUploadOffer.SignedID

	err = client.ScanFinished(meta)
	if err != nil {
		return err
	}

	return nil
}

func getDiscoveredFiles(config settings.Config, files []string) []string {
	filenames := make([]string, len(files))

	for i, filename := range files {
		filenames[i] = file.GetFullFilename(config.Scan.Target, filename)
	}

	return filenames
}

func createBearerGzipFileReport(
	config settings.Config,
	reportData *types.ReportData,
) (*string, *string, error) {
	tempDir, err := os.MkdirTemp("", "reports")
	if err != nil {
		return nil, nil, err
	}

	file, err := os.CreateTemp(tempDir, "security-*.json.gz")
	if err != nil {
		return &tempDir, nil, err
	}

	content, _ := util.ReportJSON(reportData.SaasReport)
	gzWriter := gzip.NewWriter(file)
	_, err = gzWriter.Write([]byte(*content))
	if err != nil {
		return nil, nil, err
	}
	gzWriter.Close()

	filename := file.Name()

	return &tempDir, &filename, nil
}

func getMeta(config settings.Config) (*saas.Meta, error) {
	sha, err := getSha(config.Scan.Target)
	if err != nil {
		return nil, err
	}

	currentBranch, err := getCurrentBranch(config.Scan.Target)
	if err != nil {
		return nil, err
	}

	defaultBranch, err := getDefaultBranch(config.Scan.Target)
	if err != nil {
		return nil, err
	}

	info, err := GetVCSInfo(config)
	if err != nil {
		return nil, err
	}

	return &saas.Meta{
		ID:                 info.ID,
		Host:               string(info.Host),
		Username:           info.Username,
		Name:               info.Name,
		FullName:           info.FullName,
		URL:                info.Raw,
		Target:             config.Scan.Target,
		SHA:                *sha,
		CurrentBranch:      *currentBranch,
		DefaultBranch:      *defaultBranch,
		DiffBaseBranch:     config.Scan.DiffBaseBranch,
		BearerRulesVersion: config.BearerRulesVersion,
		BearerVersion:      build.Version,
	}, nil
}

func getSha(target string) (*string, error) {
	env := os.Getenv("SHA")
	if env != "" {
		return pointer.String(env), nil
	}
	bytes, err := exec.Command("git", "-C", target, "rev-parse", "HEAD").Output()
	if err != nil {
		log.Error().Msg("Couldn't extract git info for commit sha please set 'SHA' environment variable.")
		return nil, err
	}
	return pointer.String(strings.TrimSuffix(string(bytes), "\n")), nil
}

func getCurrentBranch(target string) (*string, error) {
	env := os.Getenv("CURRENT_BRANCH")
	if env != "" {
		return pointer.String(env), nil
	}
	bytes, err := exec.Command("git", "-C", target, "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		log.Error().Msg("Couldn't extract git info for current branch please set 'CURRENT_BRANCH' environment variable.")
		return nil, err
	}
	return pointer.String(strings.TrimSuffix(string(bytes), "\n")), nil
}

func getDefaultBranch(target string) (*string, error) {
	env := os.Getenv("DEFAULT_BRANCH")
	if env != "" {
		return pointer.String(env), nil
	}
	bytes, err := exec.Command("git", "-C", target, "rev-parse", "--abbrev-ref", "origin/HEAD").Output()
	if err != nil {
		log.Error().Msg("Couldn't extract the default branch of this repository. Please set 'DEFAULT_BRANCH' environment variable.")
		return nil, err
	}
	return pointer.String(strings.TrimPrefix(strings.TrimSuffix(string(bytes), "\n"), "origin/")), nil
}

func getRemote(target string) (*string, error) {
	env := os.Getenv("ORIGIN_URL")
	if env != "" {
		return pointer.String(env), nil
	}
	bytes, err := exec.Command("git", "-C", target, "remote", "get-url", "origin").Output()
	if err != nil {
		log.Error().Msg("Couldn't extract git info for origin url please set 'ORIGIN_URL' environment variable.")
		return nil, err
	}
	return pointer.String(strings.TrimSuffix(string(bytes), "\n")), nil
}
