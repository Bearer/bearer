package python_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/curio/pkg/detectors/python"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorPython

func TestDetectorReportDatatype(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorPython,
		Detector: python.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "datatype"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportInterfacePaths(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorPython,
		Detector: python.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "paths"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportInterfaceProjects(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorPython,
		Detector: python.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "project"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
