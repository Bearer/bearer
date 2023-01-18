package bearer.encrypted_verified

import future.keywords

ruby_encrypted contains detection if {
	detection := input.all_detections[_]
	detection.detector_type == "ruby_rails_default_encryption"
	detection.value.classification.decision.state == "valid"
}

encrypted contains detection if {
	detection := input.target_detections[_]
	detection.detector_type in ["schema_rb", "sql_lang_create_table"]
	detection.value.normalized_field_name != ""
	detection.value.normalized_object_name != ""

	some encrypted_detection in ruby_encrypted
	detection.value.normalized_object_name == encrypted_detection.value.normalized_object_name
	detection.value.normalized_field_name == encrypted_detection.value.normalized_field_name
}

verified_by contains [detection, verifications] if {
	detection := input.target_detections[_]
	detection.detector_type in ["schema_rb", "sql_lang_create_table"]
	detection.value.normalized_field_name != ""
	detection.value.normalized_object_name != ""

	verifications := [verification |
		encrypted_detection := ruby_encrypted[_]
		detection.detector_type in ["schema_rb", "sql_lang_create_table"]
		detection.value.normalized_field_name != ""
		detection.value.normalized_object_name != ""
		detection.value.normalized_object_name == encrypted_detection.value.normalized_object_name
		detection.value.normalized_field_name == encrypted_detection.value.normalized_field_name
		verification := {
			"detector": "ruby_rails_default_encryption",
			"filename": encrypted_detection.source.filename,
			"line_number": encrypted_detection.source.line_number,
		}
	]

	count(verifications) != 0
}
