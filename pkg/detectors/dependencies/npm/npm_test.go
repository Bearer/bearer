package npm_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/pkg/detectors/internal/testhelper"
	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bradleyjkemp/cupaloy"
)

const detectorType = detectors.DetectorDependencies

var registrations = testhelper.RegistrationFor(detectorType)

func TestDependenciesReport(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata"), registrations, detectorType)
	cupaloy.SnapshotT(t, report.Dependencies)
}
