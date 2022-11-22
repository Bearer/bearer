package bearer.logger_leaks

import future.keywords

result[item] {
    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    data_type = detector.data_types[_]

    some category in input.data_categories
    category.uuid == data_type.category_uuid

    location = data_type.locations[_]
    item := {
        "policy_description": input.policy_description,
        "policy_id": input.policy_id,
        "policy_name": input.policy_name,
        "data_type": data_type.name,
        "filename": location.filename,
        "line_number": location.line_number
    }
}
