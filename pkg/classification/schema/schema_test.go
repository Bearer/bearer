package schema_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/report/detectors"
	reportschema "github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/schema/datatype"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/stretchr/testify/assert"
)

func TestSchemaObjectClassification(t *testing.T) {
	knownObjectDataType := db.DataType{
		DataCategoryName: "Unique Identifier",
		DefaultCategory:  "Identification",
		Id:               86,
		UUID:             "12d44ae0-1df7-4faf-9fb1-b46cc4b4dce9",
	}
	tests := []struct {
		Name  string
		Input schema.DataTypeDetection
		Want  schema.Classification
	}{
		{
			Name: "from vendors folder",
			Input: schema.DataTypeDetection{
				Filename:     "vendor/vendor.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &datatype.DataType{
					Name: "User",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"sku_name": &datatype.DataType{
							Name: "sku_name",
							Type: reportschema.SimpleTypeString,
							UUID: "3",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "user",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "included_in_vendor_folder",
				},
			},
		},
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
				Name:     "user",
				DataType: &knownObjectDataType,
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
				Name:     "user",
				DataType: &knownObjectDataType,
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
				Name:     "seller fiscal information",
				DataType: &knownObjectDataType,
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
				Name:     "seller fiscal information",
				DataType: &knownObjectDataType,
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
				Name:     "applicant",
				DataType: &knownObjectDataType,
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
				Name:     "bank accounts",
				DataType: &knownObjectDataType,
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
				Name:     "bank accounts",
				DataType: &knownObjectDataType,
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
				Name:     "patients",
				DataType: &knownObjectDataType,
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
				Name:     "accounts",
				DataType: &knownObjectDataType,
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
				Name:     "accounts",
				DataType: &knownObjectDataType,
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
				Name:     "accounts",
				DataType: &knownObjectDataType,
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
				Name:     "accounts",
				DataType: &knownObjectDataType,
				Decision: classify.ClassificationDecision{
					State:  classify.Potential,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "Unknown object matching stop word",
			Input: schema.DataTypeDetection{
				Filename:     "application.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "prop types",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"lastname": &datatype.DataType{
							Name: "lastname",
							Type: reportschema.SimpleTypeString,
							UUID: "2",
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "prop types",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "stop_word",
				},
			},
		},
		{
			Name: "Unknown object with invalid properties",
			Input: schema.DataTypeDetection{
				Filename:     "db/structure.sql",
				DetectorType: detectors.DetectorRails,
				Value: &datatype.DataType{
					Name: "foo",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"way_name": &datatype.DataType{
							Name: "way_name",
							Type: reportschema.SimpleTypeString,
						},
						"municipality_name": &datatype.DataType{
							Name: "municipality_name",
							Type: reportschema.SimpleTypeString,
						},
						"linnworks_account_name": &datatype.DataType{
							Name: "linnworks_account_name",
							Type: reportschema.SimpleTypeString,
						},
						"return_company_name": &datatype.DataType{
							Name: "return_company_name",
							Type: reportschema.SimpleTypeString,
						},
						"brandName": &datatype.DataType{
							Name: "brandName",
							Type: reportschema.SimpleTypeString,
						},
						"fileName": &datatype.DataType{
							Name: "fileName",
							Type: reportschema.SimpleTypeString,
						},
						"InvalidNumberOfParams": &datatype.DataType{
							Name: "InvalidNumberOfParams",
							Type: reportschema.SimpleTypeString,
						},
						"prescription_template_id": &datatype.DataType{
							Name: "prescription_template_id",
							Type: reportschema.SimpleTypeString,
						},
						"last_name_score": &datatype.DataType{
							Name: "last_name_score",
							Type: reportschema.SimpleTypeNumber,
						},
						"studentApplicant.id": &datatype.DataType{
							Name: "studentApplicant.id",
							Type: reportschema.SimpleTypeString,
						},
						"studentapplicantTest.id": &datatype.DataType{
							Name: "studentapplicantTest.id",
							Type: reportschema.SimpleTypeString,
						},
						"return_seller_id": &datatype.DataType{
							Name: "return_seller_id",
							Type: reportschema.SimpleTypeString,
						},
						"seller_account_id": &datatype.DataType{
							Name: "seller_account_id",
							Type: reportschema.SimpleTypeString,
						},
						"customerReturnId": &datatype.DataType{
							Name: "customerReturnId",
							Type: reportschema.SimpleTypeString,
						},
						"lab_test_result_id": &datatype.DataType{
							Name: "lab_test_result_id",
							Type: reportschema.SimpleTypeString,
						},
						"master_product_category_updated_by_user_uuid": &datatype.DataType{
							Name: "master_product_category_updated_by_user_uuid",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "foo",
				Decision: classify.ClassificationDecision{
					State:  classify.Potential, // database
					Reason: "unknown_data_object",
				},
			},
		},
		{
			Name: "Unknown object with db identifier but no associated object properties",
			Input: schema.DataTypeDetection{
				Filename:     "db/structure.sql",
				DetectorType: detectors.DetectorRails,
				Value: &datatype.DataType{
					Name: "foo",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"applicants_id": &datatype.DataType{
							Name: "applicants_id",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "foo",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "only_db_identifiers",
				},
			},
		},
		{
			Name: "Unknown object with db identifier and an associated object property",
			Input: schema.DataTypeDetection{
				Filename:     "db/structure.sql",
				DetectorType: detectors.DetectorRails,
				Value: &datatype.DataType{
					Name: "foo",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"applicants_id": &datatype.DataType{
							Name: "applicants_id",
							Type: reportschema.SimpleTypeString,
						},
						"message": &datatype.DataType{
							Name: "message",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "foo",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "invalid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "ClientConfig object",
			Input: schema.DataTypeDetection{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "ClientConfig",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"service_name": &datatype.DataType{
							Name: "service_name",
							Type: reportschema.SimpleTypeString,
						},
						"instance_name": &datatype.DataType{
							Name: "instance_name",
							Type: reportschema.SimpleTypeString,
						},
						"matchLocationFromPhraseSetName": &datatype.DataType{
							Name: "matchLocationFromPhraseSetName",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "client config",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "unknown_data_object",
				},
			},
		},
		{
			Name: "IProps object with invalid properties",
			Input: schema.DataTypeDetection{
				Filename:     "application.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "IProps",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"fontFamily": &datatype.DataType{
							Name: "fontFamily",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "i props",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "unknown_data_object",
				},
			},
		},
		{
			Name: "IProps object with valid properties",
			Input: schema.DataTypeDetection{
				Filename:     "application.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "IProps",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"national_identifier": &datatype.DataType{
							Name: "national_identifier",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "i props",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "invalid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "kbv_hm_diagnosis_groups object",
			Input: schema.DataTypeDetection{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "kbv_hm_diagnosis_groups",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"max_prescription_amount": &datatype.DataType{
							Name: "max_prescription_amount",
							Type: reportschema.SimpleTypeString,
						},
						"kbv_hma_prescription_requirement_id": &datatype.DataType{
							Name: "kbv_hma_prescription_requirement_id",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "kbv hm diagnosis groups",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "unknown_data_object",
				},
			},
		},
		{
			Name: "MiningPoolShares object",
			Input: schema.DataTypeDetection{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "MiningPoolShares",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"propertyName": &datatype.DataType{
							Name: "propertyName",
							Type: reportschema.SimpleTypeString,
						},
						"accountName": &datatype.DataType{
							Name: "accountName",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "mining pool shares",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "unknown_data_object",
				},
			},
		},
		{
			Name: "MiningPoolShares object",
			Input: schema.DataTypeDetection{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "MiningPoolShares",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"full_name": &datatype.DataType{
							Name: "full_name",
							Type: reportschema.SimpleTypeString,
						},
						"_AuthorName": &datatype.DataType{
							Name: "_AuthorName",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "mining pool shares",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "invalid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "With valid extended data properties",
			Input: schema.DataTypeDetection{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "Foo",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"person_name": &datatype.DataType{
							Name: "person_name",
							Type: reportschema.SimpleTypeString,
						},
						"fullName": &datatype.DataType{
							Name: "fullName",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "foo",
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "invalid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "With valid extended data properties",
			Input: schema.DataTypeDetection{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "Agendas",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"last_email_reminder": &datatype.DataType{
							Name: "last_email_reminder",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "agendas",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "unknown_data_object",
				},
			},
		},
		{
			Name: "With invalid object / field combination - AWS case",
			Input: schema.DataTypeDetection{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "AwsRequestSigning",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"service_name": &datatype.DataType{
							Name: "service_name",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "aws request signing",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "unknown_data_object",
				},
			},
		},
		{
			Name: "With invalid object / field combination - MetadataCredentialsFromPlugin case",
			Input: schema.DataTypeDetection{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &datatype.DataType{
					Name: "MetadataCredentialsFromPlugin",
					UUID: "1",
					Type: reportschema.SimpleTypeObject,
					Properties: map[string]datatype.DataTypable{
						"service_name": &datatype.DataType{
							Name: "service_name",
							Type: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "metadata credentials from plugin",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "unknown_data_object",
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
			output := classifier.Classify(testCase.Input)

			assert.Equal(t, testCase.Want, output.Classification)
		})
	}
}
