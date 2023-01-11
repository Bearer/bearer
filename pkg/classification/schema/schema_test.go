package schema_test

import (
	"testing"

	"github.com/bearer/curio/pkg/classification/db"
	"github.com/bearer/curio/pkg/classification/schema"
	"github.com/bearer/curio/pkg/report/detectors"
	reportschema "github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/util/classify"
	"github.com/stretchr/testify/assert"
)

func TestSchemaObjectClassification(t *testing.T) {
	knownObjectDataType := db.DataType{
		Name:         "Unique Identifier",
		CategoryUUID: "14124881-6b92-4fc5-8005-ea7c1c09592e",
		UUID:         "12d44ae0-1df7-4faf-9fb1-b46cc4b4dce9",
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
				Value: &schema.Detection{
					Name: "User",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "sku_name",
							SimpleType: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "User",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "First Name",
							Type: reportschema.SimpleTypeBool, // wrong type
						},
						{
							Name: "sku_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "strokeOpacity",
							Type: reportschema.SimpleTypeString,
						},
						{Name: "velocity",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "primaryDisplay",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "User",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "birthdate",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "First Name",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "SellerFiscalInformation",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "personName",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "_customerName",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "SellerFiscalInformation",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "erp_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "brand_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "company_front_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "_containerName",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "Applicant",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "user_id",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "bank_accounts",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "credit_card_number",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "bank_accounts",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "credit_card_number",
							Type: reportschema.SimpleTypeBool,
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
				Value: &schema.Detection{
					Name: "patients",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "place_of_birth_unknown",
							Type: reportschema.SimpleTypeBool,
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
				Value: &schema.Detection{
					Name: "accounts",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "first_name",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "accounts",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "first_name",
							Type: reportschema.SimpleTypeBool,
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
				Value: &schema.Detection{
					Name: "accounts",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "foo",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "accounts",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "bar",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "prop types",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "lastname",
							Type: reportschema.SimpleTypeString,
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
				Value: &schema.Detection{
					Name: "foo",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "way_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "municipality_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "linnworks_account_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "return_company_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "brandName",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "fileName",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "InvalidNumberOfParams",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "prescription_template_id",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "last_name_score",
							Type: reportschema.SimpleTypeNumber,
						},
						{
							Name: "studentApplicant.id",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "studentapplicantTest.id",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "return_seller_id",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "seller_account_id",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "customerReturnId",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "lab_test_result_id",
							Type: reportschema.SimpleTypeString,
						},
						{
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
				Value: &schema.Detection{
					Name: "foo",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
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
				Value: &schema.Detection{
					Name: "foo",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "applicants_id",
							Type: reportschema.SimpleTypeString,
						},
						{
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
				Value: &schema.Detection{
					Name: "ClientConfig",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "service_name",
							Type: reportschema.SimpleTypeString,
						},
						{
							Name: "instance_name",
							Type: reportschema.SimpleTypeString,
						},
						{
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
				Value: &schema.Detection{
					Name: "IProps",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
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
				Value: &schema.Detection{
					Name: "IProps",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
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
				Value: &schema.Detection{
					Name: "kbv_hm_diagnosis_groups",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "max_prescription_amount",
							Type: reportschema.SimpleTypeString,
						},
						{
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
				Value: &schema.Detection{
					Name: "MiningPoolShares",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "propertyName",
							Type: reportschema.SimpleTypeString,
						},
						{
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
				Value: &schema.Detection{
					Name: "MiningPoolShares",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "full_name",
							Type: reportschema.SimpleTypeString,
						},
						{
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
				Value: &schema.Detection{
					Name: "Foo",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name: "person_name",
							Type: reportschema.SimpleTypeString,
						},
						{
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
				Value: &schema.Detection{
					Name: "Agendas",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
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
				Value: &schema.Detection{
					Name: "AwsRequestSigning",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
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
				Value: &schema.Detection{
					Name: "MetadataCredentialsFromPlugin",
					Type: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
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
