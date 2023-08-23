package django_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/internal/detectors/internal/testhelper"
	reportdetectors "github.com/bearer/bearer/internal/report/detectors"
)

const detectorType = reportdetectors.DetectorDjango

var registrations = testhelper.RegistrationFor(detectorType)

func TestDetectorReportDatabases(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "django"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}
