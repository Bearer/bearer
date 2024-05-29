package bearer.db_encrypted

import future.keywords

sql_encrypted contains detection if {
	detection := input.all_detections[_]

	detection.detector_type in input.rule.detectors
	startswith(detection.value.normalized_field_name, input.rule.auto_encrypt_prefix)
}

encrypted contains detection if {
	detection := input.target_detections[_]

	detection.detector_type in input.rule.detectors

	detection.value.normalized_field_name != ""
	detection.value.normalized_object_name != ""

	some encrypted_detection in sql_encrypted
	detection.value.normalized_object_name == encrypted_detection.value.normalized_object_name
	detection.value.field_name == encrypted_detection.value.field_name
}

verified_by contains [detection, verifications] if {
	detection := input.target_detections[_]
	detection.value.field_name != ""
	detection.value.object_name != ""

	detection.detector_type in input.rule.detectors

	verifications := [verification |
		encrypted_detection := sql_encrypted[_]
		detection.value.normalized_object_name == encrypted_detection.value.normalized_object_name
		detection.value.field_name == encrypted_detection.value.field_name

		verification := {"detector": "db_encrypted"}
	]

	count(verifications) != 0
}
