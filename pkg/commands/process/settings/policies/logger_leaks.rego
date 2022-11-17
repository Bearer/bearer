package bearer.logger_leaks

import future.keywords

policy_id := "detect_ruby_logger"
policy_name := "Logger leak"
policy_description := "Logger leak"

locations[item] {
    some detector in input.dataflow.risks
    detector.detector_id == policy_id

    data_type = detector.data_types[_]

    some category in input.data_categories
    category.name == data_type.category

    location = data_type.locations[_]
    item := {
        "policy": policy_name,
        "data_type": data_type.name,
        "severity": category.severity,
        "filename": location.filename,
        "line_number": location.line_number
    }
}
