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
	"github.com/bearer/bearer/cmd/bearer/build"
	"github.com/bearer/bearer/pkg/commands/process/balancer/filelist"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/flag"
	"github.com/bearer/bearer/pkg/report/output/dataflow"
	"github.com/bearer/bearer/pkg/report/output/privacy"
	"github.com/bearer/bearer/pkg/report/output/security"
	"github.com/bearer/bearer/pkg/util/file"
	pointer "github.com/bearer/bearer/pkg/util/pointers"
	"github.com/gitsight/go-vcsurl"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	dataflowtypes "github.com/bearer/bearer/pkg/report/output/dataflow/types"

	"github.com/bearer/bearer/pkg/report/output/detectors"
	"github.com/bearer/bearer/pkg/report/output/stats"
	"github.com/bearer/bearer/pkg/types"
	"gopkg.in/yaml.v3"
)

var ErrUndefinedFormat = errors.New("undefined output format")

type Description struct {
	Text string `json:"text"`
}

type Configuration struct {
	Level string `json:"level"`
}

type Properties struct {
	Tags      []string `json:"tags"`
	Precision string   `json:"precision"`
}

type Help struct {
	Markdown string `json:"markdown"`
}

type SarifRule struct {
	Id                   string        `json:"id"`
	Name                 string        `json:"name"`
	Kind                 string        `json:"kind"`
	ShortDescription     Description   `json:"shortDescription"`
	FullDescription      Description   `json:"fullDescription"`
	DefaultConfiguration Configuration `json:"defaultConfiguration"`
	Properties           Properties    `json:"properties"`
	Help                 Help          `json:"help"`
}

type SarifDriver struct {
	Name       string      `json:"name"`
	SarifRules []SarifRule `json:"rules"`
}

type SarifTool struct {
	Driver SarifDriver `json:"driver"`
}

type Message struct {
	Text string `json:"text"`
}

