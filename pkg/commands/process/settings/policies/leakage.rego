package bearer.leakage

import future.keywords

sensitive_data_group_uuid := "f6a0c071-5908-4420-bac2-bba28d41223e"
personal_data_group_uuid := "e1d3135b-3c0f-4b55-abce-19f27a26cbb3"

high[item] {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    data_type = detector.data_types[_]

    some category in input.data_categories
    category.uuid == data_type.category_uuid
    category.group_uuid == sensitive_data_group_uuid

    location = data_type.locations[_]
    item := {
        "category_group": category.group_name,
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": data_type.parent.line_number,
        "parent_content": data_type.parent.content
    }
}

critical[item] {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    data_type = detector.data_types[_]

    some category in input.data_categories
    category.uuid == data_type.category_uuid
    category.group_uuid == personal_data_group_uuid

    location = data_type.locations[_]
    item := {
        "category_group": category.group_name,
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": data_type.parent.line_number,
        "parent_content": data_type.parent.content
    }
}
