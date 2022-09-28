package php_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/curio/pkg/detectors/php"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorPHP

func TestDetectorReportDatatype(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorPHP,
		Detector: php.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "datatype"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

// we should ignore binary embeding php, they cause tree sitter to explode
func TestDetectorReportIgnore(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorPHP,
		Detector: php.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "ignore"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportInterfaces1(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorPHP,
		Detector: php.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "paths"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportInterfaces2(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorPHP,
		Detector: php.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "variables"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportContext(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorPHP,
		Detector: php.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "context"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
