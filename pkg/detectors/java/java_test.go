package java_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/java"
	"github.com/bearer/curio/pkg/parser/nodeid"

	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"
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
