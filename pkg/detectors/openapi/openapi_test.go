package openapi_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/detectors/openapi"
	"github.com/bearer/curio/pkg/parser/nodeid"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"
	"github.com/bradleyjkemp/cupaloy"
)

var detectorType = detectortypes.DetectorOpenAPI
var (
	registrations = []detectors.InitializedDetector{{Type: detectortypes.DetectorOpenAPI, Detector: openapi.New(&nodeid.IntGenerator{Counter: 0})}}
)

func TestDetectorV3json(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "v3json"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Detections)
}

func TestDetectorV3yaml(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "v3yaml"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Detections)
}

func TestDetectorV2json(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "v2json"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Detections)
}

func TestDetectorV2yaml(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "v2yaml"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Detections)
}
