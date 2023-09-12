package ignore_test

import (
	"testing"

	"github.com/bearer/bearer/pkg/util/ignore"
	types "github.com/bearer/bearer/pkg/util/ignore/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestGetIgnoredFingerprints(t *testing.T) {
	t.Run("bearer.ignore does not exist", func(t *testing.T) {
		ignoredFingerprints, fileExists, err := ignore.GetIgnoredFingerprints("some_path.ignore", nil)
		assert.Equal(t, map[string]types.IgnoredFingerprint{}, ignoredFingerprints)
		assert.Equal(t, false, fileExists)
		assert.Equal(t, nil, err)
	})
}

func TestMergeIgnoredFingerprints(t *testing.T) {
	tests := []struct {
		Name                 string
		FingerprintsToIgnore map[string]types.IgnoredFingerprint
		IgnoredFingerprints  map[string]types.IgnoredFingerprint
		Force                bool
		Want                 []string
		ShouldSucceed        bool
	}{
		{
			Name: "Happy path - no duplicates",
			FingerprintsToIgnore: map[string]types.IgnoredFingerprint{
				"123": {
					IgnoredAt: "2023-08-28T09:30:01Z",
				},
			},
			IgnoredFingerprints: map[string]types.IgnoredFingerprint{
				"456": {
					IgnoredAt: "2023-08-28T09:30:01Z",
				},
			},
			Force:         false,
			Want:          []string{"123", "456"},
			ShouldSucceed: true,
		},
		{
			Name: "Duplicate entries",
			FingerprintsToIgnore: map[string]types.IgnoredFingerprint{
				"123": {
					IgnoredAt: "2023-08-28T09:30:01Z",
				},
			},
			IgnoredFingerprints: map[string]types.IgnoredFingerprint{
				"123": {
					IgnoredAt: "2023-08-28T09:30:01Z",
				},
				"456": {
					IgnoredAt: "2023-08-28T09:30:01Z",
				},
			},
			Force:         false,
			Want:          []string{"123", "456"},
			ShouldSucceed: false,
		},
		{
			Name: "Duplicate entries with force flag set",
			FingerprintsToIgnore: map[string]types.IgnoredFingerprint{
				"123": {
					IgnoredAt: "2",
				},
			},
			IgnoredFingerprints: map[string]types.IgnoredFingerprint{
				"123": {
					IgnoredAt: "2",
				},
				"456": {
					IgnoredAt: "2",
				},
			},
			Force:         true,
			Want:          []string{"123", "456"},
			ShouldSucceed: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			ignores := testCase.IgnoredFingerprints
			err := ignore.MergeIgnoredFingerprints(
				testCase.FingerprintsToIgnore, ignores, testCase.Force,
			)

			if err != nil && testCase.ShouldSucceed {
				t.Errorf("ignore returned error %s", err)
			}

			if err == nil && !testCase.ShouldSucceed {
				t.Errorf("expected error for test case %s but none was returned", testCase.Name)
			}

			assert.ElementsMatch(t, testCase.Want, maps.Keys(ignores))
		})
	}
}
