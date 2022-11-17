package bearer.logger_leaks

import future.keywords

policy_id := "detect_ruby_logger"
policy_name := input.name
policy_description := input.description

locations[item] {
    some detector in input.dataflow.risks
    detector.detector_id == policy_id

    data_type = detector.data_types[_]

    some category in input.data_categories
    category.name == data_type.category

    location = data_type.locations[_]
    item := {
        "data_type": data_type.name,
        "severity": category.severity,
        "filename": location.filename,
        "line_number": location.line_number
    }
}
