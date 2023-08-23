package gitleaks_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/bearer/internal/detectors/gitleaks"
	detectortypes "github.com/bearer/bearer/internal/report/detectors"

	"github.com/bearer/bearer/internal/detectors"
	"github.com/bearer/bearer/internal/detectors/internal/testhelper"
	"github.com/bearer/bearer/internal/parser/nodeid"
)

const detectorType = detectortypes.DetectorGitleaks

func TestSecretLeaks(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: gitleaks.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.SecretLeaks)
}
