package bearer.rails_encrypted

import future.keywords

normalized_object_name(detection) := result if {
	parts := split(detection.source.filename, "/")
	result := trim_suffix(parts[count(parts) - 1], ".rb")
}

normalized_field_name(detection) := result if {
	result := trim_prefix(detection.value.content, ":")
}

rails_encrypted contains detection if {
	detection := input.all_detections[_]
	detection.detector_type == input.rule.id
}

encrypted contains detection if {
	detection := input.target_detections[_]
	detection.detector_type in input.rule.detectors
	detection.value.field_name != ""
	detection.value.normalized_object_name != ""

	some encrypted_detection in rails_encrypted
	detection.value.field_name == normalized_field_name(encrypted_detection)
	detection.value.normalized_object_name == normalized_object_name(encrypted_detection)
}

verified_by contains [detection, verifications] if {
	detection := input.target_detections[_]
	detection.detector_type in input.rule.detectors
	detection.value.field_name != ""
	detection.value.normalized_object_name != ""

	verifications := [verification |
		encrypted_detection := rails_encrypted[_]
		detection.detector_type in input.rule.detectors

		detection.value.field_name == normalized_field_name(encrypted_detection)
		detection.value.normalized_object_name == normalized_object_name(encrypted_detection)
		verification := {
			"detector": input.rule.id,
			"filename": encrypted_detection.filename,
			"start_line_number": encrypted_detection.start_line_number,
			"start_column_number": encrypted_detection.start_column_number,
			"end_column_number": encrypted_detection.end_column_number,
		}
	]

	count(verifications) != 0
}
