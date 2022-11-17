package bearer.encryption_fails

import future.keywords


detectorID := "missing_encrypted_property"

rubyEncrypted[location] {
    some detection in input
    detection.detector_type == "detect_encrypted_ruby_class_properties"
    detection.value.classification.decision.state == "valid"
    location = detection
}

sqlTables[location] {
    some detection in input
    detection.detector_type == "detect_sql_create_public_table"
    detection.value.classification.decision.state == "valid"
    detection.value.field_name != ""
    detection.value.object_name != ""

    location = detection
}

sqlTablesEncrypted[location] {
    some detection in sqlTables
    rubyEncrypted[_].value.object_name == detection.value.object_name
    rubyEncrypted[_].value.field_name == detection.value.field_name
    
    location = detection
}

sqlTablesEncrypted[location] {
    some detection in sqlTables
    startswith(detection.value.field_name, "tanker_encrypted")
    location = detection
}

risks[risk] {
    some detection in sqlTables
    not detection in sqlTablesEncrypted
    
    risk = {
        "detector_id": detectorID,
        "data_types": [{
            "name": detection.value.classification.name,
            "stored": true,
            "locations": [{
                "filename": detection.source.filename,
                "line_number": detection.source.line_number
            }]
        }]
    }
}