type ArtifactLocation struct {
	URI string `json:"uri"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
}

type Region struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndColumn   int `json:"endColumn"`
	EndLine     int `json:"endLine"`
}

type Location struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
	Region           Region           `json:"region"`
}

type PartialFingerprints struct {
	PrimaryLocationLineHash               string `json:"primaryLocationLineHash,omitempty"`
	PrimaryLocationStartColumnFingerprint string `json:"primaryLocationStartColumnFingerprint,omitempty"`
}

type Result struct {
	RuleId              string               `json:"ruleId"`
	RuleIndex           int                  `json:"ruleIndex"`
	Message             Message              `json:"message"`
	Locations           []Location           `json:"locations"`
	PartialFingerprints *PartialFingerprints `json:"partialFingerprints,omitempty"`
}

type SarifRun struct {
	Tool    SarifTool `json:"tool"`
	Results []Result  `json:"results"`
}

type SarifOutput struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SarifRun `json:"runs"`
}

func ReportSarif(outputDetections *map[string][]security.Result) (*string, error) {
	// log.Error().Msgf("detections: %#v", outputDetections)

	output := SarifOutput{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []SarifRun{
			{
				Tool: SarifTool{
					Driver: SarifDriver{
						Name: "Bearer",
						SarifRules: []SarifRule{
							{
								Id:   "R01",
								Kind: "path-problem",
								Name: "js/unused-local-variable",
								ShortDescription: Description{
									Text: "Unused variable, import, function or class",
								},
								FullDescription: Description{
									Text: "Unused variables, imports, functions or classes may be a symptom of a bug and should be examined carefully.",
								},
								DefaultConfiguration: Configuration{
									Level: "note",
								},
								Help: Help{
									Markdown: "",
								},
								Properties: Properties{
									Tags:      []string{"maintainability"},
									Precision: "very-high",
								},
							},
						},
					},
				},
				Results: []Result{
					{
						RuleId:    "3f292041e51d22005ce48f39df3585d44ce1b0ad",
						RuleIndex: 0,
						Message: Message{
							Text: "Unused variable foo.",
						},
						Locations: []Location{
							{
								PhysicalLocation: PhysicalLocation{
									ArtifactLocation: ArtifactLocation{
										URI: "main.js",
									},
								},
								Region: Region{
									StartLine:   2,
									StartColumn: 7,
									EndColumn:   10,
								},
							},
						},
						// If you are uploading third-party SARIF files with the upload-action, the action will create partialFingerprints for you when they are not included in the SARIF file
						// PartialFingerprints: PartialFingerprints{
						// 	PrimaryLocationLineHash:               "123",
						// 	PrimaryLocationStartColumnFingerprint: "1234",
						// },
					},
				},
			},
		},
	}

	content, _ := ReportJSON(output)

	return content, nil
}

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

func GetOutput(report types.Report, config settings.Config) (any, *dataflow.DataFlow, error) {
	switch config.Report.Report {
	case flag.ReportDetectors:
		return detectors.GetOutput(report, config)
	case flag.ReportDataFlow:
		return GetDataflow(report, config, false)
	case flag.ReportSecurity:
		return reportSecurity(report, config)
	case flag.ReportSaaS:
		securityResults, dataflow, err := reportSecurity(report, config)
		if err != nil {
			return nil, nil, err
		}

		var meta *Meta
		meta, err = getMeta(config)
		if err != nil {
			meta = &Meta{
				Target: config.Scan.Target,
			}
		}

		files := getDiscoveredFiles(config)

		return BearerReport{
			Findings:   securityResults,
			DataTypes:  dataflow.Datatypes,
			Components: dataflow.Components,
			Files:      files,
			Meta:       *meta,
		}, nil, nil
	case flag.ReportPrivacy:
		return getPrivacyReportOutput(report, config)
	case flag.ReportStats:
		return reportStats(report, config)
	}

	return nil, nil, fmt.Errorf(`--report flag "%s" is not supported`, config.Report.Report)
}

func getDiscoveredFiles(config settings.Config) []string {
	filesDiscovered, _ := filelist.Discover(config.Scan.Target, config)
	files := []string{}
	for _, fileDiscovered := range filesDiscovered {
		files = append(files, file.GetFullFilename(config.Scan.Target, fileDiscovered.FilePath))
	}

	return files
}

func GetPrivacyReportCSVOutput(report types.Report, dataflow *dataflow.DataFlow, config settings.Config) (*string, error) {
	csvString, err := privacy.BuildCsvString(dataflow, config)
	if err != nil {
		return nil, err
	}

	content := csvString.String()

	return &content, nil
}

func getPrivacyReportOutput(report types.Report, config settings.Config) (*privacy.Report, *dataflow.DataFlow, error) {
	dataflow, _, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, nil, err
	}

	return privacy.GetOutput(dataflow, config)
}

func GetDataflow(report types.Report, config settings.Config, isInternal bool) (*dataflow.DataFlow, *dataflow.DataFlow, error) {
	reportedDetections, _, err := detectors.GetOutput(report, config)
	if err != nil {
		return nil, nil, err
	}

	for _, detection := range reportedDetections {
		detection.(map[string]interface{})["id"] = uuid.NewString()
	}

	return dataflow.GetOutput(reportedDetections, config, isInternal)
}

func reportStats(report types.Report, config settings.Config) (*stats.Stats, *dataflow.DataFlow, error) {
	dataflowOutput, _, err := GetDataflow(report, config, true)
	if err != nil {
		return nil, nil, err
	}

	return stats.GetOutput(report.Inputgocloc, dataflowOutput, config)
}

func reportSecurity(
	report types.Report,
	config settings.Config,
) (
	securityResults *security.Results,
	dataflow *dataflow.DataFlow,
	err error,
) {
	dataflow, _, err = GetDataflow(report, config, true)
	if err != nil {
		log.Debug().Msgf("error in dataflow %s", err)
		return
	}

	securityResults, err = security.GetOutput(dataflow, config)
	if err != nil {
		log.Debug().Msgf("error in security %s", err)
		return
	}

	if config.Client != nil && config.Client.Error == nil {
		var meta *Meta
		meta, err = getMeta(config)
		if err != nil {
			errorMessage := fmt.Sprintf("Unable to calculate Metadata. %s", err)
			log.Debug().Msgf(errorMessage)
			config.Client.Error = &errorMessage
			return
		}

		tmpDir, filename, err := createBearerGzipFileReport(config, meta, securityResults, dataflow)
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

	files := getDiscoveredFiles(config)

	content, _ := ReportJSON(&BearerReport{
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

func getMeta(config settings.Config) (*Meta, error) {
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

	return &Meta{
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
		BearerRulesVersion: config.BearerRulesVersion,
		BearerVersion:      build.Version,
	}, nil
}

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
	SignedID           string `json:"signed_id,omitempty" yaml:"signed_id,omitempty"`
	BearerRulesVersion string `json:"bearer_rules_version,omitempty" yaml:"bearer_rules_version,omitempty"`
	BearerVersion      string `json:"bearer_version,omitempty" yaml:"bearer_version,omitempty"`
}

type BearerReport struct {
	Meta       Meta                          `json:"meta" yaml:"meta"`
	Findings   *map[string][]security.Result `json:"findings" yaml:"findings"`
	DataTypes  []dataflowtypes.Datatype      `json:"data_types" yaml:"data_types"`
	Components []dataflowtypes.Component     `json:"components" yaml:"components"`
	Files      []string                      `json:"files" yaml:"files"`
}
