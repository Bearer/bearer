package sql_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/sql"
	"github.com/bearer/curio/pkg/parser/nodeid"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"

	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bradleyjkemp/cupaloy"
)

var detectorType = detectortypes.DetectorGraphQL

func TestCreateTable(t *testing.T) {
	var (
		registrations = []detectors.InitializedDetector{{detectorType, sql.New(&nodeid.IntGenerator{Counter: 0})}}
	)
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "create_table"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestCreateView(t *testing.T) {
	var (
		registrations = []detectors.InitializedDetector{{detectorType, sql.New(&nodeid.IntGenerator{Counter: 0})}}
	)
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "create_view"), registrations, detectorType)

	jsonOutput, err := json.MarshalIndent(detectorReport.CreateView, "", "\t")
	if err != nil {
		t.Error(err)
	}

	cupaloy.SnapshotT(t, jsonOutput)
}
