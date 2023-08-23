package symfony_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"

	"github.com/bearer/bearer/internal/detectors/internal/testhelper"

	reportdetectors "github.com/bearer/bearer/internal/report/detectors"
)

const detectorType = reportdetectors.DetectorSymfony

var registrations = testhelper.RegistrationFor(detectorType)

func TestBuildReportFramework(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{`php_not_symfony`, false},
		{`not_php`, false},
		{`symfony`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := testhelper.Extract(
				t,
				filepath.Join("testdata", tt.name),
				registrations,
				detectorType,
			)

			if tt.expected {
				assert.Greater(t, len(report.Frameworks), 0)
			} else {
				assert.Len(t, report.Frameworks, 0)
			}
		})
	}
}

func TestBuildReportDataStores(t *testing.T) {
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "symfony"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Frameworks)
}
