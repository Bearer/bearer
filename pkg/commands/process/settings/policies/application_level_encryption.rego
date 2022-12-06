package bearer.application_level_encryption

import data.bearer.common

import future.keywords

policy_breach contains item if {
    some datatype in input.dataflow.data_types
    some detector in datatype.detectors
    detector.name == input.policy_id

    some location in detector.locations
    not location.encrypted

    item := {
        "category_groups": data.bearer.common.groups_for_datatype(datatype),
        "severity": data.bearer.common.severity_of_datatype(datatype),
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": detector.parent.line_number,
        "parent_content": detector.parent.content
    }
}