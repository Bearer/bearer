package classify_test

import (
	"testing"

	"github.com/bearer/bearer/pkg/report/detectors"
	"github.com/bearer/bearer/pkg/util/classify"
	"github.com/stretchr/testify/assert"
)

func TestIsVendored(t *testing.T) {
	tests := []struct {
		Name, Input string
		Want        bool
	}{
		{
			Name:  "In a vendor folder",
			Input: "/vendor/package.js",
			Want:  true,
		},
		{
			Name:  "In a migrations folder",
			Input: "migrations/user.txt",
			Want:  true,
		},
		{
			Name:  "In public folder",
			Input: "app/public/package.js",
			Want:  true,
		},
		{
			Name:  "Not in any vendor folder",
			Input: "db/schema.rb",
			Want:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Want, classify.IsVendored(testCase.Input))
		})
	}
}

func TestIsPotentialDetector(t *testing.T) {
	tests := []struct {
		Name  string
		Input detectors.Type
		Want  bool
	}{
		{
			Name:  "ENV file detector type",
			Input: detectors.DetectorEnvFile,
			Want:  true,
		},
		{
			Name:  "YAML config detector type",
			Input: detectors.DetectorYamlConfig,
			Want:  true,
		},
		{
			Name:  "Other detector type",
			Input: detectors.DetectorCSharp,
			Want:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Want, classify.IsPotentialDetector(testCase.Input))
		})
	}
}
