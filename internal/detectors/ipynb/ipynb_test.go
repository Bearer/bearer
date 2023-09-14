package ipynb_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/internal/detectors/ipynb"
	detectortypes "github.com/bearer/bearer/internal/report/detectors"

	"github.com/bearer/bearer/internal/detectors"
	"github.com/bearer/bearer/internal/detectors/internal/testhelper"
	"github.com/bearer/bearer/internal/parser/nodeid"
)

const detectorType = detectortypes.DetectorIPYNB

func TestDetectorReportInterfaces(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorIPYNB,
		Detector: ipynb.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "notebooks"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
