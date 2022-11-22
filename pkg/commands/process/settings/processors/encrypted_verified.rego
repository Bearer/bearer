package bearer.encrypted_verified

import future.keywords


default encrypted := false


rubyEncrypted[location] {
    some detection in input.all_detections
    detection.detector_type == "detect_encrypted_ruby_class_properties"
    detection.value.classification.decision.state == "valid"
    location = detection
}

encrypted = true {
    some detection in rubyEncrypted
    detection.value.object_name == input.target.value.object_name
    detection.value.field_name == input.target.value.field_name
    input.target.value.field_name != ""
    input.target.value.object_name != ""
}

verifiedBy[verification] {
    some detection in rubyEncrypted
    detection.value.object_name == input.target.value.object_name
    detection.value.field_name == input.target.value.field_name

    verification = {
        "detector": "detect_encrypted_ruby_class_properties",
        "filename": detection.source.filename,
        "line_number": detection.source.line_number
    }
}