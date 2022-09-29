package yamlconfig_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/report/detectors"
)

const detectorType = detectors.DetectorYamlConfig

var registrations = testhelper.RegistrationFor(detectorType)

func TestDetectorReportVariables(t *testing.T) {
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "project"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportInterfaces(t *testing.T) {
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "project"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
