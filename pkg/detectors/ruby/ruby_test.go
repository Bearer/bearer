package ruby_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/curio/pkg/detectors/ruby"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorRuby

func TestDetectorReportDatatype(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorRuby,
		Detector: ruby.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "datatype"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportInterfacesPaths(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorRuby,
		Detector: ruby.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "paths"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportInterfacesVariables(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorRuby,
		Detector: ruby.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "variables"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
