package schema_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/report/detectors"
	reportschema "github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/stretchr/testify/assert"

	"github.com/bearer/curio/pkg/classification/schema"
)

func TestSchemaObjectClassification(t *testing.T) {
	tests := []struct {
		Name  string
		Input schema.DataTypeDetection
		Want  schema.Classification
	}{
		{
			Name: "known object with no valid properties",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "User",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"First Name": &datatype.DataType{
							Name: "First Name",
							Type: reportschema.SimpleTypeBool, // wrong type
							UUID: "2",
						},
						"sku_name": &datatype.DataType{
							Name: "sku_name",
							Type: reportschema.SimpleTypeString,
							UUID: "3",
						},
						"strokeOpacity": &datatype.DataType{
							Name: "strokeOpacity",
							Type: reportschema.SimpleTypeString,
							UUID: "4",
						},
						"velocity": &datatype.DataType{
							Name: "velocity",
							Type: reportschema.SimpleTypeString,
							UUID: "5",
						},
						"primaryDisplay": &datatype.DataType{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
							UUID: "6",
						},
						"requestPassword": &datatype.DataType{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
							UUID: "7",
						},
						"resetPassword": &datatype.DataType{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
							UUID: "8",
						},
						"resolveAmbiguousRoles": &datatype.DataType{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
							UUID: "9",
						},
						"combinePartialBlindedSignatures": &datatype.DataType{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
							UUID: "10",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "user",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with valid properties",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "User",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"birthdate": &datatype.DataType{
							Name: "birthdate",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
						"First Name": &datatype.DataType{
							Name: "First Name",
							Type: reportschema.SimpleTypeString,
							UUID: "3",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "user",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object is SellerFiscalInformation with valid properties",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "SellerFiscalInformation",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"name": &datatype.DataType{
							Name: "name",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
						"personName": &datatype.DataType{
							Name: "personName",
							Type: reportschema.SimpleTypeString,
							UUID: "3",
						},
						"_customerName": &datatype.DataType{
							Name: "_customerName",
							Type: reportschema.SimpleTypeString,
							UUID: "3",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "seller fiscal information",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object is SellerFiscalInformation with no valid properties",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "SellerFiscalInformation",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"erp_name": &datatype.DataType{
							Name: "erp_name",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
						"brand_name": &datatype.DataType{
							Name: "brand_name",
							Type: reportschema.SimpleTypeString,
							UUID: "3",
						},
						"company_front_name": &datatype.DataType{
							Name: "company_front_name",
							Type: reportschema.SimpleTypeString,
							UUID: "3",
						},
						"_containerName": &datatype.DataType{
							Name: "_containerName",
							Type: reportschema.SimpleTypeString,
							UUID: "4",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "seller fiscal information",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that match a db identifier",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "Applicant",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"user_id": &datatype.DataType{
							Name: "user_id",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "applicant",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object with properties that match exclude patterns",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "bank_accounts",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"credit_card_number": &datatype.DataType{
							Name: "credit_card_number",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "bank accounts",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object with properties that match exclude patterns but are the wrong type (bool) - case #1",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "bank_accounts",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"credit_card_number": &datatype.DataType{
							Name: "credit_card_number",
							Type: reportschema.SimpleTypeBool,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "bank accounts",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that match exclude patterns but are the wrong type (bool) - case #2",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "patients",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"place_of_birth_unknown": &datatype.DataType{
							Name: "place_of_birth_unknown",
							Type: reportschema.SimpleTypeBool,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "patients",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that do not match exclude patterns",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "accounts",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"first_name": &datatype.DataType{
							Name: "first_name",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "accounts",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object with properties that do not match exclude criteria but are the wrong type (bool)",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "accounts",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"first_name": &datatype.DataType{
							Name: "first_name",
							Type: reportschema.SimpleTypeBool,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "accounts",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that do not match any patterns",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "accounts",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"foo": &datatype.DataType{
							Name: "foo",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "accounts",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "Known object - database detection - with invalid properties",
			Input: schema.DataTypeDetection{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorSQL,
				Value: &datatype.DataType{
					Name: "accounts",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"bar": &datatype.DataType{
							Name: "bar",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "accounts",
				Decision: classify.ClassificationDecision{
					State:  classify.Potential,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
	}
	classifier := schema.New(
		schema.Config{
			DataTypes:                      db.Default().DataTypes,
			DataTypeClassificationPatterns: db.Default().DataTypeClassificationPatterns,
			KnownPersonObjectPatterns:      db.Default().KnownPersonObjectPatterns,
		},
	)

	for _, testCase := range tests {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := classifier.Classify(testCase.Input)
			if err != nil {
				t.Errorf("classifier returned error %s", err)
			}

			assert.Equal(t, testCase.Want, output.Classification)
		})
	}
}
