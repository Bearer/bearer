package schema_test

import (
	"testing"

	"github.com/bearer/bearer/internal/classification/db"
	"github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/report/detectors"
	reportschema "github.com/bearer/bearer/internal/report/schema"
	"github.com/bearer/bearer/internal/util/classify"
	"github.com/stretchr/testify/assert"
)

func TestSchemaObjectClassification(t *testing.T) {
	tests := []struct {
		Name  string
		Input schema.ClassificationRequest
		Want  schema.Classification
	}{
		{
			Name: "from vendors folder",
			Input: schema.ClassificationRequest{
				Filename:     "vendor/vendor.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "User",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "User",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
				Name:        "user",
				DataType:    nil,
				SubjectName: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with valid properties",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "User",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
				Name:        "user",
				DataType:    nil,
				SubjectName: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object is SellerFiscalInformation with valid properties",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "SellerFiscalInformation",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
				Name:        "seller fiscal information",
				DataType:    nil,
				SubjectName: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object is SellerFiscalInformation with no valid properties",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "SellerFiscalInformation",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
				Name:        "seller fiscal information",
				DataType:    nil,
				SubjectName: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that match a db identifier",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "Applicant",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "user_id",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:        "applicant",
				DataType:    nil,
				SubjectName: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object with properties that matches exclude patterns",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "users",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "last_synced_mobile_app",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:        "users",
				DataType:    nil,
				SubjectName: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that match exclude patterns",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "bank_accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "credit_card_number",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:        "bank accounts",
				DataType:    nil,
				SubjectName: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object with properties that match exclude patterns but are the wrong type (bool) - case #1",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "bank_accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "credit_card_number",
							SimpleType: reportschema.SimpleTypeBool,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:        "bank accounts",
				SubjectName: nil,
				DataType:    nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that match exclude patterns but are the wrong type (bool) - case #2",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "patients",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "place_of_birth_unknown",
							SimpleType: reportschema.SimpleTypeBool,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:        "patients",
				SubjectName: nil,
				DataType:    nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that do not match exclude patterns",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "first_name",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:     "accounts",
				DataType: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Valid,
					Reason: "valid_object_with_valid_properties",
				},
			},
		},
		{
			Name: "known object with properties that do not match exclude criteria but are the wrong type (bool)",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "first_name",
							SimpleType: reportschema.SimpleTypeBool,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:     "accounts",
				DataType: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "known object with properties that do not match any patterns",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "foo",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:     "accounts",
				DataType: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "Known object - database detection - with invalid properties",
			Input: schema.ClassificationRequest{
				Filename:     "db/schema.rb",
				DetectorType: detectors.DetectorSchemaRb,
				Value: &schema.ClassificationRequestDetection{
					Name:       "accounts",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "bar",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name:     "accounts",
				DataType: nil,
				Decision: classify.ClassificationDecision{
					State:  classify.Potential,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		},
		{
			Name: "Unknown object matching stop word",
			Input: schema.ClassificationRequest{
				Filename:     "application.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "prop types",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "db/structure.sql",
				DetectorType: detectors.DetectorRails,
				Value: &schema.ClassificationRequestDetection{
					Name:       "foo",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "db/structure.sql",
				DetectorType: detectors.DetectorRails,
				Value: &schema.ClassificationRequestDetection{
					Name:       "foo",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "db/structure.sql",
				DetectorType: detectors.DetectorRails,
				Value: &schema.ClassificationRequestDetection{
					Name:       "foo",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "ClientConfig",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "application.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "IProps",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "application.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "IProps",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "kbv_hm_diagnosis_groups",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "MiningPoolShares",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "MiningPoolShares",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "Foo",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "Agendas",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "AwsRequestSigning",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
			Input: schema.ClassificationRequest{
				Filename:     "app.js",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "MetadataCredentialsFromPlugin",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
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
		{
			Name: "full name in unknown object - case 1",
			Input: schema.ClassificationRequest{
				Filename:     "lib/api/validations/validators/absence.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "@scope",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "full_name",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "@scope",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "unknown_data_object",
				},
			},
		},
		{
			Name: "full name in unknown object - case 1",
			Input: schema.ClassificationRequest{
				Filename:     "lib/gitlab/bitbucket_import/importer.rb",
				DetectorType: detectors.DetectorRuby,
				Value: &schema.ClassificationRequestDetection{
					Name:       "project",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "full_name",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			},
			Want: schema.Classification{
				Name: "project",
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

	t.Run("incorrect Personal health history classification", func(t *testing.T) {
		output := classifier.Classify(
			schema.ClassificationRequest{
				Filename:     "src/components/order/client-transaction-modals/withdraw-vault/index.jsx",
				DetectorType: detectors.DetectorJavascript,
				Value: &schema.ClassificationRequestDetection{
					Name:       "selfManagedCustodySigningSteps",
					SimpleType: reportschema.SimpleTypeObject,
					Properties: []*schema.ClassificationRequestDetection{
						{
							Name:       "cold",
							SimpleType: reportschema.SimpleTypeString,
						},
					},
				},
			})

		assert.Equal(t, &schema.ClassifiedDatatype{
			Name: "selfManagedCustodySigningSteps",
			Properties: []*schema.ClassifiedDatatype{
				{
					Name:       "cold",
					Properties: nil,
					Classification: schema.Classification{
						Name:     "cold",
						DataType: nil,
						Decision: classify.ClassificationDecision{
							State:  classify.Invalid,
							Reason: "invalid_property",
						},
					},
				},
			},
			Classification: schema.Classification{
				Name: "self managed custody signing steps",
				Decision: classify.ClassificationDecision{
					State:  classify.Invalid,
					Reason: "valid_object_with_invalid_properties",
				},
			},
		}, output)
	})
}
