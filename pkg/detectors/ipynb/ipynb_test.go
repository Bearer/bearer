package ipynb_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/curio/pkg/detectors/ipynb"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorIPYNB

func TestDetectorReportInterfaces(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorIPYNB,
		Detector: ipynb.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "notebooks"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
