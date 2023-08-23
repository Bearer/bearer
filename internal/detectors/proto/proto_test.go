package proto_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/bearer/internal/detectors"
	"github.com/bearer/bearer/internal/detectors/proto"
	"github.com/bearer/bearer/internal/parser/nodeid"
	detectortypes "github.com/bearer/bearer/internal/report/detectors"

	"github.com/bearer/bearer/internal/detectors/internal/testhelper"
	"github.com/bradleyjkemp/cupaloy"
)

var detectorType = detectortypes.DetectorGraphQL
var (
	registrations = []detectors.InitializedDetector{{Type: detectorType, Detector: proto.New(&nodeid.IntGenerator{Counter: 0})}}
)

func TestBuildReportSchema(t *testing.T) {
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "protos"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
