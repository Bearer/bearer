package bearer.encrypted_verified

import future.keywords

ruby_encrypted[detection] {
    detection := input.all_detections[_]
    detection.detector_type == "detect_encrypted_ruby_class_properties"
    detection.value.classification.decision.state == "valid"
}

encrypted[detection] {
    detection := input.target_detections[_]
    detection.value.field_name != ""
    detection.value.object_name != ""

    some encrypted_detection in ruby_encrypted
    detection.value.object_name == encrypted_detection.value.object_name
    detection.value.field_name == encrypted_detection.value.field_name
}

verified_by[[detection, verifications]] {
    detection := input.target_detections[_]
    detection.value.field_name != ""
    detection.value.object_name != ""

    verifications := [verification |
                        encrypted_detection := ruby_encrypted[_]
                        detection.value.object_name == encrypted_detection.value.object_name
                        detection.value.field_name == encrypted_detection.value.field_name
                        verification := {
                            "detector": "detect_encrypted_ruby_class_properties",
                            "filename": encrypted_detection.source.filename,
                            "line_number": encrypted_detection.source.line_number
                        }]

    count(verifications) != 0
}
