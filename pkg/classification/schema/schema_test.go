package schema_test

import (
	"testing"

	reportschema "github.com/bearer/curio/pkg/report/schema"
	"github.com/stretchr/testify/assert"

	"github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/parser/datatype"
)

type testCase struct {
	Name  string
	Input datatype.DataType
	Want  schema.ClassifiedDatatype
}

func TestSchema(t *testing.T) {
	tests := []testCase{
		{
			Name: "simple path",
			Input: datatype.DataType{
				Name: "user",
				Type: reportschema.SimpleTypeObject,
				Properties: map[string]*datatype.DataType{
					"address": {
						Type: reportschema.SimpleTypeString,
						UUID: "2",
					},
					"age": {
						Type: reportschema.SimpleTypeString,
						UUID: "3",
					},
				},
				UUID: "1",
			},
			Want: schema.ClassifiedDatatype{
				DataType: &datatype.DataType{
					UUID: "1",
				},
				Classification: schema.Classification{
					Name: "personal data",
				},
				Properties: map[string]schema.ClassifiedDatatype{
					"address": {
						Classification: schema.Classification{
							Name: "personal data",
						},
						DataType: &datatype.DataType{
							UUID: "2",
						},
					},
					"age": {
						Classification: schema.Classification{},
					},
				},
			},
		},
	}

	classifier := schema.New(schema.Config{})

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := classifier.Classify(testCase.Input)
			if err != nil {
				t.Errorf("classifier returned error %s", err)
			}

			assert.Equal(t, testCase.Want, output)
		})
	}
}
