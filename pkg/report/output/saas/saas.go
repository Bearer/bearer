package saas

import (
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bearer/bearer/api"
	"github.com/bearer/bearer/api/s3"
	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/commands/process/orchestrator/filelist"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	saas "github.com/bearer/bearer/pkg/report/output/saas/types"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/util/file"
	util "github.com/bearer/bearer/pkg/util/output"
	pointer "github.com/bearer/bearer/pkg/util/pointers"
	"github.com/gitsight/go-vcsurl"
	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"
)

func GetReport(
	config settings.Config,
	securityResults *map[string][]security.Result,
	dataflow *dataflow.DataFlow,
	goclocResult *gocloc.Result,
) (saas.BearerReport, *dataflow.DataFlow, error) {
	var meta *saas.Meta
	meta, err := getMeta(config)
	if err != nil {
		meta = &saas.Meta{
			Target: config.Scan.Target,
		}
	}

	files := getDiscoveredFiles(config, goclocResult)

	return saas.BearerReport{
		Findings:   securityResults,
		DataTypes:  dataflow.Datatypes,
		Components: dataflow.Components,
		Errors:     dataflow.Errors,
		Files:      files,
		Meta:       *meta,
		// Dependencies: dataflow.Dependencies,
	}, dataflow, nil
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

	gitRemote, err := getRemote(config.Scan.Target)
	if err != nil {
		return nil, err
	}

	info, err := vcsurl.Parse(*gitRemote)
	if err != nil {
		log.Debug().Msgf("couldn't parse origin url %s", err)
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

func SendReport(
	config settings.Config,
	securityResults *map[string][]security.Result,
	goclocResult *gocloc.Result,
	dataflow *dataflow.DataFlow,
) {
	var meta *saas.Meta
	meta, err := getMeta(config)
	if err != nil {
		errorMessage := fmt.Sprintf("Unable to calculate Metadata. %s", err)
		log.Debug().Msgf(errorMessage)
		config.Client.Error = &errorMessage
		return
	}

	tmpDir, filename, err := createBearerGzipFileReport(config, meta, securityResults, goclocResult, dataflow)
	if err != nil {
		config.Client.Error = pointer.String("Could not compress report.")
		log.Debug().Msgf("error creating report %s", err)
	}

	defer os.RemoveAll(*tmpDir)

	err = sendReportToBearer(config.Client, meta, filename)
	if err != nil {
		config.Client.Error = pointer.String("Report upload failed.")
		log.Debug().Msgf("error sending report to Bearer cloud: %s", err)
	}
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

func getDiscoveredFiles(config settings.Config, goclocResult *gocloc.Result) []string {
	filesDiscovered, _ := filelist.Discover(config.Scan.Target, goclocResult, config)
	files := []string{}
	for _, fileDiscovered := range filesDiscovered {
		files = append(files, file.GetFullFilename(config.Scan.Target, fileDiscovered.FilePath))
	}

	return files
}

func createBearerGzipFileReport(
	config settings.Config,
	meta *saas.Meta,
	securityResults *security.Results,
	goclocResult *gocloc.Result,
	dataflow *dataflow.DataFlow,
) (*string, *string, error) {
	tempDir, err := os.MkdirTemp("", "reports")
	if err != nil {
		return nil, nil, err
	}

	file, err := os.CreateTemp(tempDir, "security-*.json.gz")
	if err != nil {
		return &tempDir, nil, err
	}

	files := getDiscoveredFiles(config, goclocResult)

	content, _ := util.ReportJSON(&saas.BearerReport{
		Findings:   securityResults,
		DataTypes:  dataflow.Datatypes,
		Components: dataflow.Components,
		Files:      files,
		Meta:       *meta,
	})

	gzWriter := gzip.NewWriter(file)
	_, err = gzWriter.Write([]byte(*content))
	if err != nil {
		return nil, nil, err
	}
	gzWriter.Close()

	filename := file.Name()

	return &tempDir, &filename, nil
}

func getSha(target string) (*string, error) {
	env := os.Getenv("SHA")
	if env != "" {
		return pointer.String(env), nil
	}
	bytes, err := exec.Command("git", "-C", target, "rev-parse", "HEAD").Output()
	if err != nil {
		log.Debug().Msg("Couldn't extract git info for commit sha please set 'SHA' environment variable.")
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
		log.Debug().Msg("Couldn't extract git info for current branch please set 'CURRENT_BRANCH' environment variable.")
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
		log.Debug().Msg("Couldn't extract the default branch of this repository please set 'DEFAULT_BRANCH' environment variable.")
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
		log.Debug().Msg("Couldn't extract git info for origin url please set 'ORIGIN_URL' environment variable.")
		return nil, err
	}
	return pointer.String(strings.TrimSuffix(string(bytes), "\n")), nil
}
