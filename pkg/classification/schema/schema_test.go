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
					Name:       "User",
					SimpleType: reportschema.SimpleTypeObject,
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
					Name:       "User",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "First Name",
							SimpleType: reportschema.SimpleTypeBool, // wrong type
						},
						{
							Name:       "sku_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "strokeOpacity",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "velocity",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "primaryDisplay",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "primaryDisplay",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "primaryDisplay",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "primaryDisplay",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "primaryDisplay",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "User",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "birthdate",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "First Name",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "SellerFiscalInformation",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "personName",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "_customerName",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "SellerFiscalInformation",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "erp_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "brand_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "company_front_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "_containerName",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "Applicant",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "user_id",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "bank_accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "credit_card_number",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "bank_accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "credit_card_number",
							SimpleType: reportschema.SimpleTypeBool,
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
					Name:       "patients",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "place_of_birth_unknown",
							SimpleType: reportschema.SimpleTypeBool,
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
					Name:       "accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "first_name",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "first_name",
							SimpleType: reportschema.SimpleTypeBool,
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
					Name:       "accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "foo",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "bar",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "prop types",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "lastname",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "foo",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "way_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "municipality_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "linnworks_account_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "return_company_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "brandName",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "fileName",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "InvalidNumberOfParams",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "prescription_template_id",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "last_name_score",
							SimpleType: reportschema.SimpleTypeNumber,
						},
						{
							Name:       "studentApplicant.id",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "studentapplicantTest.id",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "return_seller_id",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "seller_account_id",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "customerReturnId",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "lab_test_result_id",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "master_product_category_updated_by_user_uuid",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "foo",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "applicants_id",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "foo",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "applicants_id",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "message",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "ClientConfig",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "service_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "instance_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "matchLocationFromPhraseSetName",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "IProps",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "fontFamily",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "IProps",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "national_identifier",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "kbv_hm_diagnosis_groups",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "max_prescription_amount",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "kbv_hma_prescription_requirement_id",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "MiningPoolShares",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "propertyName",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "accountName",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "MiningPoolShares",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "full_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "_AuthorName",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "Foo",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "person_name",
							SimpleType: reportschema.SimpleTypeString,
						},
						{
							Name:       "fullName",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "Agendas",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "last_email_reminder",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "AwsRequestSigning",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "service_name",
							SimpleType: reportschema.SimpleTypeString,
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
					Name:       "MetadataCredentialsFromPlugin",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.Detection{
						{
							Name:       "service_name",
							SimpleType: reportschema.SimpleTypeString,
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
