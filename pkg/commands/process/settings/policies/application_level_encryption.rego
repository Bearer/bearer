package bearer.application_level_encryption

import data.bearer.common

import future.keywords

policy_failure contains item if {
    datatype = input.dataflow.data_types[_]
    datatype.name != "Unique Identifier"

    detector = datatype.detectors[_]
    location = detector.locations[_]

    location.stored

    not location.encrypted


    item := {
        "category_groups": data.bearer.common.groups_for_datatype(datatype),
        "severity": "low",
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}