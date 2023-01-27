[
	{
		"CustomDetector": "logger",
		"DetectionType": "custom_classified",
		"Source": {
			"filename": "object.rb",
			"language": "Ruby",
			"language_type": "programming",
			"line_number": 1,
			"column_number": 13,
			"text": null
		},
		"Value": {
			"object_name": "user",
			"field_name": "",
			"field_type": "",
			"field_type_simple": "",
			"classification": {
				"name": "user",
				"data_type": {
					"name": "Unique Identifier",
					"uuid": "12d44ae0-1df7-4faf-9fb1-b46cc4b4dce9",
					"category_uuid": "14124881-6b92-4fc5-8005-ea7c1c09592e"
				},
				"decision": {
					"state": "valid",
					"reason": "valid_object_with_valid_properties"
				}
			},
			"parent": {
				"line_number": 1,
				"content": "user.name"
			},
			"normalized_object_name": "user"
		}
	},
	{
		"CustomDetector": "logger",
		"DetectionType": "custom_classified",
		"Source": {
			"filename": "object.rb",
			"language": "Ruby",
			"language_type": "programming",
			"line_number": 1,
			"column_number": 13,
			"text": null
		},
		"Value": {
			"object_name": "user",
			"field_name": "name",
			"field_type": "",
			"field_type_simple": "",
			"classification": {
				"name": "name",
				"data_type": {
					"name": "Fullname",
					"uuid": "1617291b-bc22-4267-ad5e-8054b2505958",
					"category_uuid": "14124881-6b92-4fc5-8005-ea7c1c09592e"
				},
				"decision": {
					"state": "valid",
					"reason": "known_pattern"
				}
			},
			"parent": {
				"line_number": 1,
				"content": "user.name"
			},
			"normalized_object_name": "user",
			"normalized_field_name": "name"
		}
	}
]
