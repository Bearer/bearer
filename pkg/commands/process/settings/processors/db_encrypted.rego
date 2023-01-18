package bearer.db_encrypted

import future.keywords

sql_encrypted contains detection if {
	detection := input.all_detections[_]

	detection.detector_type in ["schema_rb", "sql_lang_create_table"]
	detection.value.classification.decision.state == "valid"
	startswith(detection.value.normalized_field_name, "encrypted_")
}

encrypted contains detection if {
	detection := input.target_detections[_]

	detection.detector_type in ["schema_rb", "sql_lang_create_table"]

	detection.value.normalized_field_name != ""
	detection.value.normalized_object_name != ""

	some encrypted_detection in sql_encrypted
	detection.value.normalized_object_name == encrypted_detection.value.normalized_object_name
	detection.value.normalized_field_name == encrypted_detection.value.normalized_field_name
}

verified_by contains [detection, verifications] if {
	detection := input.target_detections[_]
	detection.value.field_name != ""
	detection.value.object_name != ""

	detection.detector_type in ["schema_rb", "sql_lang_create_table"]

	verifications := [verification |
		encrypted_detection := sql_encrypted[_]
		detection.value.normalized_object_name == encrypted_detection.value.normalized_object_name
		detection.value.normalized_field_name == encrypted_detection.value.normalized_field_name

		verification := {"detector": "db_encrypted"}
	]

	count(verifications) != 0
}
