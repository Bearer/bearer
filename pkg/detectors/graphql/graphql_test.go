package graphql_test

import (
	"path/filepath"
	"testing"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/graphql"
	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/parser/nodeid"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"
	"github.com/bradleyjkemp/cupaloy"
)

var detectorType = detectortypes.DetectorGraphQL
var (
	registrations = []detectors.InitializedDetector{{Type: detectorType, Detector: graphql.New(&nodeid.IntGenerator{Counter: 0})}}
)

func TestBuildReportSchema(t *testing.T) {
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "schemas"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}
