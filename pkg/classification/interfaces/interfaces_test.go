package interfaces_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/interfaces"
	"github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/parser/datatype"
	reportinterfaces "github.com/bearer/curio/pkg/report/interfaces"
	"github.com/bearer/curio/pkg/report/values"
	"github.com/bearer/curio/pkg/report/variables"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name  string
	Input reportinterfaces.Interface
	Want  interfaces.ClassifiedInterface
}

func TestSchema(t *testing.T) {
	tests := []testCase{
		testCase{
			Name: "simple path",
			Input: reportinterfaces.Interface{
				Type: interfaces.TypeURL,
				Value: &values.Value{
					Parts: []values.Part{
						&values.String{
							Type:  values.PartTypeString,
							Value: "http://",
						},
						&values.VariableReference{
							Type: values.PartTypeVariableReference,
							Identifier: variables.Identifier{
								Type: variables.Type(values.PartTypeString),
								Name: "$URL",
							},
						},
					},
				},
			},
			Want: interfaces.ClassifiedInterface{
				DataType: &datatype.DataType{
					UUID: "1",
				},
				Classification: schema.Classification{
					Name: "personal data",
				},
				Properties: map[string]schema.ClassifiedDatatype{
					"address": schema.ClassifiedDatatype{
						Classification: schema.Classification{
							Name: "personal data",
						},
						DataType: &datatype.DataType{
							UUID: "2",
						},
					},
					"age": schema.ClassifiedDatatype{
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
