package bearer.http_get_parameters

import future.keywords

sensitive_data_group_uuid := "f6a0c071-5908-4420-bac2-bba28d41223e"
personal_data_group_uuid := "e1d3135b-3c0f-4b55-abce-19f27a26cbb3"

item_in_data_category contains [category_group_uuid, item] if {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    data_type = detector.data_types[_]

    some category in input.data_categories
    category.uuid == data_type.category_uuid
    category_group_uuid := category.group_uuid

    location = data_type.locations[_]
    item := {
        "category_group": category.group_name,
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}

high[item] {
    item_in_data_category[[sensitive_data_group_uuid, item]]
}

critical[item] {
    item_in_data_category[[personal_data_group_uuid, item]]
}
