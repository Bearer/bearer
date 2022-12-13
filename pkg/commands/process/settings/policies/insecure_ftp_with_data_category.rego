package bearer.insecure_ftp_with_data_category

import data.bearer.common

import future.keywords

policy_failure contains item if {
    some risk in input.dataflow.risks
    risk.detector_id == "detect_rails_insecure_ftp_data"

    data_type = risk.data_types[_]
    location = data_type.locations[_]

    item := {
        "category_groups": data.bearer.common.groups_for_datatype(data_type),
        "severity": data.bearer.common.severity_of_datatype(data_type),
        "filename": location.filename,
        "line_number": location.line_number,
        "parent_line_number": location.parent.line_number,
        "parent_content": location.parent.content
    }
}
