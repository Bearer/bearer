package gitleaks_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/bearer/curio/pkg/detectors/gitleaks"
	detectortypes "github.com/bearer/curio/pkg/report/detectors"

	"github.com/bearer/curio/pkg/detectors"
	"github.com/bearer/curio/pkg/detectors/internal/testhelper"
	"github.com/bearer/curio/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorGitleaks

func TestSecretLeaks(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: gitleaks.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.SecretLeaks)
}
