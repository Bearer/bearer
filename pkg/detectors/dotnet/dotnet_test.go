package dotnet_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/pkg/detectors/internal/testhelper"

	reportdetectors "github.com/bearer/bearer/pkg/report/detectors"
)

const detectorType = reportdetectors.DetectorDotnet

var registrations = testhelper.RegistrationFor(detectorType)

func TestDetectorReportDbContexts(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "project", "db_contexts", "multiple"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}
