package bearer.application_level_encryption

import future.keywords

sensitive_data_group_uuid := "f6a0c071-5908-4420-bac2-bba28d41223e"
personal_data_group_uuid := "e1d3135b-3c0f-4b55-abce-19f27a26cbb3"

high[item] {
    some datatype in input.dataflow.data_types    
    some detector in datatype.detectors
    detector.name == input.policy_id
    
    some location in detector.locations
    not location.encrypted

    some category in input.data_categories
    category.uuid == datatype.category_uuid
    category.group_uuid == sensitive_data_group_uuid

    item = {
        "category_group":  category.group_name,
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": detector.parent.line_number,
        "parent_content": detector.parent.content

    }
}

critical[item] {
    some datatype in input.dataflow.data_types    
    some detector in datatype.detectors
    detector.name == input.policy_id
    
    some location in detector.locations
    not location.encrypted

    some category in input.data_categories
    category.uuid == datatype.category_uuid
    category.group_uuid == personal_data_group_uuid

    item = {
        "category_group":  category.group_name,
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": detector.parent.line_number,
        "parent_content": detector.parent.content
    }
}