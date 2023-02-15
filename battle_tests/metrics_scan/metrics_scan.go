package metricsscan

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/bearer/curio/battle_tests/scan"
	"github.com/rs/zerolog/log"
)

type DataType struct {
	Name        string  `json:"name" yaml:"name"`
	Occurrences float64 `json:"occurrences" yaml:"occurrences"`
}

type ScanReport struct {
	NumberOfLines     int        `json:"number_of_lines" yaml:"number_of_lines"`
	NumberOfDataTypes int        `json:"number_of_data_types" yaml:"number_of_data_types"`
	DataTypes         []DataType `json:"data_types" yaml:"data_types"`
}

type PolicyScanReport struct {
	// PolicyName      string `json:"policy_name" yaml:"policy_name"`
	PolicyDsrid     string `json:"rule_dsrid" yaml:"rule_dsrid"`
	PolicyDisplayId string `json:"rule_display_id" yaml:"rule_display_id"`
	// PolicyDescription string   `json:"policy_description" yaml:"policy_description"`
	LineNumber     int      `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename       string   `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroups []string `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	// ParentLineNumber  int      `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	// ParentContent     string   `json:"parent_content,omitempty" yaml:"parent_content,omitempty"`
	// OmitParent        bool     `json:"omit_parent" yaml:"omit_parent"`
	// DetailedContext   string   `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
}

type MetricsReport struct {
	URL                string                        `json:"url"`          // full url
	RepoSizeKB         float64                       `json:"repo_size_kb"` // repository size in KB
	Time               float64                       `json:"time"`         // time it took in Seconds
	Language           string                        `json:"language"`
	NumberOfDataTypes  float64                       `json:"nb_data_type"`
	NumberOfLineOfCode float64                       `json:"nb_line_of_code"` // number of line of code
	DataTypes          []DataType                    `json:"data_types"`      // datatypes detected
	PolicyBreaches     map[string][]PolicyScanReport `json:"policy_breaches"` // policy breaches, grouped by severity
}

func FakeScanRepository(repositoryUrl string, reportingChan chan *MetricsReport) {
	metrics := &MetricsReport{
		URL: repositoryUrl,
	}

	log.Debug().Msgf("processing fake repository %s", repositoryUrl)

	reportData := &ScanReport{
		NumberOfLines:     12,
		NumberOfDataTypes: 43,
		DataTypes: []DataType{
			{
				Name:        "Emails",
				Occurrences: 43,
			},
		},
	}

	metrics.NumberOfDataTypes = float64(reportData.NumberOfDataTypes)
	metrics.NumberOfLineOfCode = float64(reportData.NumberOfLines)
	metrics.DataTypes = reportData.DataTypes

	reportingChan <- metrics
}

func ScanRepository(repositoryUrl string, language string, reportingChan chan *MetricsReport) {
	metrics := &MetricsReport{
		URL:      repositoryUrl,
		Language: language,
	}

	log.Debug().Msgf("processing repository %s", repositoryUrl)

	scanner := scan.NewScan(repositoryUrl)
	scanCrashedHandler := func(err error) {
		scanner.Cleanup()
		log.Printf("scan failed %s", err)
		reportingChan <- metrics
	}

	output, startTime, err := scanner.Start("stats")

	if err != nil {
		log.Debug().Msgf("stats output failed: %s", err)
		return
	}

	metrics.RepoSizeKB = math.Round(float64(scanner.FSSize) / 1024)

	endTime := time.Now()
	timeDuration := endTime.Sub(*startTime).Seconds()
	metrics.Time = float64(timeDuration)

	var reportData ScanReport
	err = json.Unmarshal(output, &reportData)

	if err != nil {
		log.Debug().Msgf("stats marshalling failed: %s", err)
		log.Debug().Msg(string(output[:]))
		scanCrashedHandler(fmt.Errorf("got report response which errored %w", err))
		return
	}

	metrics.NumberOfDataTypes = float64(reportData.NumberOfDataTypes)
	metrics.NumberOfLineOfCode = float64(reportData.NumberOfLines)
	metrics.DataTypes = reportData.DataTypes


	// Run summary
	policiesOutput, _, err := scanner.Start("summary")
	if err != nil {
		log.Debug().Msgf("summary output failed: %s", err)
		return
	}

	var policiesReportData map[string][]PolicyScanReport
	err = json.Unmarshal(policiesOutput, &policiesReportData)
	if err != nil {
		log.Debug().Msgf("summary marshalling failed: %s", err)
		log.Debug().Msg(string(policiesOutput[:]))
		scanCrashedHandler(fmt.Errorf("got report response which errored %w", err))
		return
	}

	metrics.PolicyBreaches = policiesReportData

	scanner.Cleanup()
	reportingChan <- metrics
}
