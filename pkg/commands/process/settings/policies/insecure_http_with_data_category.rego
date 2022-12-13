package bearer.insecure_http_with_data_category

import data.bearer.common

import future.keywords

insecure_http_with_data contains [data_detector_id, insecure_detector_id, item] if {
    some data_risk in input.dataflow.risks
    data_detector_id := data_risk.detector_id

    data_type = data_risk.data_types[_]
    data_location = data_type.locations[_]

    some insecure_risk in input.dataflow.risks
    insecure_detector_id := insecure_risk.detector_id
    some insecure_location in insecure_risk.locations
    data_location.filename == insecure_location.filename
    data_location.parent.line_number == insecure_location.parent.line_number

    item := {
        "category_groups": data.bearer.common.groups_for_datatype(data_type),
        "severity": data.bearer.common.severity_of_datatype(data_type),
        "filename": data_location.filename,
        "line_number": data_location.line_number,
        "parent_line_number": data_location.parent.line_number,
        "parent_content": data_location.parent.content
    }
}

policy_failure contains item if {
  insecure_http_with_data[["ruby_http_get_detection", "ruby_http_get_insecure", item]]
}

policy_failure contains item if {
  insecure_http_with_data[["ruby_http_post_detection", "ruby_http_post_insecure", item]]
}
