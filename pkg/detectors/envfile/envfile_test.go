package envfile_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/pkg/detectors/internal/testhelper"
	"github.com/bearer/bearer/pkg/report/detectors"
)

const detectorType = detectors.DetectorEnvFile

var registrations = testhelper.RegistrationFor(detectorType)

func TestDetectorReportVariables(t *testing.T) {
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "variables"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
