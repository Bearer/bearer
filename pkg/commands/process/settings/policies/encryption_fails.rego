package bearer.encryption_fails

import future.keywords


rubyEncrypts[location] {
    some detector in input
    detector.detector_type == "detect_encrypted_ruby_class_properties"
    detector.value.classification.decision.state == "valid"
    location = detector
}



sqlTables[location] {
    some detector in input
    detector.detector_type == "detect_sql_create_public_table"
    detector.value.classification.decision.state == "valid"
    detector.value.field_name != ""
    detector.value.object_name != ""

    location = detector
}

sqlTablesEncrypted[location] {
    some detection in sqlTables
    rubyEncrypts[_].value.object_name == detection.value.object_name
    rubyEncrypts[_].value.field_name == detection.value.field_name
    
    location = detection
}

nonEncrypted[detection] {
    some detection in sqlTables
    not detection in sqlTablesEncrypted
}