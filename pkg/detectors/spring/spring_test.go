package spring_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"

	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	reportdetectors "github.com/bearer/curio/pkg/report/detectors"
)

const detectorType = reportdetectors.DetectorSpring

var registrations = testhelper.RegistrationFor(detectorType)

func TestBuildReportFramework(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{`java_not_spring`, false},
		{`not_java`, false},
		{`spring`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := testhelper.Extract(t,
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
	report := testhelper.Extract(t, filepath.Join("testdata", "spring"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}
