package schema_test

import (
	"encoding/json"
	"testing"

	reportschema "github.com/bearer/curio/pkg/report/schema"

	"github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/util/jsonlinesreader"
	"github.com/bradleyjkemp/cupaloy"
)

func TestSchemaJSON(t *testing.T) {
	config := schema.Config{}
	classifier := schema.New(config)

	var output []schema.ClassifiedSchema

	inputData, err := jsonlinesreader.ReadAllText("./testdata/schema.jsonl")
	if err != nil {
		t.Fatalf("failed to read input %e", err)
	}

	for _, interfaceData := range inputData {
		var toClassify reportschema.Schema
		json.Unmarshal([]byte(interfaceData), &toClassify)

		classifiedSchema, err := classifier.Classify(toClassify)
		if err != nil {
			t.Fatalf("failed to classify schema:%s with error:%e", interfaceData, err)
		}

		output = append(output, classifiedSchema)
	}

	jsonOuput, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("failed to marshal json %e", err)
	}

	cupaloy.SnapshotT(t, string(jsonOuput))
}
