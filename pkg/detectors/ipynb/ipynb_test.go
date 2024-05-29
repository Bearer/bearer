package ipynb_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/pkg/detectors/ipynb"
	detectortypes "github.com/bearer/bearer/pkg/report/detectors"

	"github.com/bearer/bearer/pkg/detectors"
	"github.com/bearer/bearer/pkg/detectors/internal/testhelper"
	"github.com/bearer/bearer/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorIPYNB

func TestDetectorReportInterfaces(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorIPYNB,
		Detector: ipynb.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "notebooks"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
