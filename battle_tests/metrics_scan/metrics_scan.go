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
	PolicyName        string   `json:"policy_name" yaml:"policy_name"`
	PolicyDisplayId   string   `json:"policy_display_id" yaml:"policy_display_id"`
	// PolicyDescription string   `json:"policy_description" yaml:"policy_description"`
	LineNumber        int      `json:"line_number,omitempty" yaml:"line_number,omitempty"`
	Filename          string   `json:"filename,omitempty" yaml:"filename,omitempty"`
	CategoryGroups    []string `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	// ParentLineNumber  int      `json:"parent_line_number,omitempty" yaml:"parent_line_number,omitempty"`
	// ParentContent     string   `json:"parent_content,omitempty" yaml:"parent_content,omitempty"`
	// OmitParent        bool     `json:"omit_parent" yaml:"omit_parent"`
	// DetailedContext   string   `json:"detailed_context,omitempty" yaml:"detailed_context,omitempty"`
}

type MetricsReport struct {
	URL        string  // full url
	RepoSizeKB float64 // repository size in KB
	Time       float64 // time it took in Seconds
	Language   string
	Memory     float64 // maximum memory used in MB

	// Start of statistics report
	NumberOfDataTypes  float64    // number of data types detected
	NumberOfLineOfCode float64    // number of line of code
	DataTypes          []DataType // datatypes detected
	// End of statistics report

	PolicyBreaches map[string]PolicyScanReport // policy breaches, grouped by severity
}

func FakeScanRepository(repositoryUrl string, reportingChan chan *MetricsReport) {
	metrics := &MetricsReport{
		URL:    repositoryUrl,
		Memory: 0,
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
		Memory:   0,
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
		return
	}

	metrics.RepoSizeKB = math.Round(float64(scanner.FSSize) / 1024)

	endTime := time.Now()
	timeDuration := endTime.Sub(*startTime).Seconds()
	metrics.Time = float64(timeDuration)

	var reportData ScanReport
	err = json.Unmarshal(output, &reportData)

	if err != nil {
		return
	}

	metrics.NumberOfDataTypes = float64(reportData.NumberOfDataTypes)
	metrics.NumberOfLineOfCode = float64(reportData.NumberOfLines)
	metrics.DataTypes = reportData.DataTypes

	// @todo FIXME: Reachable?
	if err != nil {
		scanCrashedHandler(fmt.Errorf("got report response which errored %w", err))
		return
	}

	if language == "ruby" {
		// Run a policies scan
		policiesOutput, _, err := scanner.Start("policies")
		if err != nil {
			return
		}

		var policiesReportData map[string]PolicyScanReport
		err = json.Unmarshal(policiesOutput, &policiesReportData)
		if err != nil {
			return
		}

		metrics.PolicyBreaches = policiesReportData
	}

	scanner.Cleanup()
	reportingChan <- metrics
}
