package metricsscan

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/bearer/curio/battle_tests/scan"
)

type DataType struct {
	Name        string  `json:"name"`
	Occurrences float64 `json:"occurrences"`
}

type ScanReport struct {
	NumberOfLines     int        `json:"number_of_lines"`
	NumberOfDataTypes int        `json:"number_of_data_types"`
	DataTypes         []DataType `json:"data_types"`
}

type MetricsReport struct {
	URL        string  // full url
	RepoSizeKB float64 // repository size in KB
	Time       float64 // time it took in Seconds
	Memory     float64 // maximum memory used in MB

	// Start of statistics report
	NumberOfDataTypes  float64    // number of data types detected
	NumberOfLineOfCode float64    // number of line of code
	DataTypes          []DataType // datatypes detected
	// End of statistics report
}

func FakeScanRepository(repositoryUrl string, reportingChan chan *MetricsReport) {
	metrics := &MetricsReport{
		URL:    repositoryUrl,
		Memory: 0,
	}

	log.Printf("processing repository %s", repositoryUrl)

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

func ScanRepository(repositoryUrl string, reportingChan chan *MetricsReport) {
	metrics := &MetricsReport{
		URL:    repositoryUrl,
		Memory: 0,
	}

	log.Printf("processing repository %s", repositoryUrl)

	scanner := scan.NewScan(repositoryUrl)
	scanCrashedHandler := func(err error) {
		scanner.Cleanup()
		log.Printf("scan failed %s", err)
		reportingChan <- metrics
	}

	output, startTime, err := scanner.Start()

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

	if err != nil {
		scanCrashedHandler(fmt.Errorf("got report response which errored %w", err))
		return
	}

	scanner.Cleanup()
	reportingChan <- metrics
}
