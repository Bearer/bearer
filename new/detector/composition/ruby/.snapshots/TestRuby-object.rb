[
	{
		"DetectorType": "logger",
		"MatchNode": {},
		"Data": {
			"Pattern": "logger.info($\u003cDATA_TYPE\u003e)\n",
			"Datatypes": [
				{
					"DetectorType": "datatype",
					"MatchNode": {},
					"Data": {
						"Name": "user",
						"Classification": {
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
						"Properties": [
							{
								"Name": "name",
								"Detection": {
									"DetectorType": "object",
									"MatchNode": {},
									"Data": {
										"Name": "name"
									}
								},
								"Classification": {
									"name": "user",
									"decision": {
										"state": "valid",
										"reason": "valid_object_with_valid_properties"
									}
								},
								"Properties": [
									{
										"Name": "name",
										"Detection": {
											"MatchNode": {},
											"ContextNode": null,
											"Data": {
												"Name": "name"
											}
										},
										"Classification": {
											"name": "name",
											"subject_name": "User",
											"data_type": {
												"name": "Fullname",
												"uuid": "1617291b-bc22-4267-ad5e-8054b2505958",
												"category_uuid": "14124881-6b92-4fc5-8005-ea7c1c09592e"
											},
											"decision": {
												"state": "valid",
												"reason": "known_pattern"
											}
										}
									}
								}
							}
						]
					}
				}
			]
		}
	}
]
