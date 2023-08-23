package java_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/internal/detectors"
	"github.com/bearer/bearer/internal/detectors/java"
	"github.com/bearer/bearer/internal/parser/nodeid"

	"github.com/bearer/bearer/internal/detectors/internal/testhelper"
	detectortypes "github.com/bearer/bearer/internal/report/detectors"
)

const detectorType = detectortypes.DetectorJava

func TestDetectorReportDataTypes(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: java.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "datatype"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportPaths(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: java.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "paths"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportVariables(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: java.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "variables"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
