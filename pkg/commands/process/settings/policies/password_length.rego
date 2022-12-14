package bearer.password_length

import data.bearer.common

import future.keywords

policy_breach contains item if {
    some data_type in input.dataflow.data_types

    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    location = detector.locations[_]
    item := {
        "category_groups": data.bearer.common.groups_for_datatypes(input.dataflow.data_types),
        "severity": "medium",
        "filename": location.filename,
        "line_number": location.line_number,
        "omit_parent": true
    }
}
