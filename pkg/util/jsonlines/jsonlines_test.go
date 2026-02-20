package jsonlines_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/bearer/bearer/pkg/util/jsonlines"
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

	_, err = file.Seek(0, 0)
	if err != nil {
		t.Fatalf("failed to seek file to begning %s", err)
	}

	decodedValue := make([]interface{}, 0)

	err = jsonlines.Decode(file, &decodedValue)
	if err != nil {
		t.Fatalf("failed to encode %s", err)
	}

	decodedObjects := []TestFile{}

	for _, v := range decodedValue {
		var object TestFile
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("failed to marshal to json %s", err)
		}
		err = json.Unmarshal(b, &object)
		if err != nil {
			t.Fatalf("failed to unmarshal from json %s", err)
		}

		decodedObjects = append(decodedObjects, object)
	}

	assert.Equal(t, originalValue, decodedObjects)
}
