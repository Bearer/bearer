package bearer.insecure_smtp

import data.bearer.common

import future.keywords

policy_breach contains item if {
    some data_type in input.dataflow.data_types
    data.bearer.common.has_sensitive_data(data_type)

    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    location = detector.locations[_]
    item := {
        "category_group": data.bearer.common.groups_for_datatype(data_type),
        "severity": "medium",
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}
