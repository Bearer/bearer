package bearer.logger_leaks

import future.keywords

sensitive_group_uuid := "f6a0c071-5908-4420-bac2-bba28d41223e"

result[item] {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    data_type = detector.data_types[_]

    some category in input.data_categories
    category.uuid == data_type.category_uuid
    category.group_uuid == sensitive_group_uuid

    location = data_type.locations[_]
    item := {
        "policy_id": input.policy_id,
        "policy_name": input.policy_name,
        "policy_description": input.policy_description,
        "severity": "critical",
        "data_type": data_type.name,
        "filename": location.filename,
        "line_number": location.line_number
    }
}

result[item] {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    data_type = detector.data_types[_]

    some category in input.data_categories
    category.uuid == data_type.category_uuid
    category.group_uuid != sensitive_group_uuid

    location = data_type.locations[_]
    item := {
        "policy_id": input.policy_id,
        "policy_name": input.policy_name,
        "policy_description": input.policy_description,
        "severity": "medium",
        "data_type": data_type.name,
        "filename": location.filename,
        "line_number": location.line_number
    }
}
