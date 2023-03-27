package output

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bearer/bearer/api"
	"github.com/bearer/bearer/api/s3"
	"github.com/bearer/bearer/pkg/commands/process/balancer/filelist"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	dataflowoutput "github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/gitsight/go-vcsurl"
	"github.com/google/uuid"

	"github.com/bearer/bearer/pkg/report/output/detectors"
	"github.com/bearer/bearer/pkg/report/output/stats"
	"github.com/bearer/bearer/pkg/types"
	"gopkg.in/yaml.v3"

	"github.com/hhatto/gocloc"
	"github.com/rs/zerolog/log"
)

var ErrUndefinedFormat = errors.New("undefined output format")

func ReportJSON(outputDetections any) (*string, error) {
	jsonBytes, err := json.Marshal(&outputDetections)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal detections: %s", err)
	}

	content := string(jsonBytes)
	return &content, nil
}

func ReportYAML(outputDetections any) (*string, error) {
	yamlBytes, err := yaml.Marshal(&outputDetections)
	if err != nil {
		return nil, fmt.Errorf("failed to yaml marshal detections: %s", err)
	}

	content := string(yamlBytes)
	return &content, nil
}

func GetOutput(report types.Report, config settings.Config) (any, *gocloc.Result, *dataflow.DataFlow, error) {
	switch config.Report.Report {
	case flag.ReportDetectors:
		return detectors.GetOutput(report, config)
	case flag.ReportDataFlow:
		return GetDataflow(report, config, false)
	case flag.ReportSecurity:
		return reportSecurity(report, config)
	case flag.ReportPrivacy:
		return getPrivacyReportOutput(report, config)
	case flag.ReportStats:
		return reportStats(report, config)
	}

	return nil, nil, nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
}

func GetPrivacyReportCSVOutput(report types.Report, lineOfCodeOutput *gocloc.Result, dataflow *dataflow.DataFlow, config settings.Config) (*string, error) {
	csvString, err := privacy.BuildCsvString(dataflow, lineOfCodeOutput, config)
	if err != nil {
		return nil, err
	}

	content := csvString.String()

	return &content, nil
}

func getPrivacyReportOutput(report types.Report, config settings.Config) (*privacy.Report, *gocloc.Result, *dataflow.DataFlow, error) {
	lineOfCodeOutput, err := stats.GoclocDetectorOutput(config.Scan.Target)
	if err != nil {
		log.Error().Msgf("error in line of code output %s", err)
		return nil, nil, nil, err
	}

	dataflow, _, _, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, nil, nil, err
	}

	return privacy.GetOutput(dataflow, lineOfCodeOutput, config)
}

func GetDataflow(report types.Report, config settings.Config, isInternal bool) (*dataflow.DataFlow, *gocloc.Result, *dataflow.DataFlow, error) {
	reportedDetections, _, _, err := detectors.GetOutput(report, config)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, detection := range reportedDetections {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}

	return dataflow.GetOutput(reportedDetections, config, isInternal)
}

func reportStats(report types.Report, config settings.Config) (*stats.Stats, *gocloc.Result, *dataflow.DataFlow, error) {
	lineOfCodeOutput, err := stats.GoclocDetectorOutput(config.Scan.Target)
	if err != nil {
		return nil, nil, nil, err
	}

	dataflowOutput, _, _, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, nil, nil, err
	}

	return stats.GetOutput(lineOfCodeOutput, dataflowOutput, config)
}

func reportSecurity(
	report types.Report,
	config settings.Config,
) (
	securityResults *security.Results,
	lineOfCodeOutput *gocloc.Result,
	dataflow *dataflow.DataFlow,
	err error,
) {
	lineOfCodeOutput, err = stats.GoclocDetectorOutput(config.Scan.Target)
	if err != nil {
		log.Error().Msgf("error in line of code output %s", err)
		return
	}

	dataflow, _, _, err = GetDataflow(report, config, true)
	if err != nil {
		log.Error().Msgf("error in dataflow %s", err)
		return
	}

	securityResults, err = security.GetOutput(dataflow, config)
	if err != nil {
		log.Error().Msgf("error in security %s", err)
		return
	}

	if config.Client != nil {
		meta, err := getMeta(config)
		if err != nil {
			log.Error().Msgf("couldn't get meta for repo %s", err)
		}

		tmpDir, filename, err := createBearerGzipFileReport(config, meta, securityResults, dataflow)
		if err != nil {
			log.Error().Msgf("error creating report %s", err)
		}

		defer os.RemoveAll(*tmpDir)

		err = sendReportToBearer(config.Client, meta, filename)

		if err != nil {
			log.Error().Msgf("error sending report to Bearer")
		}
	}

	return
}

func sendReportToBearer(client *api.API, meta *Meta, filename *string) error {
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

func createBearerGzipFileReport(
	config settings.Config,
	meta *Meta,
	securityResults *security.Results,
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

	filesDiscovered, _ := filelist.Discover(config.Scan.Target, config)
	files := []string{}
	for _, fileDiscovered := range filesDiscovered {
		files = append(files, dataflowoutput.GetFullFilename(config.Scan.Target, fileDiscovered.FilePath))
	}

	content, _ := ReportJSON(&BearerReport{
		Findings:   securityResults,
		DataTypes:  &dataflow.Datatypes,
		Components: &dataflow.Components,
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

func getMeta(config settings.Config) (*Meta, error) {
	output, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		log.Error().Msgf("couldn't get git info %s", err)
		return nil, err
	}

	info, err := vcsurl.Parse(strings.TrimSuffix(string(output), "\n"))
	if err != nil {
		log.Error().Msgf("couldn't parse url %s", err)
		return nil, err
	}

	return &Meta{
		ID:       info.ID,
		Host:     string(info.Host),
		Username: info.Username,
		Name:     info.Name,
		FullName: info.FullName,
		URL:      info.Raw,
		Target:   config.Scan.Target,
	}, nil
}

type Meta struct {
	ID       string `json:"id" yaml:"id"`
	Host     string `json:"host" yaml:"host"`
	Username string `json:"username" yaml:"username"`
	Name     string `json:"name" yaml:"name"`
	URL      string `json:"url" yaml:"url"`
	FullName string `json:"full_name" yaml:"full_name"`
	Target   string `json:"target" yaml:"target"`
	SignedID string `json:"signed_id,omitempty" yaml:"signed_id,omitempty"`
}

type BearerReport struct {
	Meta       Meta                          `json:"meta" yaml:"meta"`
	Findings   *map[string][]security.Result `json:"findings" yaml:"findings"`
	DataTypes  *dataflow.DataTypes           `json:"data_types" yaml:"data_types"`
	Components *dataflow.Components          `json:"components" yaml:"components"`
	Files      []string                      `json:"files" yaml:"files"`
}
