package bearer.weak_password_encryption

import data.bearer.common

import future.keywords

password_uuid := "02bb0d3a-2c8c-4842-be1c-c057f0dccd63"

policy_breach contains item if {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    data_type = detector.data_types[_]
    data_type.uuid == password_uuid

    location = data_type.locations[_]
    item := {
        "category_groups": data.bearer.common.groups_for_datatype(data_type),
        "severity": data.bearer.common.severity_of_datatype(data_type),
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}
