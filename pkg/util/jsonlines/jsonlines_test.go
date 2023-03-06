package jsonlines_test

import (
	"os"
	"testing"

	"github.com/bearer/bearer/pkg/util/jsonlines"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestJsonlines(t *testing.T) {
	type TestFile struct {
		Name  string
		Order int
	}

	originalValue := []TestFile{
		{
			Name:  "test struct 1",
			Order: 1,
		},
		{
			Name:  "test struct 2",
			Order: 2,
		},
		{
			Name:  "test struct 3 \n",
			Order: 3,
		},
	}

	file, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatalf("failed to create temp file %s", err)
	}
	defer file.Close()

	err = jsonlines.Encode(file, &originalValue)
	if err != nil {
		t.Fatalf("failed to encode %s", err)
	}

	file.Seek(0, 0)

	decodedValue := make([]interface{}, 0)

	err = jsonlines.Decode(file, &decodedValue)
	if err != nil {
		t.Fatalf("failed to encode %s", err)
	}

	decodedObjects := []TestFile{}

	for _, v := range decodedValue {
		var object TestFile
		err := mapstructure.Decode(v, &object)
		if err != nil {
			t.Fatalf("failed to decode mapstructure %s", err)
		}

		decodedObjects = append(decodedObjects, object)
	}

	assert.Equal(t, originalValue, decodedObjects)
}
