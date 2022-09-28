package simple_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/report/detectors"

	"github.com/bradleyjkemp/cupaloy"
)

const detectorType = detectors.DetectorSimple

var registrations = testhelper.RegistrationFor(detectorType)

func TestBuildReportInterfaces(t *testing.T) {
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "project"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
