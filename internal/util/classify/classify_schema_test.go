package classify_test

import (
	"testing"

	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/util/classify"
	"github.com/stretchr/testify/assert"
)

func IsDatabase(t *testing.T) {
	tests := []struct {
		Name  string
		Input detectors.Type
		Want  bool
	}{
		{
			Name:  "Rails detector",
			Input: detectors.DetectorRails,
			Want:  true,
		},
		{
			Name:  "SQL detector",
			Input: "migrations/user.txt",
			Want:  true,
		},
		{
			Name:  "Other detector",
			Input: detectors.DetectorDjango,
			Want:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Want, classify.IsDatabase(testCase.Input))
		})
	}
}

func ObjectStopWordDetected(t *testing.T) {
	tests := []struct {
		Name, Input string
		Want        bool
	}{
		{
			Name:  "Object stop word",
			Input: "prop types",
			Want:  true,
		},
		{
			Name:  "Not an object stop word",
			Input: "hello world",
			Want:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Want, classify.ObjectStopWordDetected(testCase.Input))
		})
	}
}

func PropertyStopWordDetected(t *testing.T) {
	tests := []struct {
		Name, Input string
		Want        bool
	}{
		{
			Name:  "Property stop word",
			Input: "disable click",
			Want:  true,
		},
		{
			Name:  "Not a property stop word",
			Input: "hello world",
			Want:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Want, classify.PropertyStopWordDetected(testCase.Input))
		})
	}
}

func IsExpectedIdentifierDataTypeId(t *testing.T) {
	tests := []struct {
		Name  string
		Input int
		Want  bool
	}{
		{
			Name:  "Expected identifier data type id",
			Input: 132,
			Want:  true,
		},
		{
			Name:  "Expected identifier data type id",
			Input: 13,
			Want:  true,
		},
		{
			Name:  "Not an identifier data type id",
			Input: 244,
			Want:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			assert.Equal(t, testCase.Want, classify.IsExpectedIdentifierDataTypeId(testCase.Input))
		})
	}
}
